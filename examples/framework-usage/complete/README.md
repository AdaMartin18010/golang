# 框架完整使用示例

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 概述

本示例展示如何在一个完整的 HTTP 服务中集成和使用框架的所有核心能力。

---

## 🎯 功能展示

本示例展示了以下框架能力：

1. ✅ **可观测性（OTLP）** - OpenTelemetry 集成
2. ✅ **采样机制** - 请求采样
3. ✅ **追踪定位** - 分布式追踪和错误定位
4. ✅ **数据库抽象** - 通用数据库接口
5. ✅ **精细控制** - 功能开关、速率控制、熔断器
6. ✅ **数据转换** - 请求数据转换
7. ✅ **反射/自解释** - 元数据检查

---

## 🚀 运行示例

### 1. 启动服务

```bash
cd examples/framework-usage/complete
go run main.go
```

### 2. 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 获取用户列表
curl http://localhost:8080/api/v1/users

# 实验性功能
curl http://localhost:8080/api/v1/experimental
```

---

## 📚 相关文档

- [框架核心能力总结](../../../docs/framework/07-框架核心能力总结.md)
- [核心能力使用示例](../../../docs/framework/08-核心能力使用示例.md)
- [框架能力完整集成示例](../../../docs/framework/13-框架能力完整集成示例.md)

---

**最后更新**: 2025-01-XX
