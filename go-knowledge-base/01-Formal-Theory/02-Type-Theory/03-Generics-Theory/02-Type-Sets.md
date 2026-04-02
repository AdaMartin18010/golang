# 类型集合 (Type Sets)

> **分类**: 形式理论

---

## 定义

Go 1.18+ 接口定义**类型集合**:

```go
type Signed interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

Signed = {int, int8, int16, int32, int64}
```

---

## 操作

### 并集

```
A | B = A ∪ B
```

### 交集

```go
type Both interface {
    Signed
    Integer
}

Both = Signed ∩ Integer
```

---

## 近似符号 (~)

```go
~int 表示 所有底层类型为 int 的类型

包括: int, MyInt (type MyInt int)
```

---

## 约束满足

```
T satisfies C  ⟺  T ∈ typeset(C)
```
