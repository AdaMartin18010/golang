# Go 1.26 性能优化模式

**版本**: Go 1.26
**性质**: 性能优化实战指南
**目标**: 写出高性能Go代码

---

## 目录

- [Go 1.26 性能优化模式](#go-126-性能优化模式)
  - [目录](#目录)
  - [1. 内存分配优化](#1-内存分配优化)
    - [1.1 对象池模式](#11-对象池模式)
    - [1.2 栈逃逸避免](#12-栈逃逸避免)
    - [1.3 预分配策略](#13-预分配策略)
  - [2. 切片和Map优化](#2-切片和map优化)
    - [2.1 切片追加优化](#21-切片追加优化)
    - [2.2 Map优化技巧](#22-map优化技巧)
  - [3. 字符串优化](#3-字符串优化)
    - [3.1 字符串构建](#31-字符串构建)
    - [3.2 字符串转换优化](#32-字符串转换优化)
  - [4. 并发优化](#4-并发优化)
    - [4.1 Goroutine池](#41-goroutine池)
    - [4.2 减少锁竞争](#42-减少锁竞争)
    - [4.3 Channel优化](#43-channel优化)
  - [5. 反射优化](#5-反射优化)
    - [5.1 缓存反射结果](#51-缓存反射结果)
    - [5.2 避免反射](#52-避免反射)
  - [6. 算法优化](#6-算法优化)
    - [6.1 使用合适的数据结构](#61-使用合适的数据结构)
    - [6.2 字符串搜索优化](#62-字符串搜索优化)
  - [7. Benchmark技巧](#7-benchmark技巧)
    - [7.1 正确编写Benchmark](#71-正确编写benchmark)
    - [7.2 内存分析](#72-内存分析)
    - [7.3 Profiling](#73-profiling)
  - [优化检查清单](#优化检查清单)

---

## 1. 内存分配优化

### 1.1 对象池模式

```go
package main

import (
    "sync"
)

// BufferPool 高频率分配的对象池
type BufferPool struct {
    pool sync.Pool
}

// NewBufferPool 创建指定大小的缓冲池
func NewBufferPool(size int) *BufferPool {
    return &BufferPool{
        pool: sync.Pool{
            New: func() interface{} {
                return make([]byte, 0, size)
            },
        },
    }
}

// Get 获取缓冲区
func (p *BufferPool) Get() []byte {
    return p.pool.Get().([]byte)
}

// Put 归还缓冲区
func (p *BufferPool) Put(b []byte) {
    // 重置但保留容量
    b = b[:0]
    p.pool.Put(b)
}

// 使用示例
var pool = NewBufferPool(4096)

func process(data []byte) []byte {
    buf := pool.Get()
    defer pool.Put(buf)

    buf = append(buf, data...)
    // 处理...
    return buf
}
```

### 1.2 栈逃逸避免

```go
// ❌ 导致逃逸
func escape() *int {
    x := 1
    return &x  // x逃逸到堆
}

// ✅ 保持在栈上
func noEscape() int {
    x := 1
    return x  // x在栈上分配
}

// ❌ interface{}导致逃逸
func interfaceEscape(x int) interface{} {
    return x  // 装箱，堆分配
}

// ✅ 使用具体类型
func concrete(x int) int {
    return x  // 栈分配
}

// ❌ 闭包捕获导致逃逸
func closure() func() int {
    x := 1
    return func() int {
        return x
    }
}

// ✅ 传递参数而非捕获
func noClosure(x int) func() int {
    return func() int {
        return x  // x作为参数传入
    }
}
```

### 1.3 预分配策略

```go
// 切片预分配
func filter(items []int, fn func(int) bool) []int {
    // 最坏情况: 全部匹配
    result := make([]int, 0, len(items))

    for _, item := range items {
        if fn(item) {
            result = append(result, item)
        }
    }

    return result
}

// Map预分配
func groupBy(items []Item) map[string][]Item {
    groups := make(map[string][]Item, len(items))

    for _, item := range items {
        key := item.Category
        groups[key] = append(groups[key], item)
    }

    return groups
}

// 精确预分配
func toStrings(nums []int) []string {
    result := make([]string, len(nums))  // 精确大小
    for i, n := range nums {
        result[i] = strconv.Itoa(n)
    }
    return result
}
```

---

## 2. 切片和Map优化

### 2.1 切片追加优化

```go
// 批量追加优于逐个追加
func appendBatch() {
    var s []int

    // ❌ 多次分配
    for i := 0; i < 1000; i++ {
        s = append(s, i)  // 可能多次扩容
    }

    // ✅ 预分配
    s = make([]int, 0, 1000)
    for i := 0; i < 1000; i++ {
        s = append(s, i)
    }

    // ✅ 直接索引赋值
    s = make([]int, 1000)
    for i := 0; i < 1000; i++ {
        s[i] = i
    }
}

// 使用copy进行批量操作
func copyOptimized() {
    src := make([]int, 10000)
    dst := make([]int, len(src))

    // ✅ copy比循环快
    copy(dst, src)
}
```

### 2.2 Map优化技巧

```go
// 预分配Map容量
func createMap(n int) map[int]string {
    // ✅ 预分配避免rehash
    m := make(map[int]string, n)

    for i := 0; i < n; i++ {
        m[i] = strconv.Itoa(i)
    }

    return m
}

// 使用struct{}作为Set
func useSet() {
    // ✅ struct{}不占内存
    set := make(map[string]struct{})
    set["key"] = struct{}{}

    // 检查存在性
    if _, ok := set["key"]; ok {
        // 存在
    }
}

// 删除所有元素: 重建vs遍历删除
func clearMap(m map[int]string) {
    // ❌ 遍历删除
    for k := range m {
        delete(m, k)
    }

    // ✅ 重建(大量元素时更快)
    m = make(map[int]string)
}
```

---

## 3. 字符串优化

### 3.1 字符串构建

```go
// strings.Builder 优于 +=
func buildString(n int) string {
    // ❌ 每次+都创建新字符串
    var s string
    for i := 0; i < n; i++ {
        s += strconv.Itoa(i)
    }
    return s
}

// ✅ 使用strings.Builder
func buildStringOptimized(n int) string {
    var b strings.Builder
    b.Grow(n * 4)  // 预分配

    for i := 0; i < n; i++ {
        b.WriteString(strconv.Itoa(i))
    }
    return b.String()
}

// bytes.Buffer 可重用
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func buildWithPool(data []string) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)

    buf.Reset()
    for _, s := range data {
        buf.WriteString(s)
    }
    return buf.String()
}
```

### 3.2 字符串转换优化

```go
// []byte和string转换的零拷贝技巧
func stringToBytes(s string) []byte {
    // ✅ 无拷贝转换 (unsafe)
    return *(*[]byte)(unsafe.Pointer(&struct {
        string
        Cap int
    }{s, len(s)}))
}

func bytesToString(b []byte) string {
    // ✅ 无拷贝转换
    return *(*string)(unsafe.Pointer(&b))
}

// 安全版本 (Go 1.20+)
func bytesToStringSafe(b []byte) string {
    return string(b)  // 编译器优化，可能不拷贝
}
```

---

## 4. 并发优化

### 4.1 Goroutine池

```go
package main

import (
    "context"
    "sync"
)

// WorkerPool goroutine池
type WorkerPool struct {
    workers int
    jobs    chan func()
    wg      sync.WaitGroup
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workers int) *WorkerPool {
    wp := &WorkerPool{
        workers: workers,
        jobs:    make(chan func()),
    }

    for i := 0; i < workers; i++ {
        wp.wg.Add(1)
        go wp.worker()
    }

    return wp
}

func (wp *WorkerPool) worker() {
    defer wp.wg.Done()
    for job := range wp.jobs {
        job()
    }
}

// Submit 提交任务
func (wp *WorkerPool) Submit(job func()) {
    wp.jobs <- job
}

// Close 关闭工作池
func (wp *WorkerPool) Close() {
    close(wp.jobs)
    wp.wg.Wait()
}

// 使用示例
func main() {
    pool := NewWorkerPool(10)
    defer pool.Close()

    for i := 0; i < 100; i++ {
        i := i  // 捕获循环变量
        pool.Submit(func() {
            process(i)
        })
    }
}
```

### 4.2 减少锁竞争

```go
// 分段锁减少竞争
type ShardedMap struct {
    shards [256]*shard
}

type shard struct {
    mu   sync.RWMutex
    data map[string]interface{}
}

func NewShardedMap() *ShardedMap {
    sm := &ShardedMap{}
    for i := 0; i < 256; i++ {
        sm.shards[i] = &shard{data: make(map[string]interface{})}
    }
    return sm
}

func (sm *ShardedMap) getShard(key string) *shard {
    hash := fnv32(key)
    return sm.shards[hash%256]
}

func (sm *ShardedMap) Get(key string) (interface{}, bool) {
    shard := sm.getShard(key)
    shard.mu.RLock()
    defer shard.mu.RUnlock()
    val, ok := shard.data[key]
    return val, ok
}

func (sm *ShardedMap) Set(key string, val interface{}) {
    shard := sm.getShard(key)
    shard.mu.Lock()
    defer shard.mu.Unlock()
    shard.data[key] = val
}

// 使用atomic替代Mutex
var counter atomic.Int64

counter.Add(1)  // 原子操作，无锁
```

### 4.3 Channel优化

```go
// 批量处理channel数据
func batchProcess(ch <-chan Item, batchSize int, timeout time.Duration) {
    batch := make([]Item, 0, batchSize)
    timer := time.NewTimer(timeout)
    defer timer.Stop()

    for {
        select {
        case item, ok := <-ch:
            if !ok {
                // 处理剩余数据
                if len(batch) > 0 {
                    processBatch(batch)
                }
                return
            }
            batch = append(batch, item)
            if len(batch) >= batchSize {
                processBatch(batch)
                batch = batch[:0]
                timer.Reset(timeout)
            }

        case <-timer.C:
            if len(batch) > 0 {
                processBatch(batch)
                batch = batch[:0]
            }
            timer.Reset(timeout)
        }
    }
}

// 使用有缓冲channel
ch := make(chan Item, 1000)  // 减少阻塞
```

---

## 5. 反射优化

### 5.1 缓存反射结果

```go
type structCache struct {
    fields []fieldInfo
}

type fieldInfo struct {
    index int
    name  string
}

var cache sync.Map  // map[reflect.Type]*structCache

func getStructCache(t reflect.Type) *structCache {
    if c, ok := cache.Load(t); ok {
        return c.(*structCache)
    }

    // 构建缓存
    sc := &structCache{
        fields: make([]fieldInfo, t.NumField()),
    }
    for i := 0; i < t.NumField(); i++ {
        sc.fields[i] = fieldInfo{
            index: i,
            name:  t.Field(i).Name,
        }
    }

    cache.Store(t, sc)
    return sc
}

// 使用缓存的反射
func processStruct(v interface{}) {
    t := reflect.TypeOf(v)
    cache := getStructCache(t)

    // 使用缓存的fieldInfo，避免重复反射
    for _, field := range cache.fields {
        // 处理字段
    }
}
```

### 5.2 避免反射

```go
// 使用代码生成替代反射
//go:generate go run gen.go

// 使用接口替代反射
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}

// 使用类型开关替代反射
func process(v interface{}) {
    switch x := v.(type) {
    case string:
        processString(x)
    case int:
        processInt(x)
    case Marshaler:
        x.MarshalJSON()
    }
}
```

---

## 6. 算法优化

### 6.1 使用合适的数据结构

```go
// 查找操作: Map vs Slice
// O(1) vs O(n)

// ✅ 频繁查找使用Map
lookup := map[string]int{"a": 1, "b": 2}
if v, ok := lookup["a"]; ok {
    // O(1)
}

// ✅ 有序数据使用二分查找
sort.Ints(nums)
idx := sort.SearchInts(nums, target)  // O(log n)

// ✅ 去重使用Set
set := make(map[int]struct{})
for _, v := range items {
    set[v] = struct{}{}
}
```

### 6.2 字符串搜索优化

```go
// 使用strings.Builder
var b strings.Builder
b.Grow(expectedSize)

// 预编译正则表达式
var re = regexp.MustCompile(`pattern`)

// 使用strings.Index代替正则
idx := strings.Index(s, substr)

// Boyer-Moore算法 (大量搜索)
// 使用github.com/cloudflare/buffer
```

---

## 7. Benchmark技巧

### 7.1 正确编写Benchmark

```go
// 基准测试函数
func BenchmarkXxx(b *testing.B) {
    // 准备数据
    data := make([]int, 1000)
    for i := range data {
        data[i] = i
    }

    b.ResetTimer()  // 重置计时器

    for i := 0; i < b.N; i++ {
        process(data)
    }
}

// 比较不同实现
func BenchmarkProcessOld(b *testing.B) {
    for i := 0; i < b.N; i++ {
        processOld(data)
    }
}

func BenchmarkProcessNew(b *testing.B) {
    for i := 0; i < b.N; i++ {
        processNew(data)
    }
}

// 使用benchstat比较结果
// go test -bench=. -count=5 > old.txt
// 修改代码
// go test -bench=. -count=5 > new.txt
// benchstat old.txt new.txt
```

### 7.2 内存分析

```go
// 内存基准测试
func BenchmarkAlloc(b *testing.B) {
    b.ReportAllocs()  // 报告内存分配

    for i := 0; i < b.N; i++ {
        _ = make([]byte, 1024)
    }
}

// 运行: go test -bench=BenchmarkAlloc -benchmem
```

### 7.3 Profiling

```go
// CPU Profile
// go test -cpuprofile=cpu.prof -bench=.
// go tool pprof cpu.prof

// Memory Profile
// go test -memprofile=mem.prof -bench=.
// go tool pprof mem.prof

// Trace
// go test -trace=trace.out -bench=.
// go tool trace trace.out
```

---

## 优化检查清单

- [ ] 使用sync.Pool复用高频分配的对象
- [ ] 预分配切片和Map的容量
- [ ] 使用strings.Builder代替字符串拼接
- [ ] 避免不必要的接口装箱
- [ ] 减少锁粒度，使用分段锁或原子操作
- [ ] 缓存反射结果，或避免使用反射
- [ ] 选择合适的数据结构
- [ ] 批量处理channel数据
- [ ] 使用有缓冲channel
- [ ] 使用对象池管理goroutine

---

**文档版本**: 1.0
**最后更新**: 2026-03-08
