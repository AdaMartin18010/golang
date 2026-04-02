# 类型断言 (Type Assertions)

> **分类**: 语言设计

---

## 基本断言

```go
var i interface{} = "hello"

// 断言为具体类型
s := i.(string)  // "hello"

// 带检查的断言
s, ok := i.(string)  // ok = true
n, ok := i.(int)     // ok = false

// 安全断言
if s, ok := i.(string); ok {
    fmt.Println(s)
} else {
    fmt.Println("not a string")
}
```

---

## 类型开关

```go
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("int: %d\n", v)
    case string:
        fmt.Printf("string: %q\n", v)
    case bool:
        fmt.Printf("bool: %t\n", v)
    case []int:
        fmt.Printf("slice of ints: %v\n", v)
    case map[string]int:
        fmt.Printf("map: %v\n", v)
    case nil:
        fmt.Println("nil")
    case Person:
        fmt.Printf("person: %s\n", v.Name)
    case *Person:
        fmt.Printf("person pointer: %s\n", v.Name)
    default:
        fmt.Printf("unknown type: %T\n", v)
    }
}
```

---

## 接口断言

```go
type Reader interface {
    Read() ([]byte, error)
}

type Closer interface {
    Close() error
}

func process(r Reader) error {
    // 检查是否也实现了 Closer
    if c, ok := r.(Closer); ok {
        defer c.Close()
    }

    data, err := r.Read()
    // ...
}
```

---

## 空接口判断

```go
func isNil(i interface{}) bool {
    if i == nil {
        return true
    }

    // 检查底层值是否为 nil
    v := reflect.ValueOf(i)
    return v.Kind() == reflect.Ptr && v.IsNil()
}

// 使用
var p *Person
fmt.Println(isNil(p))  // true
fmt.Println(isNil(nil)) // true
```

---

## 性能考虑

```go
// 编译器优化：直接类型断言
var r io.Reader = strings.NewReader("hello")

// 编译器已知类型，优化为直接检查
if sr, ok := r.(*strings.Reader); ok {
    // 无运行时开销
}

// 多次断言缓存类型
var typ = reflect.TypeOf((*MyInterface)(nil)).Elem()

func checkType(i interface{}) bool {
    return reflect.TypeOf(i).Implements(typ)
}
```

---

## 最佳实践

```go
// ✅ 优先使用类型开关
switch v := i.(type) {
case int:
    // ...
case string:
    // ...
}

// ✅ 检查 ok 值避免 panic
v, ok := i.(MyType)
if !ok {
    return errors.New("type mismatch")
}

// ❌ 避免裸断言
v := i.(MyType)  // 可能 panic

// ✅ 断言前检查 nil
if i == nil {
    return nil
}
v, ok := i.(MyType)
```
