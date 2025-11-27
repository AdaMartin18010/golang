# 1.2.1.1 Go 1.21 标准库增强：结构化日志 `slog` 深度解析

<!-- TOC START -->
- [1.2.1.1 Go 1.21 标准库增强：结构化日志 `slog` 深度解析](#1211-go-121-标准库增强结构化日志-slog-深度解析)
  - [1.2.1.1.1 🎯 **核心概念：什么是结构化日志？**](#12111--核心概念什么是结构化日志)
  - [1.2.1.1.2 ✨ **`slog` 包的核心组件**](#12112--slog-包的核心组件)
  - [1.2.1.1.3 💡 **`slog` 相比 `log` 和第三方库的优势**](#12113--slog-相比-log-和第三方库的优势)
  - [1.2.1.1.4 📝 **基本用法**](#12114--基本用法)
  - [1.2.1.1.5 🚀 **总结**](#12115--总结)
<!-- TOC END -->

## 1.2.1.1.1 🎯 **核心概念：什么是结构化日志？**

传统的日志（如标准库的 `log` 包）通常是无格式的纯文本字符串。而**结构化日志**则将日志信息记录为一种机器友好的格式，最常见的是 **JSON**。每条日志记录都包含固定的字段（如时间戳 `time`、日志级别 `level`、消息 `msg`）以及一组自定义的键值对属性（Attributes）。

**为什么需要结构化日志？**

1. **机器可读性**: 日志可以被日志聚合系统（如 ELK Stack, Splunk, Datadog）轻松地解析、索引和查询。
2. **强大的查询能力**: 可以对特定字段进行过滤、聚合和分析。例如，"查询所有 `level` 为 `error` 且 `user_id` 为 `123` 的日志"。
3. **标准化**: 统一团队和项目中的日志格式，降低维护和分析成本。
4. **上下文信息丰富**: 能够方便地将请求ID、用户信息等上下文信息附加到每条日志中。

## 1.2.1.1.2 ✨ **`slog` 包的核心组件**

`slog` 的设计是模块化和可扩展的，主要由以下几个核心组件构成：

1. **`Logger`**:
    - 日志记录器实例，是与 `slog` 交互的主要入口。它提供了 `Info()`, `Debug()`, `Warn()`, `Error()` 等方法来记录不同级别的日志。
    - 可以通过 `With()` 方法创建附带固定属性（如 `request_id`）的子 Logger。

2. **`Handler`**:
    - 处理器，是 `slog` 的后端。它负责接收 `Logger` 生成的日志记录（`Record`），并决定如何格式化（如转为 JSON 或纯文本）和输出（如写入 `os.Stdout` 或文件）。
    - `slog` 内置了两个 `Handler`:
        - `slog.TextHandler`: 输出人类可读的 `key=value` 格式文本。
        - `slog.JSONHandler`: 输出 JSON 格式文本。
    - 开发者可以实现自己的 `Handler` 来对接不同的日志系统或实现自定义格式。

3. **`Record`**:
    - 代表一条完整的日志事件的结构体。它包含了时间、级别、消息和所有的键值对属性。`Logger` 创建 `Record`，然后将其传递给 `Handler` 处理。

4. **`Level`**:
    - 日志级别（如 `Debug`, `Info`, `Warn`, `Error`）。`Logger` 和 `Handler` 都可以配置一个最低日志级别，低于该级别的日志将被忽略，从而实现对日志输出量的控制。

5. **`Attr` (Attribute)**:
    - 代表一个键值对（Key-Value Pair），是结构化日志的核心。`slog` 提供了 `slog.String()`, `slog.Int()`, `slog.Duration()` 等便捷的函数来创建不同类型的 `Attr`。

## 1.2.1.1.3 💡 **`slog` 相比 `log` 和第三方库的优势**

- **官方标准**: 作为标准库的一部分，`slog` 为 Go 生态系统提供了一个统一的日志接口，有助于库和应用之间的解耦。库可以只依赖 `slog` 接口，而最终应用可以选择使用哪个 `Handler` 实现。
- **高性能**: `slog` 的设计非常注重性能，尤其是在日志级别被禁用时，其开销极小。`JSONHandler` 的性能也与主流的第三方库（如 `zerolog`）相当。
- **易用性与灵活性**: 提供了简洁的顶层函数（如 `slog.Info()`）和面向对象的 `Logger` 实例，同时通过可插拔的 `Handler` 机制保证了极高的灵活性。
- **上下文传递**: `With()` 方法使得在日志中添加和传递上下文信息（如追踪ID、用户ID）变得极其简单和高效。

## 1.2.1.1.4 📝 **基本用法**

```go
import (
    "log/slog"
    "os"
)

func main() {
    // 1. 创建一个 JSON Handler，将日志写入标准输出
    handler := slog.NewJSONHandler(os.Stdout, nil)

    // 2. 将 Handler 设置为默认的 Logger
    slog.SetDefault(slog.New(handler))

    // 3. 记录日志
    slog.Info("User logged in", "user_id", 123, "status", "success")
    slog.Error("Failed to process request", "request_id", "req-789", "error", "database connection failed")

    // 4. 创建一个带有固定上下文的子 Logger
    requestLogger := slog.With("request_id", "req-abc-123")
    requestLogger.Info("Processing started")
    requestLogger.Warn("Cache miss", "key", "user:profile:123")
}

```

## 1.2.1.1.5 🚀 **总结**

`slog` 是 Go 1.21 中最重要的标准库新增功能之一。它不仅填补了标准库在日志功能上的空白，更提供了一个高性能、高灵活性且符合现代云原生应用需求的结构化日志标准。对于所有新的 Go 项目，`slog` 都应成为默认的日志解决方案。
