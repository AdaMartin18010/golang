# 接口内部实现 (Interface Internals)

> **分类**: 语言设计
> **标签**: #interface #runtime #internals

---

## 接口结构

```go
// 空接口 (eface)
type eface struct {
    _type *_type          // 类型信息
    data  unsafe.Pointer  // 数据指针
}

// 非空接口 (iface)
type iface struct {
    tab  *itab           // 接口表
    data unsafe.Pointer  // 数据指针
}
```

---

## itab 结构

```go
type itab struct {
    inter *interfacetype  // 接口类型
    _type *_type          // 具体类型
    hash  uint32          // 类型哈希
    _     [4]byte         // 填充
    fun   [1]uintptr      // 方法表 (变长)
}
```

### 方法表布局

```
itab.fun[0] = Type.Method0
itab.fun[1] = Type.Method1
...
```

---

## 类型断言优化

```go
// 编译器优化：直接类型断言
var r io.Reader = strings.NewReader("hello")

// 编译器已知类型，直接检查
if sr, ok := r.(*strings.Reader); ok {
    // 使用 sr
}
```

### 断言实现

```go
func assertE2I(inter *interfacetype, t *_type) *itab {
    // 1. 检查空接口
    if t == nil {
        return nil
    }

    // 2. 从缓存查找
    if m := itabTable.find(inter, t); m != nil {
        return m
    }

    // 3. 生成 itab
    m := getitab(inter, t, true)
    itabTable.add(m)

    return m
}
```

---

## 接口转换成本

| 操作 | 成本 | 说明 |
|------|------|------|
| 具体 → 接口 | 低 | 分配 itab + data |
| 接口 → 具体 | 中 | 类型断言检查 |
| 接口 → 接口 | 中 | 可能需要新 itab |

---

## 优化建议

```go
// ✅ 避免频繁装箱
func process(items []int) {
    for _, item := range items {
        fmt.Println(item)  // Println 接收 interface{}，每个 int 都装箱
    }
}

// ✅ 使用类型特定函数
func process(items []int) {
    var b strings.Builder
    for _, item := range items {
        b.WriteString(strconv.Itoa(item))
    }
    fmt.Println(b.String())
}

// ✅ 批量处理减少装箱
func processBatch(items []int, handler func([]int)) {
    handler(items)
}
```

---

## 调试技巧

```go
// 查看接口内部
func inspectInterface(i interface{}) {
    e := (*eface)(unsafe.Pointer(&i))
    fmt.Printf("Type: %v\n", e._type)
    fmt.Printf("Data: %p\n", e.data)
}

// 使用反射
func interfaceInfo(i interface{}) {
    t := reflect.TypeOf(i)
    v := reflect.ValueOf(i)
    fmt.Printf("Type: %v\n", t)
    fmt.Printf("Kind: %v\n", t.Kind())
    fmt.Printf("Value: %v\n", v)
}
```
