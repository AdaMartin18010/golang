# AD-008: 系统性能优化模式 (System Performance Optimization)

> **维度**: Application Domains
> **级别**: S (16+ KB)
> **标签**: #performance #optimization #profiling #caching #scalability
> **权威来源**: [Systems Performance](https://www.brendangregg.com/systems-performance-2nd-edition.html) - Brendan Gregg

---

## 性能优化层次

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Performance Optimization Layers                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 架构层 (Architecture)                                                    │
│     ├── 水平扩展 (Sharding/Partitioning)                                     │
│     ├── 读写分离                                                             │
│     ├── 缓存策略 (CDN/Redis/Local)                                           │
│     └── 异步处理 (Queue/Event-driven)                                        │
│                                                                              │
│  2. 算法层 (Algorithm)                                                       │
│     ├── 时间复杂度优化                                                        │
│     ├── 空间换时间                                                           │
│     └── 数据结构选择                                                         │
│                                                                              │
│  3. 代码层 (Code)                                                            │
│     ├── 减少内存分配                                                         │
│     ├── 避免热点锁                                                           │
│     └── 向量化/SIMD                                                          │
│                                                                              │
│  4. 系统层 (System)                                                          │
│     ├── CPU 亲和性                                                           │
│     ├── 零拷贝                                                               │
│     └── 系统调用优化                                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Go 性能优化

### 内存优化

```go
package perf

import (
    "sync"
    "sync/atomic"
)

// ObjectPool 对象池，减少 GC 压力
type ObjectPool[T any] struct {
    pool sync.Pool
    reset func(T) T
}

func NewObjectPool[T any](newFunc func() T, reset func(T) T) *ObjectPool[T] {
    return &ObjectPool[T]{
        pool: sync.Pool{
            New: func() interface{} {
                return newFunc()
            },
        },
        reset: reset,
    }
}

func (p *ObjectPool[T]) Get() T {
    return p.pool.Get().(T)
}

func (p *ObjectPool[T]) Put(x T) {
    if p.reset != nil {
        x = p.reset(x)
    }
    p.pool.Put(x)
}

// 使用示例
var bufferPool = NewObjectPool(
    func() []byte { return make([]byte, 0, 4096) },
    func(b []byte) []byte { return b[:0] },
)

func processData() {
    buf := bufferPool.Get()
    defer bufferPool.Put(buf)

    // 使用 buf...
}
```

### 并发优化

```go
// ShardedMap 分片 Map，减少锁竞争
type ShardedMap[K comparable, V any] struct {
    shards []*shard[K, V]
    hash   func(K) uint64
}

type shard[K comparable, V any] struct {
    mu   sync.RWMutex
    data map[K]V
}

func NewShardedMap[K comparable, V any](shardCount int, hash func(K) uint64) *ShardedMap[K, V] {
    shards := make([]*shard[K, V], shardCount)
    for i := range shards {
        shards[i] = &shard[K, V]{data: make(map[K]V)}
    }
    return &ShardedMap[K, V]{
        shards: shards,
        hash:   hash,
    }
}

func (m *ShardedMap[K, V]) getShard(key K) *shard[K, V] {
    return m.shards[m.hash(key)%uint64(len(m.shards))]
}

func (m *ShardedMap[K, V]) Get(key K) (V, bool) {
    s := m.getShard(key)
    s.mu.RLock()
    defer s.mu.RUnlock()
    v, ok := s.data[key]
    return v, ok
}

func (m *ShardedMap[K, V]) Set(key K, value V) {
    s := m.getShard(key)
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key] = value
}

// 无锁计数器 (Atomic)
type AtomicCounter struct {
    value atomic.Int64
}

func (c *AtomicCounter) Inc() int64 {
    return c.value.Add(1)
}

func (c *AtomicCounter) Value() int64 {
    return c.value.Load()
}
```

---

## 缓存策略

### 多级缓存

```
请求流程:

Client
   │
   ▼
┌─────────────┐     Miss    ┌─────────────┐     Miss    ┌─────────────┐
│   L1 Cache  │────────────►│   L2 Cache  │────────────►│    Origin   │
│  (Local)    │             │  (Redis)    │             │   (DB/API)  │
│  < 1ms      │             │  < 5ms      │             │  < 100ms    │
└─────────────┘             └─────────────┘             └─────────────┘
       │                           │
       │ Hit                       │ Hit
       ▼                           ▼
   直接返回                    返回并回填 L1
```

### 缓存模式

```go
package perf

import (
    "context"
    "time"

    "github.com/patrickmn/go-cache"
)

// CacheAside Cache-Aside 模式
type CacheAside[K comparable, V any] struct {
    localCache *cache.Cache
    remoteCache RemoteCache
    loader      func(ctx context.Context, key K) (V, error)
    ttl         time.Duration
}

func (c *CacheAside[K, V]) Get(ctx context.Context, key K) (V, error) {
    var zero V

    // 1. 查本地缓存
    if val, found := c.localCache.Get(key); found {
        return val.(V), nil
    }

    // 2. 查远程缓存
    if val, err := c.remoteCache.Get(ctx, key); err == nil {
        c.localCache.Set(key, val, c.ttl)
        return val, nil
    }

    // 3. 加载源数据
    val, err := c.loader(ctx, key)
    if err != nil {
        return zero, err
    }

    // 4. 回填缓存
    c.remoteCache.Set(ctx, key, val, c.ttl)
    c.localCache.Set(key, val, c.ttl)

    return val, nil
}

// 缓存更新策略
// Write-Through: 同时写缓存和数据库
// Write-Behind: 先写缓存，异步写数据库
// Write-Around: 直接写数据库，删除缓存
```

---

## 性能分析

### Go Profiling

```bash
# CPU Profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory Profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Goroutine Profiling
curl http://localhost:6060/debug/pprof/goroutine

# Block Profiling (锁等待)
go test -blockprofile=block.prof

# 火焰图
go tool pprof -http=:8080 cpu.prof
```

### 关键指标

| 指标 | 工具 | 目标 |
|------|------|------|
| CPU | pprof | 热点函数优化 |
| Memory | pprof, trace | 减少分配 |
| Latency | trace | 识别阻塞 |
| Throughput | Benchmark | 持续优化 |

---

## 性能测试

```go
package perf

import (
    "sync"
    "testing"
    "time"
)

// BenchmarkMap 对比不同 Map 实现
func BenchmarkSyncMap(b *testing.B) {
    var m sync.Map
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            m.Store(i, i)
            m.Load(i)
            i++
        }
    })
}

func BenchmarkShardedMap(b *testing.B) {
    m := NewShardedMap[int, int](64, func(k int) uint64 {
        return uint64(k)
    })
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            m.Set(i, i)
            m.Get(i)
            i++
        }
    })
}

// 负载测试
func LoadTest(concurrency int, duration time.Duration, fn func()) {
    var wg sync.WaitGroup
    stop := make(chan struct{})

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case <-stop:
                    return
                default:
                    fn()
                }
            }
        }()
    }

    time.Sleep(duration)
    close(stop)
    wg.Wait()
}
```

---

## 参考文献

1. [Systems Performance](https://www.brendangregg.com/systems-performance-2nd-edition.html) - Brendan Gregg
2. [High Performance Go](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)
3. [Go Performance Book](https://github.com/dgryski/go-perfbook)
