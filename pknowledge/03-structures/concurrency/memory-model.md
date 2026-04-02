# Go 内存模型

> **创建**: 2026-04-02
> **状态**: 持续更新
> **关联**: [[CSP Model]], [[DRF-SC Guarantee]]

---

## 核心原则

**Do not communicate by sharing memory; share memory by communicating.**

---

## Happens-Before 关系

### 定义

偏序关系 `happens-before` (→hb) 定义事件顺序。

### 公理

1. **初始化顺序**

   ```
   如果包 p 导入 q，则 q 的 init 先于 p 的 init
   ```

2. **Goroutine 创建**

   ```
   go stmt →hb stmt 的执行
   ```

3. **Channel 通信**

   ```
   send(c, v) →hb receive(c) → v
   close(c) →hb receive(c) → zero
   ```

4. **锁同步**

   ```
   unlock(m) →hb lock(m)
   ```

### 图示

```
Goroutine 1          Goroutine 2
    |                      |
send(ch, v) ─────────→ receive(ch)
    |                      |
   ...                    ...
    ↓                      ↓
  操作 A                  操作 B

    A →hb B (通过 channel)
```

---

## DRF-SC 保证

### 定理

**Data-Race-Free programs behave Sequentially Consistently.**

### 证明概要

1. 定义 happens-before 关系
2. 证明无数据竞争程序有全序
3. 全序等价于顺序一致性

### 意义

- 无竞争程序无需考虑内存重排序
- 简化并发编程心智模型

---

## 数据竞争

### 定义

对同一变量的两个访问，其中至少一个写入，且无 happens-before 关系。

### 检测

```bash
go run -race myprogram.go
```

### 避免

- 使用 channel 通信
- 使用 sync.Mutex
- 使用 atomic 操作

---

## 示例分析

### 正确用法

```go
var wg sync.WaitGroup
var result int

wg.Add(1)
go func() {
    defer wg.Done()
    result = compute()  // 写
}()

wg.Wait()  // Wait 内部有同步
fmt.Println(result)  // 读 - happens-after
```

### 错误用法

```go
var done bool
var result int

go func() {
    result = compute()
    done = true  // 写 - 无同步
}()

for !done {  // 读 - 数据竞争
}
fmt.Println(result)  // 可能看到旧值
```

---

## 形式化验证

### 目标

验证程序无数据竞争。

### 方法

1. **模型检测** - 探索所有执行路径
2. **类型系统** - 编译时检查
3. **运行时检测** - Race Detector

---

## 关联

- [[CSP Model]]
- [[Channel Semantics]]
- [[Mutex Implementation]]
