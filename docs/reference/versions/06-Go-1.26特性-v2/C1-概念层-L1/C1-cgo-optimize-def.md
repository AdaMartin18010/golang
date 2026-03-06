# cgo优化 (CGO Optimizations)

> **文档层级**: C1-概念层 (Concept Layer L1)
> **文档类型**: 概念定义 (Concept Definition)
> **最后更新**: 2026-03-06

---

## 一、概念定义

### 1.1 cgo概述

```
cgo: Go与C代码互操作的机制

Go 1.26优化:
  - 更快的cgo调用
  - 优化的C字符串处理
  - 改进的内存管理
```

### 1.2 性能提升

| 优化项 | 改进幅度 | 说明 |
|--------|----------|------|
| 调用开销 | 30-50% | 减少上下文切换 |
| C字符串转换 | 20-40% | 优化C.GoString等 |
| 内存分配 | 15-25% | 减少临时分配 |

---

## 二、核心优化

### 2.1 调用路径优化

```go
// Go 1.25: 每次cgo调用都有较重开销
// Go 1.26: 优化后的轻量级调用

// 批处理模式减少调用次数
func BatchProcess(items []Item) {
    // 一次性传递多个项目，减少cgo调用次数
    C.process_batch(
        unsafe.Pointer(&items[0]),
        C.size_t(len(items)),
    )
}
```

### 2.2 字符串处理优化

```go
// 优化的C字符串转换
func optimizedString() {
    cs := C.CString("hello")  // 分配C字符串
    defer C.free(unsafe.Pointer(cs))

    // Go 1.26优化了C.GoString性能
    gs := C.GoString(cs)  // 更快的转换
}
```

---

## 三、最佳实践

### 3.1 减少调用开销

```go
// ✅ 批处理
func processBatch(data []byte) {
    C.process_large_buffer(
        (*C.char)(unsafe.Pointer(&data[0])),
        C.size_t(len(data)),
    )
}

// ❌ 避免频繁小调用
func processSlow(items []int) {
    for _, item := range items {
        C.process_one(C.int(item))  // 每次都有开销
    }
}
```

### 3.2 内存管理

```go
// ✅ 使用C内存池
var bufferPool = sync.Pool{
    New: func() interface{} {
        return C.malloc(1024)
    },
}

// ✅ 避免重复分配
func processWithBuffer() {
    buf := bufferPool.Get().(unsafe.Pointer)
    defer bufferPool.Put(buf)

    C.process(buf)
}
```

---

**概念分类**: 运行时 - FFI优化
**Go版本**: 1.26 (自动启用)
