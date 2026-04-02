# Go 并发语义

> **分类**: 形式理论

---

## Goroutine 语义

### 创建

```
(go e) 创建新 goroutine 执行 e

⟨go e, σ⟩ → ⟨(), σ'⟩  with new goroutine g executing e
```

### Happens-Before

```
go f() →hb f()的第一条语句
```

---

## Channel 语义

### 语法

```
e ::= make(chan T, n)  |  e₁ <- e₂  |  <-e
```

### 操作语义

#### 发送

```
⟨v₁ <- v₂, σ⟩ → ⟨(), σ[ch ↦ σ(ch) ∪ {v₂}]⟩

如果 channel 有缓冲且未满
```

#### 接收

```
⟨<-ch, σ⟩ → ⟨v, σ[ch ↦ σ(ch) \ {v}]⟩

如果 channel 有值
```

### Happens-Before

```
send(ch, v) →hb recv(ch, v)
close(ch) →hb recv(ch, zero)
```

---

## Select 语义

```
select {
case ch₁ <- v₁: e₁
case v₂ := <-ch₂: e₂
default: e₃
}
```

语义: 非确定性选择可用的 case

---

## 与 CSP 对比

| Go | CSP |
|----|-----|
| goroutine | process |
| channel | channel |
| select | external choice □ |
| go | parallel |||
