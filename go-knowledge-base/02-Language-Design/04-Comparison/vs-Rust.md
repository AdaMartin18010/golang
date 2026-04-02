# Go vs Rust 对比

> **分类**: 语言设计

---

## 概览

| 特性 | Go | Rust |
|------|-----|------|
| 发布时间 | 2009 | 2010 |
| 内存安全 | GC | 所有权系统 |
| 并发 | Goroutine + Channel | async/await |
| 编译速度 | 快 | 慢 |
| 学习曲线 | 平缓 | 陡峭 |
| 性能 | 好 | 极致 |
| 适用场景 | 网络服务、云原生 | 系统编程、嵌入式 |

---

## 内存管理

### Go: 垃圾回收

```go
func process() {
    data := make([]byte, 1024)
    // 使用 data
    // 自动释放
}
```

**优点**: 简单、安全
**缺点**: GC 停顿、内存开销

### Rust: 所有权

```rust
fn process() {
    let data = vec![0u8; 1024];
    // 使用 data
    // 编译时确定生命周期
} // 自动释放
```

**优点**: 零成本抽象、无 GC
**缺点**: 学习曲线陡峭、编译器严格

---

## 并发模型

### Go: CSP

```go
ch := make(chan int)

go func() {
    ch <- 42
}()

v := <-ch
```

### Rust: async/await

```rust
let (tx, rx) = tokio::sync::mpsc::channel(100);

tokio::spawn(async move {
    tx.send(42).await.unwrap();
});

let v = rx.recv().await.unwrap();
```

---

## 错误处理

### Go: 显式检查

```go
f, err := os.Open("file")
if err != nil {
    return err
}
defer f.Close()
```

### Rust: Result 类型

```rust
let f = File::open("file")?;
// 自动传播错误
```

---

## 何时选择

### 选择 Go

- 网络服务
- 微服务
- 云原生应用
- 快速开发
- 大型团队协作

### 选择 Rust

- 系统编程
- 嵌入式
- 性能关键
- 安全关键
- WebAssembly

---

## 总结

Go 和 Rust 都是现代系统语言，但设计目标不同：

- **Go**: 简单、高效、适合大规模工程
- **Rust**: 安全、高性能、零成本抽象
