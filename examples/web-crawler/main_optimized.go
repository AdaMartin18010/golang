package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// URL 待爬取的URL
type URL struct {
	URL   string
	Depth int
}

// Result 爬取结果
type Result struct {
	URL    string
	Status int
	Error  error
}

// CrawlerOptimized 优化后的并发爬虫
// 应用了多个并发模式：Worker Pool, Context, Graceful Shutdown
type CrawlerOptimized struct {
	maxDepth   int
	maxWorkers int
	visited    sync.Map // 使用sync.Map替代map+mutex
	client     *http.Client
}

// NewCrawlerOptimized 创建优化的爬虫
func NewCrawlerOptimized(maxDepth, maxWorkers int) *CrawlerOptimized {
	return &CrawlerOptimized{
		maxDepth:   maxDepth,
		maxWorkers: maxWorkers,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Crawl 开始爬取（使用Context和Worker Pool模式）
func (c *CrawlerOptimized) Crawl(ctx context.Context, startURL string) ([]Result, error) {
	results := make([]Result, 0)
	var resultsMu sync.Mutex

	// 创建可取消的context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Worker Pool模式
	urlQueue := make(chan URL, 100)
	var wg sync.WaitGroup

	// 启动workers
	for i := 0; i < c.maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			c.worker(ctx, workerID, urlQueue, &results, &resultsMu)
		}(i)
	}

	// 添加起始URL
	select {
	case urlQueue <- URL{URL: startURL, Depth: 0}:
	case <-ctx.Done():
		close(urlQueue)
		return results, ctx.Err()
	}

	// 优雅关闭：等待一段时间或context取消
	shutdownTimer := time.NewTimer(5 * time.Second)
	defer shutdownTimer.Stop()

	select {
	case <-shutdownTimer.C:
	case <-ctx.Done():
	}

	close(urlQueue)
	wg.Wait()

	return results, nil
}

// worker 工作协程（支持Context取消）
func (c *CrawlerOptimized) worker(
	ctx context.Context,
	id int,
	urls <-chan URL,
	results *[]Result,
	mu *sync.Mutex,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-urls:
			if !ok {
				return
			}

			// 使用LoadOrStore实现原子性检查和设置
			if _, loaded := c.visited.LoadOrStore(url.URL, true); loaded {
				continue
			}

			// 爬取URL
			result := c.fetchWithContext(ctx, url.URL)

			// 安全地添加结果
			mu.Lock()
			*results = append(*results, result)
			mu.Unlock()
		}
	}
}

// fetchWithContext 支持Context的fetch
func (c *CrawlerOptimized) fetchWithContext(ctx context.Context, url string) Result {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Result{URL: url, Error: err}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return Result{URL: url, Error: err}
	}
	defer resp.Body.Close()

	return Result{
		URL:    url,
		Status: resp.StatusCode,
	}
}

func main() {
	ctx := context.Background()
	crawler := NewCrawlerOptimized(2, 5)

	fmt.Println("开始爬取（优化版本）...")
	results, err := crawler.Crawl(ctx, "https://example.com")

	if err != nil {
		fmt.Printf("爬取错误: %v\n", err)
	}

	fmt.Printf("爬取完成，共%d个URL\n", len(results))
	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("❌ %s: %v\n", r.URL, r.Error)
		} else {
			fmt.Printf("✅ %s: %d\n", r.URL, r.Status)
		}
	}
}
