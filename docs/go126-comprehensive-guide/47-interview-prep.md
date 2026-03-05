# Go面试准备指南

> 常见问题、算法实现与系统设计面试准备

---

## 一、语言基础问题

### 1.1 类型系统

```text
Q: Go的接口是什么类型？与其他语言有何不同？
────────────────────────────────────────

回答要点：

1. 隐式实现：
   - 无需显式声明实现接口
   - 只要方法集匹配，就自动实现
   - 支持鸭子类型

2. 非侵入式设计：
   ```go
   type Reader interface {
       Read(p []byte) (n int, err error)
   }

   // 不需要 MyReader implements Reader
   type MyReader struct{}
   func (r *MyReader) Read(p []byte) (n int, err error) {
       // 实现
   }
   ```

1. 与其他语言对比：
   - Java：显式 implements 关键字
   - C++：抽象基类 + 纯虚函数
   - Go：更灵活，松耦合

Q: nil接口和nil指针有什么区别？
────────────────────────────────────────

关键点：

nil接口 = nil类型 + nil值
var r io.Reader  // nil接口

nil指针 = 具体类型 + nil值
var p *bytes.Buffer  // nil指针

常见陷阱：

func returnsNil() error {
    var p *MyError = nil
    return p  // 返回的是 (*MyError, nil) 不是 nil接口
}

err := returnsNil()
if err != nil {
    // 会进入这里！因为err包含类型信息
    fmt.Println(err != nil)  // true
}

正确做法：

func returnsNil() error {
    return nil  // 返回真正的nil
}

```

### 1.2 并发问题

```text
Q: 什么是数据竞争？如何检测？
────────────────────────────────────────

数据竞争定义：
- 两个goroutine同时访问同一内存位置
- 至少一个是写操作
- 没有同步机制

示例：
var counter int

func increment() {
    for i := 0; i < 1000; i++ {
        counter++  // 数据竞争！
    }
}

func main() {
    go increment()
    go increment()
    time.Sleep(time.Second)
    fmt.Println(counter)  // 结果不确定
}

检测方法：
1. 运行竞态检测器：
   go run -race main.go

2. 静态分析：
   go vet ./...

3. 代码审查：
   - 共享变量检查
   - 锁的配对检查

Q: channel和mutex如何选择？
────────────────────────────────────────

选择channel的场景：
- goroutine之间传递数据
- 通知/信号机制
- 分发工作
- 处理流的管道

选择mutex的场景：
- 保护共享状态
- 缓存、计数器
- 需要细粒度锁控制

对比示例：

// Channel方式 - 计数器
type Counter struct {
    inc   chan struct{}
    value chan int
}

func (c *Counter) Run() {
    count := 0
    for {
        select {
        case <-c.inc:
            count++
        case c.value <- count:
        }
    }
}

// Mutex方式 - 计数器
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

原则：
- 优先考虑channel（更符合Go哲学）
- 但mutex在保护状态时更简洁
- 避免过度工程化
```

---

## 二、算法与数据结构

### 2.1 常见算法实现

```text
LRU缓存实现：
────────────────────────────────────────

type LRUCache struct {
    capacity int
    cache    map[int]*list.Element
    ll       *list.List
}

type entry struct {
    key   int
    value int
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        cache:    make(map[int]*list.Element),
        ll:       list.New(),
    }
}

func (c *LRUCache) Get(key int) int {
    if ele, ok := c.cache[key]; ok {
        c.ll.MoveToFront(ele)
        return ele.Value.(*entry).value
    }
    return -1
}

func (c *LRUCache) Put(key int, value int) {
    if ele, ok := c.cache[key]; ok {
        c.ll.MoveToFront(ele)
        ele.Value.(*entry).value = value
        return
    }

    if c.ll.Len() >= c.capacity {
        back := c.ll.Back()
        c.ll.Remove(back)
        delete(c.cache, back.Value.(*entry).key)
    }

    ele := c.ll.PushFront(&entry{key, value})
    c.cache[key] = ele
}

实现要点：
- 使用container/list作为双向链表
- 移动元素到头部表示最近使用
- 淘汰时移除尾部元素
- 时间复杂度：Get和Put都是O(1)

并发的最小栈：
────────────────────────────────────────

type MinStack struct {
    stack    []int
    minStack []int
    mu       sync.RWMutex
}

func NewMinStack() *MinStack {
    return &MinStack{
        stack:    make([]int, 0),
        minStack: make([]int, 0),
    }
}

func (s *MinStack) Push(val int) {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.stack = append(s.stack, val)
    if len(s.minStack) == 0 || val <= s.minStack[len(s.minStack)-1] {
        s.minStack = append(s.minStack, val)
    }
}

func (s *MinStack) Pop() {
    s.mu.Lock()
    defer s.mu.Unlock()

    if len(s.stack) == 0 {
        return
    }

    top := s.stack[len(s.stack)-1]
    s.stack = s.stack[:len(s.stack)-1]

    if top == s.minStack[len(s.minStack)-1] {
        s.minStack = s.minStack[:len(s.minStack)-1]
    }
}

func (s *MinStack) GetMin() int {
    s.mu.RLock()
    defer s.mu.RUnlock()

    if len(s.minStack) == 0 {
        return 0
    }
    return s.minStack[len(s.minStack)-1]
}
```

### 2.2 并发数据结构

