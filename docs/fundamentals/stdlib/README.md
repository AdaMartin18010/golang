# Go标准库完整指南

**章节定位**: Go语言基础 > 标准库  
**难度级别**: 中级  
**预计学习时间**: 12-15小时

---

## 📖 章节概述

Go标准库是Go语言的核心组成部分，提供了丰富且高质量的包，涵盖了从基础I/O、网络通信到并发控制、加密解密等各个方面。掌握标准库是成为Go专家的必经之路。

### 🎯 学习目标

完成本章学习后，你将能够：

- ✅ 熟练使用核心包处理常见任务
- ✅ 掌握字符串和文本处理
- ✅ 精通文件和I/O操作
- ✅ 构建网络应用和HTTP服务
- ✅ 使用并发包实现同步
- ✅ 处理各种编码和解码
- ✅ 编写高质量的测试代码

---

## 📚 章节内容

### [01-核心包概览](./01-核心包概览.md)
**难度**: ⭐⭐  
**预计阅读**: 20分钟

- fmt - 格式化输入输出
- errors - 错误处理
- io - I/O原语
- os - 操作系统功能
- time - 时间处理
- flag - 命令行参数解析

**关键概念**: 接口、错误链、Reader/Writer、时间格式化

---

### [02-字符串处理](./02-字符串处理.md)
**难度**: ⭐⭐  
**预计阅读**: 20分钟

- strings - 字符串操作
- strconv - 字符串转换
- unicode - Unicode处理
- unicode/utf8 - UTF-8编码
- regexp - 正则表达式
- text/template - 文本模板

**关键概念**: 字符串不可变性、rune、编码转换、正则匹配

---

### [03-文件与IO](./03-文件与IO.md)
**难度**: ⭐⭐⭐  
**预计阅读**: 25分钟

- os - 文件操作
- io - I/O接口
- io/ioutil - I/O工具（已废弃，迁移到os和io）
- bufio - 缓冲I/O
- path/filepath - 文件路径操作
- embed - 嵌入文件

**关键概念**: File、Reader/Writer、缓冲区、路径处理

---

### [04-网络编程](./04-网络编程.md)
**难度**: ⭐⭐⭐  
**预计阅读**: 30分钟

- net - 网络I/O
- net/http - HTTP客户端和服务器
- net/url - URL解析
- net/rpc - RPC框架
- crypto/tls - TLS/SSL

**关键概念**: TCP/UDP、HTTP、TLS、连接管理、超时控制

---

### [05-并发包](./05-并发包.md)
**难度**: ⭐⭐⭐  
**预计阅读**: 25分钟

- sync - 同步原语
- sync/atomic - 原子操作
- context - 上下文管理
- runtime - 运行时接口

**关键概念**: Mutex、WaitGroup、Once、原子操作、调度器

---

### [06-编码解码](./06-编码解码.md)
**难度**: ⭐⭐⭐  
**预计阅读**: 25分钟

