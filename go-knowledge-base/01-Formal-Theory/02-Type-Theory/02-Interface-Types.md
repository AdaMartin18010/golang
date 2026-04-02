# 接口类型理论

> **分类**: 形式理论

---

## 接口作为类型集合

```
type Reader interface {
    Read(p []byte) (n int, err error)
}

Reader = { T | Read(T) 有定义 }
```

---

## 方法集计算

```
methods(t_S) = { m | m 是 t_S 的方法 } ∪
               { m | t_S 嵌入 t', m ∈ methods(t') }
```

---

## 实现关系

```
t_S implements t_I  ⟺  methods(t_I) ⊆ methods(t_S)
```

---

## 接口组合

```
type ReadWriter interface {
    Reader
    Writer
}

methods(ReadWriter) = methods(Reader) ∪ methods(Writer)
```

---

## 空接口

```
type empty interface{}

所有类型都实现 empty
```
