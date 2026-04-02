# 常量 (Constants)

> **分类**: 语言设计

---

## 基础常量

```go
const Pi = 3.14159
const MaxSize = 1024

// 多常量声明
const (
    MinInt = -1 << 63
    MaxInt = 1<<63 - 1
)
```

---

## iota 枚举

```go
const (
    Sunday = iota    // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
    Thursday         // 4
    Friday           // 5
    Saturday         // 6
)
```

### iota 技巧

```go
// 位掩码
const (
    Read = 1 << iota   // 1 (001)
    Write              // 2 (010)
    Execute            // 4 (100)
)

// 跳过值
const (
    _ = iota           // 跳过 0
    KB = 1 << (10 * iota)  // 1024
    MB = 1 << (10 * iota)  // 1048576
    GB = 1 << (10 * iota)  // 1073741824
)

// 复杂表达式
const (
    StatusOK = iota
    StatusError
    StatusPending

    StatusMax = iota  // 3，统计数量
)
```

---

## 无类型常量

```go
const Untyped = 42  // 无类型整数

var i int = Untyped
var f float64 = Untyped
var c complex128 = Untyped

// 有类型常量
const Typed int = 42

// var f float64 = Typed  // 错误：类型不匹配
```

---

## 常量规则

```go
// ✅ 常量可以是基本类型、字符串
const (
    Num = 42
    Str = "hello"
    Bool = true
)

// ❌ 常量不能是 slice、map、func、chan
// const S = []int{1, 2, 3}  // 错误
// const M = map[string]int{}  // 错误

// ❌ 常量必须是编译期确定的
// const R = rand.Int()  // 错误
// const T = time.Now()  // 错误
```

---

## 实战应用

### HTTP 状态码

```go
const (
    StatusOK           = 200
    StatusCreated      = 201
    StatusBadRequest   = 400
    StatusUnauthorized = 401
    StatusNotFound     = 404
    StatusInternal     = 500
)
```

### 时间常量

```go
const (
    Nanosecond  = 1
    Microsecond = 1000 * Nanosecond
    Millisecond = 1000 * Microsecond
    Second      = 1000 * Millisecond
    Minute      = 60 * Second
    Hour        = 60 * Minute
)
```

### 权限枚举

```go
type Permission int

const (
    PermRead Permission = 1 << iota
    PermWrite
    PermDelete
    PermAdmin
)

func (p Permission) Has(perm Permission) bool {
    return p&perm != 0
}

func (p Permission) Add(perm Permission) Permission {
    return p | perm
}

func (p Permission) Remove(perm Permission) Permission {
    return p &^ perm
}

// 使用
var userPerms = PermRead | PermWrite
if userPerms.Has(PermWrite) {
    // 允许写入
}
```