- encoding/json - JSON编解码
- encoding/xml - XML编解码
- encoding/base64 - Base64编码
- encoding/hex - 十六进制编码
- encoding/gob - Go专用编码
- compress/* - 压缩算法

**关键概念**: Marshal/Unmarshal、结构体标签、自定义编解码

---

### [07-测试包](./07-测试包.md)
**难度**: ⭐⭐⭐  
**预计阅读**: 25分钟

- testing - 测试框架
- testing/iotest - I/O测试工具
- testing/quick - 快速检查
- net/http/httptest - HTTP测试
- testing/fstest - 文件系统测试

**关键概念**: 单元测试、基准测试、Mock、测试覆盖率

---

### [08-其他重要包](./08-其他重要包.md)
**难度**: ⭐⭐  
**预计阅读**: 20分钟

- log - 日志记录
- log/slog - 结构化日志
- crypto/* - 加密算法
- hash/* - 哈希算法
- math - 数学函数
- sort - 排序算法
- container/* - 容器数据结构

**关键概念**: 日志级别、加密安全、哈希校验、排序接口

---

## 🎓 学习路径

### 初学者路线
```
核心包概览 → 字符串处理 → 文件与IO → 测试包
```
**目标**: 掌握基础标准库，能够编写简单的Go程序

### 进阶路线
```
网络编程 → 并发包 → 编码解码 → 其他重要包
```
**目标**: 精通标准库，能够构建复杂的生产级应用

### 推荐学习顺序

1. **第一周**: 核心包和字符串处理（3小时）
   - 理解fmt、errors、io等基础包
   - 掌握字符串操作和转换

2. **第二周**: 文件I/O和网络编程（5小时）
   - 掌握文件操作
   - 学习HTTP客户端和服务器

3. **第三周**: 并发和编码（4小时）
   - 深入sync包
   - 精通JSON等编解码

4. **第四周**: 测试和其他包（3小时）
   - 编写高质量测试
   - 学习加密、日志等

---

## 🔗 相关章节

### 前置知识
- [语法基础](../language/01-语法基础/) - Go基础语法
- [并发编程](../language/02-并发编程/) - 并发模型

### 后续进阶
- [Web开发](../../development/web/) - Web应用开发
- [微服务](../../advanced/architecture-practices/microservice.md) - 微服务架构
- [性能优化](../../advanced/performance/) - 性能分析和优化

---

## 💻 实战项目

完成学习后，建议实践以下项目：

1. **文件处理工具** - 使用os、io、filepath包
2. **HTTP API服务** - 使用net/http包
3. **命令行工具** - 使用flag、os包
4. **数据转换工具** - 使用encoding包
5. **并发下载器** - 综合应用sync、net/http

---

## 📊 标准库分类

### 按功能分类

```
核心功能
├── fmt          // 格式化
├── errors       // 错误处理
├── io           // I/O原语
└── os           // 系统功能

文本处理
├── strings      // 字符串
├── strconv      // 转换
├── unicode      // Unicode
├── regexp       // 正则
└── text/*       // 模板

文件系统
├── os           // 文件操作
├── io           // I/O接口
├── bufio        // 缓冲I/O
└── path/*       // 路径

网络通信
├── net          // 网络I/O
├── net/http     // HTTP
├── net/url      // URL
└── crypto/tls   // TLS

并发控制
├── sync         // 同步
├── sync/atomic  // 原子
├── context      // 上下文
└── runtime      // 运行时

数据编码
├── encoding/*   // 各种编码
├── compress/*   // 压缩
└── crypto/*     // 加密

开发工具
├── testing      // 测试
├── log          // 日志
└── flag         // 参数
```

---

## ⚠️ 常见问题

### Q1: 标准库包太多，如何系统学习？
**建议**:
1. 先学核心包（fmt、errors、io、os）
2. 根据项目需求按需学习
3. 多看官方文档和示例
4. 阅读优秀开源项目的代码

### Q2: 标准库vs第三方库如何选择？
**原则**:
- 优先使用标准库（稳定、无依赖）
- 标准库不满足时再考虑第三方库
- 关注社区活跃度和维护状态

### Q3: 如何查找标准库文档？
**资源**:
- 官方文档：https://pkg.go.dev/std
- Go by Example：https://gobyexample.com
- 源码：https://github.com/golang/go

### Q4: 标准库版本兼容性如何？
- Go 1.x版本保证API兼容
- 新版本可能添加新功能
- 关注Deprecated标记
- 及时更新到新API

---

## 📚 推荐资源

### 官方文档
- [Package Documentation](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go)

### 推荐书籍
- 《The Go Programming Language》
- 《Go语言标准库》

### 在线资源
- [Go by Example](https://gobyexample.com)
- [Go Standard Library Examples](https://github.com/SimonWaldherr/golang-examples)

---

## 🎯 下一步

完成本章学习后，你可以：

1. **深入Web开发**: 学习 [Web开发](../../development/web/)
2. **并发编程**: 学习 [并发编程高级](../language/02-并发编程/)
3. **性能优化**: 学习 [性能分析](../../advanced/performance/)
4. **架构设计**: 学习 [设计模式](../../advanced/architecture/)

---

**维护者**: Documentation Team  
**最后更新**: 2025-10-27  
**版本**: v1.0

