# Happens-Before 关系

> **分类**: 形式理论
> **难度**: 专家
> **前置知识**: 并发基础、偏序关系

---

## 概述

**Happens-Before** (→hb) 是定义程序事件偏序关系的核心概念，用于确定内存操作的可见性。

---

## 形式化定义

### 关系定义

```
→hb ⊆ Event × Event

是程序事件上的偏序关系（自反、反对称、传递）
```

### 事件类型

```
Event ::=
  | Read(var, val)     读取变量
  | Write(var, val)    写入变量
  | Send(ch, val)      发送值到通道
  | Recv(ch, val)      从通道接收值
  | Lock(m)            获取锁
  | Unlock(m)          释放锁
  | Spawn(g)           创建 goroutine
  | Join(g)            等待 goroutine
```

---

## Go 内存模型的 Happens-Before 规则

### 1. 程序顺序

```
在同一线程/goroutine中:

e₁ 在程序顺序中先于 e₂
─────────────────────────── (HB-Program)
e₁ →hb e₂
```

**Go 示例**:

```go
x = 1      // e₁
y = 2      // e₂
// e₁ →hb e₂
```

### 2. 初始化顺序

```
包 p 导入包 q
─────────────────────────── (HB-Init)
init(q) →hb init(p)
```

### 3. Goroutine 创建

```
─────────────────────────────── (HB-Spawn)
go f() →hb f()的第一条语句
```

**Go 示例**:

```go
x = 1         // e₁
go func() {   // e₂: spawn
    print(x)  // e₃
}()
// e₁ →hb e₂ →hb e₃
// 因此 print(x) 一定看到 x=1
```

### 4. Channel 发送-接收

```
Send(ch, v) →hb Recv(ch, v)    (HB-Channel)

注意: 同一值的匹配发送和接收
```

**Go 示例**:

```go
// goroutine 1
ch <- 42     // e₁

// goroutine 2
v := <-ch    // e₂
// e₁ →hb e₂, 因此 v = 42
```

### 5. Channel 关闭

```
Close(ch) →hb Recv(ch, zero)    (HB-Close)
```

### 6. 有缓冲 Channel (容量 k)

```
第 i 次接收 →hb 第 (i+k) 次发送    (HB-Buffered)

解释: 接收释放缓冲槽，允许后续发送
```

**示例** (缓冲为 2):

```
send₁ →hb send₃  (因为 recv₁ 释放槽位)
send₂ →hb send₄
recv₁ →hb send₃
recv₂ →hb send₄
```

### 7. Lock/Unlock

```
Unlock(m) →hb Lock(m)    (HB-Lock)

同一互斥锁的释放先于后续获取
```

**Go 示例**:

```go
// goroutine 1
mu.Lock()
x = 1
mu.Unlock()    // e₁

// goroutine 2
mu.Lock()      // e₂
print(x)
// e₁ →hb e₂, 因此看到 x=1
```

### 8. WaitGroup

```
WaitGroup.Done() →hb WaitGroup.Wait()返回    (HB-WaitGroup)
```

### 9. Once

```
once.Do(f)中的操作 →hb once.Do返回    (HB-Once)
```

---

## Happens-Before 图

### 示例程序

```go
var x, y int

func main() {
    x = 1                    // e1
    go func() {              // e2: spawn
        print(x)             // e3
        y = 1                // e4
        done <- true         // e5
    }()
    <-done                   // e6: recv
    print(y)                 // e7
}
```

### HB 图

```
e1: x = 1
 ↓ (程序顺序)
e2: go func()
 ↓ (spawn)
e3: print(x)
 ↓ (程序顺序)
e4: y = 1
 ↓ (程序顺序)
e5: done <- true
 ↓ (channel)
e6: <-done
 ↓ (程序顺序)
e7: print(y)
```

### 推导

```
e1 →hb e2 →hb e3  ⇒  print(x) 看到 x=1
e4 →hb e5 →hb e6 →hb e7  ⇒  print(y) 看到 y=1
```

---

## 同步关系 (Synchronizes-With)

```
同步关系是 Happens-Before 的子集，由同步原语创建:

- 发送 → 接收
- Unlock → Lock
- Spawn → 第一条语句
- 等等
```

---

## 可见性

### 读可见写

```
读 r 可以看到写 w，如果:
1. w →hb r, 且
2. 没有写 w' 使得 w →hb w' →hb r

即: r 之前的最后一次写
```

### 数据竞争

```
两个访问构成数据竞争，如果:
1. 访问同一变量
2. 至少一个是写
3. 它们之间没有 →hb 关系
```

---

## 形式化证明

### 传递性

```
Theorem: →hb 是传递的

证明:
如果 e₁ →hb e₂ 且 e₂ →hb e₃
根据定义，→hb 是传递闭包
因此 e₁ →hb e₃
```

### DRF 保证基础

```
Theorem: 没有数据竞争的程序，所有读都看到一致的值

证明概要:
- 无数据竞争 ⇒ 所有冲突访问有 →hb 关系
- →hb 是全序 ⇒ 读看到唯一确定的写
```

---

## 常见模式

### 模式 1: 通过 Channel 传递所有权

```go
ch := make(chan *Data)
// goroutine A
data := &Data{...}
ch <- data        // e1: 发送

// goroutine B
data := <-ch      // e2: 接收
data.Use()        // e3: 使用
// e1 →hb e2 →hb e3，安全访问
```

### 模式 2: 通过 Mutex 保护共享数据

```go
mu.Lock()         // e1
x = 42            // e2
mu.Unlock()       // e3
// ...
mu.Lock()         // e4 (e3 →hb e4)
print(x)          // e5，看到 42
```

---

## 参考

- "The Go Memory Model" (go.dev/ref/mem)
- "Java Memory Model" (JSR-133)
- "C++ Memory Model" (ISO/IEC 14882:2011)
- "A Tutorial Introduction to CSP"
