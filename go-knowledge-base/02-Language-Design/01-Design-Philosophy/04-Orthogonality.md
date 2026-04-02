# 正交性 (Orthogonality)

> **分类**: 语言设计

---

## 定义

**正交性**: 语言特性相互独立，可自由组合。

---

## Go 的正交特性

### 1. 类型系统正交

```go
// 任何类型可以组合
// 任何类型可以实现接口
// 接口可以组合

type MyInt int
type IntSlice []int
type IntMap map[string]int

// 都独立工作
```

### 2. 控制流正交

```go
// for 可用于所有迭代
for i := 0; i < 10; i++ { }
for k, v := range m { }
for v := range ch { }

// 无单独的 while/do-while
```

### 3. 并发正交

```go
// goroutine 可与任何函数组合
go anyFunction()

// channel 可与任何类型组合
ch := make(chan any)
```

---

## 非正交的反例

| 语言 | 非正交设计 |
|------|-----------|
| C++ | 类内函数重载规则复杂 |
| Java | 基本类型 vs 包装类型 |
| Python | 特殊方法 vs 普通方法 |

---

## Go 的正交收益

1. **学习成本低**: 掌握基本规则后可推导
2. **代码一致**: 无特殊情况处理
3. **工具友好**: 易于静态分析

---

## 示例

```go
// 接口 + 泛型 + 并发 正交组合
func ProcessAll[T any](items []T, process func(T)) {
    var wg sync.WaitGroup
    for _, item := range items {
        wg.Add(1)
        go func(i T) {
            defer wg.Done()
            process(i)
        }(item)
    }
    wg.Wait()
}
```
