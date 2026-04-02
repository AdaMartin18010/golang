# 逃逸分析 (Escape Analysis)

> **分类**: 工程与云原生
> **标签**: #performance #memory #gc

---

## 什么是逃逸分析

编译器决定变量分配在栈上还是堆上的过程。

```
栈分配: 快速，自动回收
堆分配: 慢，需要 GC
```

---

## 逃逸场景

### 1. 返回指针

```go
// ❌ 逃逸到堆
func NewUser(name string) *User {
    u := &User{Name: name}  // 逃逸
    return u
}

// ✅ 栈分配
func CreateUser(name string) User {
    u := User{Name: name}   // 栈上
    return u
}
```

### 2. 接口装箱

```go
// ❌ 逃逸
func Print(v interface{}) {
    fmt.Println(v)
}

Print(42)  // int 装箱到堆

// ✅ 避免装箱
func PrintInt(v int) {
    fmt.Println(v)
}
```

### 3. Slice 引用

```go
// ❌ 逃逸
data := make([]byte, 1024)
process(&data[0])

// ✅ 不逃逸
process(data)
```

### 4. 闭包引用

```go
// ❌ 逃逸
func makeCounter() func() int {
    count := 0           // 逃逸到堆
    return func() int {
        count++
        return count
    }
}
```

---

## 逃逸分析命令

```bash
# 查看逃逸分析
go build -gcflags="-m" main.go

# 详细输出
go build -gcflags="-m -m" main.go
```

### 输出解读

```
main.go:5:6: can inline NewUser
main.go:6:9: &User literal escapes to heap  ← 逃逸
main.go:10:6: can inline CreateUser
main.go:11:9: CreateUser … does not escape  ← 不逃逸
```

---

## 优化技巧

### 减少指针使用

```go
// ❌ 逃逸
type Config struct {
    Options *Options
}

// ✅ 内联
type Config struct {
    Options Options
}
```

### 预分配 Slice

```go
// ❌ 多次分配
var results []int
for i := 0; i < 100; i++ {
    results = append(results, i)  // 多次扩容
}

// ✅ 预分配
results := make([]int, 0, 100)
for i := 0; i < 100; i++ {
    results = append(results, i)
}
```

### 使用 sync.Pool

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func process() {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf[:4096])

    // 使用 buf
}
```

---

## 性能对比

| 分配方式 | 延迟 | 适用场景 |
|----------|------|----------|
| 栈 | ~1ns | 局部变量 |
| 堆 | ~100ns | 逃逸变量 |
| GC | ~ms | 堆回收 |

---

## 调试逃逸

```go
// 使用 runtime 查看分配
import "runtime"

func printAlloc() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("Alloc = %v KB\n", m.Alloc/1024)
}
```
