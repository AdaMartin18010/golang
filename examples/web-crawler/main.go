package main

import (
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

// Crawler 简单的并发爬虫（未优化版本）
type Crawler struct {
	maxDepth   int
	maxWorkers int
	visited    map[string]bool
	mu         sync.Mutex
	wg         sync.WaitGroup
}

// NewCrawler 创建爬虫
func NewCrawler(maxDepth, maxWorkers int) *Crawler {
	return &Crawler{
		maxDepth:   maxDepth,
		maxWorkers: maxWorkers,
		visited:    make(map[string]bool),
	}
}

// Crawl 开始爬取
func (c *Crawler) Crawl(startURL string) []Result {
	results := make([]Result, 0)
	resultsChan := make(chan Result)
	urlQueue := make(chan URL, 100)

	// 启动工作协程
	for i := 0; i < c.maxWorkers; i++ {
		c.wg.Add(1)
		go c.worker(urlQueue, resultsChan)
	}

	// 启动结果收集协程
	done := make(chan struct{})
	go func() {
		for result := range resultsChan {
			results = append(results, result)
		}
		close(done)
	}()

	// 添加起始URL
	urlQueue <- URL{URL: startURL, Depth: 0}

	// 等待所有工作完成
	time.Sleep(5 * time.Second) // 简单等待（不够优雅）
	close(urlQueue)
	c.wg.Wait()
	close(resultsChan)
	<-done

	return results
}

// worker 工作协程
func (c *Crawler) worker(urls <-chan URL, results chan<- Result) {
	defer c.wg.Done()

	for url := range urls {
		// 检查是否已访问
		c.mu.Lock()
		if c.visited[url.URL] {
			c.mu.Unlock()
			continue
		}
		c.visited[url.URL] = true
		c.mu.Unlock()

		// 爬取URL
		result := c.fetch(url.URL)
		results <- result
	}
}

// fetch 获取URL内容
func (c *Crawler) fetch(url string) Result {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
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
	crawler := NewCrawler(2, 5)
	results := crawler.Crawl("https://example.com")

	fmt.Printf("爬取完成，共%d个URL\n", len(results))
	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("❌ %s: %v\n", r.URL, r.Error)
		} else {
			fmt.Printf("✅ %s: %d\n", r.URL, r.Status)
		}
	}
}
