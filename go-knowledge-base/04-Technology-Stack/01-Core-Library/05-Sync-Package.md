# sync 包详解

> **分类**: 开源技术堆栈

---

## Mutex

```go
var mu sync.Mutex
var count int

func increment() {
    mu.Lock()
    defer mu.Unlock()
    count++
}
```

## RWMutex

```go
var rwmu sync.RWMutex
var data map[string]string

func read(key string) string {
    rwmu.RLock()
    defer rwmu.RUnlock()
    return data[key]
}

func write(key, value string) {
    rwmu.Lock()
    defer rwmu.Unlock()
    data[key] = value
}
```

---

## WaitGroup

```go
var wg sync.WaitGroup

for i := 0; i < 3; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        work()
    }()
}

wg.Wait()
```

---

## Once

```go
var once sync.Once
var instance *Singleton

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

---

## Map

```go
var m sync.Map

m.Store("key", "value")
value, ok := m.Load("key")
m.Delete("key")
m.Range(func(key, value interface{}) bool {
    // 遍历
    return true
})
```
