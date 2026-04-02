# 函子 (Functors)

> **分类**: 形式理论模型

---

## 定义

函子是范畴之间的映射，保持结构。

```
F: C → D
```

---

## Go 中的函子

### Map 函子

```go
// []T 是一个函子
func Map[A, B any](fa []A, f func(A) B) []B {
    fb := make([]B, len(fa))
    for i, a := range fa {
        fb[i] = f(a)
    }
    return fb
}

// 使用
nums := []int{1, 2, 3}
doubles := Map(nums, func(n int) int { return n * 2 })
```

### Option 函子

```go
type Option[T any] struct {
    value T
    valid bool
}

func (o Option[T]) Map(f func(T) T) Option[T] {
    if !o.valid {
        return Option[T]{}
    }
    return Option[T]{value: f(o.value), valid: true}
}
```

---

## 函子定律

```go
// 1. 恒等律: map(id) == id
Map(slice, func(x int) int { return x })  // 等于原 slice

// 2. 结合律: map(f ∘ g) == map(f) ∘ map(g)
```
