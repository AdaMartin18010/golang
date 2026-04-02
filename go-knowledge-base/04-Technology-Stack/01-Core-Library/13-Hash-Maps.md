# Map 内部实现与优化 (Hash Maps)

> **分类**: 开源技术堆栈  
> **标签**: #map #internals #performance

---

## Map 结构

```go
// runtime 中的 hmap 定义
type hmap struct {
    count     int           // 元素个数
    flags     uint8         // 标志位
    B         uint8         // 桶数 = 2^B
    noverflow uint16        // 溢出桶数
    hash0     uint32        // 哈希种子
    
    buckets    unsafe.Pointer  // 桶数组
    oldbuckets unsafe.Pointer  // 扩容时的旧桶
    nevacuate  uintptr         // 扩容进度
}

// 桶结构
type bmap struct {
    tophash [8]uint8  // 哈希高8位
    // 后面跟着 key/value/overflow 指针
}
```

---

## 哈希冲突解决

### 链地址法

```
Bucket 0: [k1|v1] -> [k2|v2] -> nil
Bucket 1: [k3|v3] -> nil
Bucket 2: nil
...
```

---

## 扩容机制

### 触发条件

```
负载因子 > 6.5  或
溢出桶数量 > 桶数量
```

### 渐进式扩容

```go
// 扩容时创建新桶
func hashGrow(t *maptype, h *hmap) {
    // 分配新桶数组
    nextOverflow := make([]bmap, 1<<(h.B+1))
    
    // 标记旧桶
    for i := range h.buckets {
        b := &h.buckets[i]
        b.tophash[0] = evacuatedEmpty
    }
    
    // 交换
    h.oldbuckets = h.buckets
    h.buckets = nextOverflow
    h.B++
}

// 渐进式迁移
func evacuate(t *maptype, h *hmap, oldbucket uintptr) {
    b := (*bmap)(add(h.oldbuckets, oldbucket*uintptr(t.bucketsize)))
    
    for b != nil {
        // 重新哈希到新桶
        for i := 0; i < 8; i++ {
            if b.tophash[i] == empty {
                continue
            }
            
            k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
            v := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.valuesize))
            
            // 计算新桶位置
            hash := t.hasher(k, uintptr(h.hash0))
            y := &xy[hash&newbit]  // 新桶
            
            // 迁移
            y.keys[idx] = k
            y.values[idx] = v
        }
        
        b = b.overflow(t)
    }
}
```

---

## 并发安全

### 为什么 map 不是线程安全的

```go
// ❌ 数据竞争
var m = make(map[string]int)

go func() {
    m["key"] = 1  // 写
}()

go func() {
    _ = m["key"]   // 读 - 竞争!
}()
```

### 解决方案

```go
// 1. sync.Map
var sm sync.Map
sm.Store("key", value)
v, ok := sm.Load("key")

// 2. 读写锁
var (
    mu sync.RWMutex
    m  = make(map[string]int)
)

mu.RLock()
v := m["key"]
mu.RUnlock()

mu.Lock()
m["key"] = v
mu.Unlock()

// 3. 分片锁（高性能）
type ShardedMap struct {
    shards [32]*shard
}

type shard struct {
    mu sync.RWMutex
    m  map[string]interface{}
}

func (sm *ShardedMap) getShard(key string) *shard {
    hash := fnv32(key)
    return sm.shards[hash%32]
}
```

---

## 性能优化

### 预分配

```go
// ✅ 预分配减少扩容
m := make(map[string]int, 1000)

// ❌ 多次扩容
m := make(map[string]int)
for i := 0; i < 1000; i++ {
    m[fmt.Sprintf("key%d", i)] = i
}
```

### 值类型 vs 指针

```go
// ✅ 小值类型直接存储
type Config struct {
    Timeout int
    Retries int
}
m := make(map[string]Config)

// ❌ 大结构体用指针
type LargeData struct {
    Data [1024]byte
}
m := make(map[string]*LargeData)
```

---

## 内存布局

```
Bucket:
┌─────────────────────────────┐
│ tophash[0] tophash[1] ...   │  8 bytes
├─────────────────────────────┤
│ key[0] key[1] ...           │  8*keySize
├─────────────────────────────┤
│ value[0] value[1] ...       │  8*valueSize
├─────────────────────────────┤
│ overflow pointer            │  8 bytes
└─────────────────────────────┘
```

---

## 遍历顺序

```go
// 遍历顺序随机！
for k, v := range m {
    fmt.Println(k, v)
}

// 如需顺序，先排序 keys
var keys []string
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)

for _, k := range keys {
    fmt.Println(k, m[k])
}
```
