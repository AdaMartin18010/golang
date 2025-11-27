# 1.2.3.1 Go 标准库增强：并发原语与新模式

<!-- TOC START -->
- [1.2.3.1 Go 标准库增强：并发原语与新模式](#1231-go-标准库增强并发原语与新模式)
  - [1.2.3.1.1 🎯 **核心思想：持续演进的并发能力**](#12311--核心思想持续演进的并发能力)
  - [1.2.3.1.2 ✨ **`context` 包的增强 (Go 1.21+)**](#12312--context-包的增强-go-121)
    - [1.2.3.1.2.1 1. `context.AfterFunc`](#123121-1-contextafterfunc)
    - [1.2.3.1.2.2 2. `context.WithoutCancel`](#123122-2-contextwithoutcancel)
  - [1.2.3.1.3 ✨ **`sync` 包的增强 (Go 1.21+)**](#12313--sync-包的增强-go-121)
    - [1.2.3.1.3.1 `sync.OnceFunc` / `sync.OnceValue`](#123131-synconcefunc--synconcevalue)
  - [1.2.3.1.4 ✨ **`atomic` 包的类型安全演进 (Go 1.19+)**](#12314--atomic-包的类型安全演进-go-119)
    - [1.2.3.1.4.1 `atomic.Int64`, `atomic.Pointer`, `atomic.Bool` 等类型](#123141-atomicint64-atomicpointer-atomicbool-等类型)
  - [1.2.3.1.5 🚀 **总结**](#12315--总结)
<!-- TOC END -->

## 1.2.3.1.1 🎯 **核心思想：持续演进的并发能力**

Go 的并发模型虽然从一开始就非常强大，但 Go 团队仍在持续对其进行优化和增强。这些增强通常不是颠覆性的，而是通过引入新的辅助函数、类型和模式，来解决开发者在实践中遇到的常见问题，减少样板代码，并提升代码的健壮性和可读性。

本节将介绍几个近年来加入到标准库中、非常实用的并发原语和功能。

## 1.2.3.1.2 ✨ **`context` 包的增强 (Go 1.21+)**

`context` 包是 Go 并发编程的基石，Go 1.21 为其带来了几个重要的新函数。

### 1.2.3.1.2.1 1. `context.AfterFunc`

- **功能**: 注册一个函数，当给定的 `context` 完成（被取消或超时）时，该函数会被自动调用。
- **解决了什么问题?**: 在此之前，如果想在 context 取消时执行清理操作，通常需要启动一个专门的 goroutine，在其中使用 `select` 监听 `ctx.Done()`。`AfterFunc` 将这个常见模式封装成了一个简单的函数调用。
- **用法**:

  ```go
  // 当 ctx 完成时，打印一条日志
  stop := context.AfterFunc(ctx, func() {
      log.Println("Context is done. Performing cleanup.")
  })

  // 如果想在 context 完成前主动停止回调，可以调用 stop()
  // stop()
  ```

### 1.2.3.1.2.2 2. `context.WithoutCancel`

- **功能**: 创建一个新的 context，它继承了父 context 的所有值（Values），但**忽略**了父 context 的取消信号。
- **解决了什么问题?**: 有时，即使一个请求的 context 被取消了（例如，客户端断开连接），我们仍然希望执行某些不能被中断的操作，比如将关键日志写入数据库、更新最终状态等。`WithoutCancel` 允许这些"收尾"工作在一个不会被意外取消的环境中运行。
- **用法**:

  ```go
  go func() {
      // 即使 reqCtx 被取消，下面的日志记录操作也会继续执行
      logCtx := context.WithoutCancel(reqCtx)
      writeFinalLog(logCtx, "Request finished.")
  }()
  ```

## 1.2.3.1.3 ✨ **`sync` 包的增强 (Go 1.21+)**

### 1.2.3.1.3.1 `sync.OnceFunc` / `sync.OnceValue`

- **功能**: `OnceFunc` 将一个无参数、无返回值的函数包装成一个确保只执行一次的新函数。`OnceValue` 类似，但用于包装返回一个值的函数。
- **解决了什么问题?**: 这是对 `sync.Once` 的一种更现代、更符合函数式编程风格的封装。它避免了需要显式创建一个 `sync.Once` 实例的样板代码。
- **用法**:

  ```go
  // 传统方式
  var once sync.Once
  func init() { once.Do(setup) }

  // OnceFunc 方式
  var init = sync.OnceFunc(setup)
  // 在需要时调用 init() 即可，它会自动处理"只执行一次"的逻辑
  ```

## 1.2.3.1.4 ✨ **`atomic` 包的类型安全演进 (Go 1.19+)**

### 1.2.3.1.4.1 `atomic.Int64`, `atomic.Pointer`, `atomic.Bool` 等类型

- **功能**: Go 1.19 引入了一系列原子类型，如 `atomic.Int64`, `atomic.Uint32`, `atomic.Pointer[T]` 等。这些类型将原子操作（如 `Load`, `Store`, `Add`, `CompareAndSwap`）封装为其方法。
- **解决了什么问题?**: 在此之前，原子操作是通过 `atomic` 包的顶层函数（如 `atomic.LoadInt64(&val)`）来完成的。这种方式存在两个问题：
    1. **类型不匹配**: 很容易将 `LoadInt64` 用在一个 `int32` 变量上，编译器无法发现这种错误。
    2. **忘记原子性**: 开发者可能会不小心直接对共享变量进行读写（如 `v = 10`），而不是使用原子操作，从而引发竞争条件。
- **新类型的优势**:
  - **类型安全**: `atomic.Int64` 类型只能调用 `Int64` 相关的原子方法，避免了类型混用。
  - **API清晰**: 将操作封装为方法（如 `val.Store(10)`），比 `atomic.StoreInt64(&val, 10)` 更清晰。
  - **强制原子性**: `atomic.Int64` 类型没有提供直接的读写方式，强制开发者必须使用其方法，从而保证了操作的原子性。
- **用法**:

  ```go
  // 旧方式
  var counter int64
  atomic.AddInt64(&counter, 1)

  // 新方式 (Go 1.19+)
  var counter atomic.Int64
  counter.Add(1)
  fmt.Println(counter.Load())
  ```

## 1.2.3.1.5 🚀 **总结**

这些看似微小的改进，实际上极大地提升了 Go 并发编程的体验。它们通过将常见的、易错的模式封装到标准库中，减少了样板代码，提高了代码的安全性和可读性。熟悉并使用这些新的并发原语是每个 Go 开发者保持代码现代化的重要一步。
