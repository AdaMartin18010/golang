# Slice 内部实现 (Slice Internals)

> **分类**: 语言设计
> **标签**: #slice #runtime #internals

---

## Slice 结构

```go
// runtime 中的 slice 定义
type slice struct {
    array unsafe.Pointer  // 底层数组指针
    len   int             // 长度
    cap   int             // 容量
}
```

```
Slice Header
┌─────────────────┐
│ array (pointer) │ ──> ┌───┬───┬───┬───┬───┐
│ len = 3         │     │ A │ B │ C │ D │ E │  (底层数组)
│ cap = 5         │     └───┴───┴───┴───┴───┘
└─────────────────┘       ▲   ▲   ▲
                          │   │   │
                         [0] [1] [2]
```

---

## 创建与扩容

### make 分配

```go
s := make([]int, 3, 5)
// len=3, cap=5
// 底层数组: [0, 0, 0, _, _]
```

### 扩容策略

```go
// 扩容规则
cap < 1024:    新 cap = 旧 cap * 2
cap >= 1024:   新 cap = 旧 cap * 1.25

// 示例
s := make([]int, 0, 100)
s = append(s, 1)  // cap=100 (不变)

s := make([]int, 0, 1024)
s = append(s, 1)  // cap=1024 (不变)
s = append(s, make([]int, 1024)...)  // cap=1280 (1024 * 1.25)
```

### 扩容源码逻辑

```go
func growslice(et *_type, old slice, cap int) slice {
    newcap := old.cap
    doublecap := newcap + newcap

    if cap > doublecap {
        newcap = cap
    } else {
        if old.cap < 1024 {
            newcap = doublecap
        } else {
            for newcap < cap {
                newcap += newcap / 4
            }
        }
    }

    // 内存对齐
    // ...
}
```

---

## 切片操作

### 切片表达式

```go
s := []int{0, 1, 2, 3, 4}

s[1:3]   // [1, 2]       len=2, cap=4
s[1:3:3] // [1, 2]       len=2, cap=2 (限制 cap)
s[:0]    // []           len=0, cap=5
```

### 共享底层数组

```go
s1 := []int{1, 2, 3, 4, 5}
s2 := s1[1:3]  // [2, 3]

s2[0] = 100
// s1: [1, 100, 3, 4, 5]
// s2: [100, 3]

// 解决: 复制
s2 := make([]int, 2)
copy(s2, s1[1:3])
```

---

## 内存泄漏

```go
// ❌ 内存泄漏
func process() {
    data := make([]byte, 1024*1024)  // 1MB

    // 只使用一小部分
    header := data[:100]

    // header 仍然引用整个 1MB 数组
    return header
}

// ✅ 正确做法
func process() []byte {
    data := make([]byte, 1024*1024)

    // 复制需要的数据
    header := make([]byte, 100)
    copy(header, data[:100])

    return header  // 只返回 100 字节
}
```

---

## 性能优化

### 预分配容量

```go
// ❌ 多次分配
var s []int
for i := 0; i < 1000; i++ {
    s = append(s, i)
}

// ✅ 一次分配
s := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    s = append(s, i)
}
```

### 复用切片

```go
var pool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func process() {
    buf := pool.Get().([]byte)
    defer pool.Put(buf[:4096])  // 重置长度

    // 使用 buf
}
```

---

## nil vs 空切片

```go
var s1 []int        // nil slice
s2 := []int(nil)    // nil slice
s3 := []int{}       // empty slice
s4 := make([]int, 0) // empty slice

// 区别
s1 == nil  // true
s3 == nil  // false

// JSON 序列化
json.Marshal(s1) // null
json.Marshal(s3) // []

// 最佳实践: 返回空切片
func getItems() []Item {
    items := []Item{}  // 返回 [] 而不是 null
    return items
}
```