```text
无锁队列（Ring Buffer）：
────────────────────────────────────────

type RingBuffer struct {
    buffer []interface{}
    head   uint64
    tail   uint64
    size   uint64
}

func NewRingBuffer(size int) *RingBuffer {
    return &RingBuffer{
        buffer: make([]interface{}, size),
        size:   uint64(size),
    }
}

func (r *RingBuffer) Enqueue(item interface{}) bool {
    for {
        tail := atomic.LoadUint64(&r.tail)
        head := atomic.LoadUint64(&r.head)

        if (tail+1)%r.size == head%r.size {
            return false  // 满
        }

        if atomic.CompareAndSwapUint64(&r.tail, tail, tail+1) {
            r.buffer[tail%r.size] = item
            return true
        }
    }
}

func (r *RingBuffer) Dequeue() (interface{}, bool) {
    for {
        head := atomic.LoadUint64(&r.head)
        tail := atomic.LoadUint64(&r.tail)

        if head == tail {
            return nil, false  // 空
        }

        item := r.buffer[head%r.size]
        if atomic.CompareAndSwapUint64(&r.head, head, head+1) {
            return item, true
        }
    }
}
```

---

## 三、系统设计

### 3.1 设计短链接服务

```text
需求：
────────────────────────────────────────

功能：
- 输入长URL，生成短链接
- 访问短链接，重定向到长URL
- 支持自定义短链接
- 统计访问次数

非功能需求：
- 高可用
- 低延迟（<10ms）
- 支持10亿条记录

API设计：
────────────────────────────────────────

POST /api/shorten
Request: {"url": "https://example.com/very/long/url", "custom": "mylink"}
Response: {"short_url": "abc123", "expires_at": "2025-01-01"}

GET /{shortCode}
Response: 302 Redirect to original URL

GET /api/stats/{shortCode}
Response: {"clicks": 1000, "last_access": "2024-01-15"}

数据库设计：
────────────────────────────────────────

表结构：

CREATE TABLE short_urls (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    long_url VARCHAR(2048) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    click_count BIGINT DEFAULT 0
);

索引：
- short_code (唯一索引，用于查找)
- expires_at (用于清理过期数据)

短码生成：
────────────────────────────────────────

方案1：Base62编码
- 自增ID转Base62 (0-9a-zA-Z)
- abc123 → 62^6 = 568亿组合

方案2：Hash
- MD5/SHA1前6位
- 冲突解决：+1重试

方案3：随机
- 随机6位字符
- 检查数据库是否已存在

推荐方案1（简单可靠）：

func encode(id uint64) string {
    const base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    var result []byte
    for id > 0 {
        result = append([]byte{base62[id%62]}, result...)
        id /= 62
    }
    return string(result)
}

系统架构：
────────────────────────────────────────

                ┌─────────────┐
                │   Client    │
                └──────┬──────┘
                       │
                ┌──────▼──────┐
                │   CDN       │
                └──────┬──────┘
                       │
                ┌──────▼──────┐
                │  Load Balancer
                └──────┬──────┘
                       │
            ┌──────────┼──────────┐
            │          │          │
       ┌────▼───┐ ┌────▼───┐ ┌────▼───┐
       │ Go App │ │ Go App │ │ Go App │
       └───┬────┘ └───┬────┘ └───┬────┘
           │          │          │
       ┌───┴──────────┴──────────┴───┐
       │        Redis Cluster        │
       │     (缓存热数据)            │
       └─────────────┬───────────────┘
                     │
              ┌──────▼──────┐
              │   MySQL     │
              │  (主从复制) │
              └─────────────┘

缓存策略：
- 读：先查Redis，miss再查MySQL
- 写：先写MySQL，再删Redis缓存
- 热点数据预热
```

### 3.2 设计Rate Limiter

```text
算法选择：
────────────────────────────────────────

Token Bucket（推荐）：
- 容量固定，匀速产生token
- 突发流量友好
- 实现简单

Sliding Window：
- 精确计数
- 内存占用大
- 需要记录每次请求

Fixed Window：
- 简单，但可能有边界突刺

分布式实现：
────────────────────────────────────────

使用Redis + Lua：

-- rate_limiter.lua
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])

local current = redis.call("GET", key)
if current == false then
    current = 0
else
    current = tonumber(current)
end

if current >= limit then
    return 0  -- 拒绝
end

redis.call("INCR", key)
redis.call("EXPIRE", key, window)
return 1  -- 允许

Go实现：

type RedisLimiter struct {
    client *redis.Client
    limit  int
    window time.Duration
    script *redis.Script
}

func NewRedisLimiter(client *redis.Client, limit int, window time.Duration) *RedisLimiter {
    return &RedisLimiter{
        client: client,
        limit:  limit,
        window: window,
        script: redis.NewScript(luaScript),
    }
}

func (l *RedisLimiter) Allow(ctx context.Context, key string) bool {
    result, err := l.script.Run(ctx, l.client, []string{key},
        l.limit, int(l.window.Seconds())).Result()
    if err != nil {
        return true  // 降级策略：允许
    }
    return result.(int64) == 1
}
```

---

## 四、代码Review问题

```text
找出代码问题：
────────────────────────────────────────

func processItems(items []Item) error {
    for _, item := range items {
        go func() {
            if err := process(item); err != nil {
                return err  // 编译错误！
            }
        }()
    }
    return nil
}

问题：
1. 循环变量捕获：item被所有goroutine共享
2. 错误处理：goroutine内return不会返回给调用者
3. 没有等待goroutine完成
4. 可能创建过多goroutine

修复：

func processItems(items []Item) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(items))
    sem := make(chan struct{}, 10)  // 限制并发数

    for _, item := range items {
        wg.Add(1)
        go func(i Item) {
            defer wg.Done()

            sem <- struct{}{}
            defer func() { <-sem }()

            if err := process(i); err != nil {
                errChan <- err
            }
        }(item)
    }

    go func() {
        wg.Wait()
        close(errChan)
    }()

    for err := range errChan {
        if err != nil {
            return err
        }
    }
    return nil
}
```

---

*本章提供了Go面试的核心知识点和常见问题。*
