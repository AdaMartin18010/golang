# LD-007: Go 反射与接口内部机制 (Go Reflection & Interface Internals)

> **维度**: Language Design
> **级别**: S (15+ KB)
> **标签**: #go-reflection #interface #type-assertion #dynamic-dispatch
> **权威来源**: [Go Data Structures: Interfaces](https://research.swtch.com/interfaces), [reflect package](https://golang.org/pkg/reflect/)

---

## 接口内部表示

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Interface Internal Representation                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  空接口 (interface{})                                                        │
│  ┌─────────────────────────────────────┐                                    │
│  │  tab *itab  (或 *struct{Type, Hash}) │  类型信息                        │
│  │  data unsafe.Pointer                 │  数据指针                        │
│  └─────────────────────────────────────┘                                    │
│                                                                              │
│  非空接口 (io.Reader)                                                        │
│  ┌─────────────────────────────────────┐                                    │
│  │  tab *itab                          │  接口表 (类型 + 方法表)           │
│  │  data unsafe.Pointer                 │  数据指针                        │
│  └─────────────────────────────────────┘                                    │
│                                                                              │
│  itab 结构:                                                                  │
│  ┌─────────────────────────────────────┐                                    │
│  │  inter *interfacetype               │  接口类型                        │
│  │  _type *_type                       │  动态类型                        │
│  │  hash uint32                        │  类型哈希 (用于 switch)          │
│  │  _ [4]byte                          │  对齐                            │
│  │  fun [1]uintptr                     │  方法表 (变长)                   │
│  │      fun[0] = Read                  │                                  │
│  └─────────────────────────────────────┘                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 反射机制

### reflect.Type 和 reflect.Value

```go
package main

import (
    "fmt"
    "reflect"
)

type Person struct {
    Name string
    Age  int
}

func (p Person) Greet() string {
    return fmt.Sprintf("Hello, I'm %s", p.Name)
}

func main() {
    p := Person{Name: "Alice", Age: 30}

    // 获取 reflect.Type
    t := reflect.TypeOf(p)
    fmt.Println(t.Name())        // "Person"
    fmt.Println(t.Kind())        // "struct"
    fmt.Println(t.NumField())    // 2

    // 获取 reflect.Value
    v := reflect.ValueOf(p)
    fmt.Println(v.Field(0))      // "Alice"

    // 遍历字段
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        fmt.Printf("%s: %v\n", field.Name, value)
    }

    // 调用方法
    method := v.MethodByName("Greet")
    result := method.Call(nil)
    fmt.Println(result[0])       // "Hello, I'm Alice"
}
```

---

## 动态类型断言

```go
// 类型断言
func process(v interface{}) {
    // 方式1：直接断言
    if s, ok := v.(string); ok {
        fmt.Println("String:", s)
        return
    }

    // 方式2：switch
    switch x := v.(type) {
    case int:
        fmt.Println("Int:", x)
    case string:
        fmt.Println("String:", x)
    case io.Reader:
        fmt.Println("Reader:", x)
    default:
        fmt.Printf("Unknown: %T\n", x)
    }
}

// 底层实现：
// 1. 比较接口的 tab._type 与目标类型
// 2. 如果匹配，返回 data 指针
// 3. 不匹配，panic 或返回 false (ok 形式)
```

---

## 反射性能

| 操作 | 直接调用 | 反射 | 开销 |
|------|---------|------|------|
| 方法调用 | 1x | 100-1000x | 高 |
| 字段访问 | 1x | 50-100x | 高 |
| 类型断言 | 1x | 1x | 低 |
| 接口赋值 | 1x | 1x | 低 |

### 优化：避免反射

```go
// 反模式：运行时反射
func Serialize(v interface{}) []byte {
    val := reflect.ValueOf(v)
    // ... 反射处理
}

// 正解：代码生成 (如 msgp, easyjson)
//go:generate msgp

// 或使用泛型 (Go 1.18+)
func Serialize[T any](v T) []byte {
    // 编译期确定类型
}
```

---

## 参考文献

1. [Go Data Structures: Interfaces](https://research.swtch.com/interfaces) - Russ Cox
2. [The Laws of Reflection](https://blog.golang.org/laws-of-reflection) - Rob Pike
3. [reflect package](https://golang.org/pkg/reflect/)
