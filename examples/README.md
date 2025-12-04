# 示例代码

**更新**: 2025-12-03
**状态**: 完整

---

## 🎯 示例分类

### 1. 完整集成示例 ⭐⭐⭐⭐⭐

**[complete-integration/](./complete-integration/)** - **推荐首先查看**
- 展示所有核心功能的集成
- OpenTelemetry + eBPF + JWT + RBAC
- 生产就绪的示例
- **新增**: 2025-12-03

### 2. 可观测性示例 ⭐⭐⭐⭐⭐

**[observability/](./observability/)**
- **[app/](./observability/app/)** - 基础 OTLP 集成
- **[ebpf-monitoring/](./observability/ebpf-monitoring/)** - eBPF 监控 ✨ 新增
- **[system-monitoring/](./observability/system-monitoring/)** - 系统监控
- **[complete/](./observability/complete/)** - 完整可观测性
- **[docker-compose.yaml](./observability/docker-compose.yaml)** - 完整栈 ✨ 更新

### 3. 安全示例 ⭐⭐⭐⭐⭐

**[security/](./security/)** ✨ 新增
- **[auth-example/](./security/auth-example/)** - JWT + RBAC 完整示例

### 4. 框架使用示例 ⭐⭐⭐⭐

**[framework-usage/](./framework-usage/)**
- **[complete/](./framework-usage/complete/)** - 框架完整示例
- **[user-domain/](./framework-usage/user-domain/)** - 领域模型示例
- **[wire-example/](./framework-usage/wire-example/)** - Wire DI 示例

### 5. 现代特性示例 ⭐⭐⭐⭐

**[modern-features/](./modern-features/)**
- **[06-architecture-patterns/](./modern-features/06-architecture-patterns/)** - 架构模式
  - Clean Architecture
  - DDD
  - CQRS
  - Event Sourcing

### 6. 并发示例 ⭐⭐⭐

**[concurrency/](./concurrency/)**
- Goroutine 使用
- Channel 模式
- Context 应用
- 并发模式

### 7. 消息队列示例 ⭐⭐⭐

**[messaging/](./messaging/)**
- **[kafka/](./messaging/kafka/)** - Kafka 生产者/消费者
- **[nats/](./messaging/nats/)** - NATS 发布/订阅
- **[mqtt/](./messaging/mqtt/)** - MQTT 客户端

### 8. gRPC 示例 ⭐⭐⭐

**[grpc/](./grpc/)**
- gRPC 客户端
- gRPC 服务器
- 拦截器

### 9. Web 爬虫示例 ⭐⭐

**[web-crawler/](./web-crawler/)**
- 并发爬虫
- 限流控制

### 10. 测试框架示例 ⭐⭐⭐

**[testing-framework/](./testing-framework/)**
- 单元测试
- 集成测试
- Mock 使用

---

## 🚀 快速开始

### 推荐学习路径

#### 新手路径（5分钟）

```bash
# 1. 完整集成示例
cd complete-integration
go run main.go

# 2. 测试 API
curl http://localhost:8080/health
```

#### 进阶路径（30分钟）

```bash
# 1. 可观测性完整栈
cd observability
docker-compose up -d
cd app && go run main.go

# 2. eBPF 监控（Linux）
cd observability/ebpf-monitoring
sudo go run main.go

# 3. 安全示例
cd security/auth-example
go run main.go
```

#### 深度路径（2小时）

- 阅读所有示例代码
- 理解架构设计
- 查看相关文档
- 实践自己的应用

---

## 📊 示例统计

| 类别 | 示例数 | 更新状态 |
|------|--------|---------|
| 完整集成 | 1个 | ✨ 新增 |
| 可观测性 | 7个 | ✨ 更新 |
| 安全 | 1个 | ✨ 新增 |
| 框架使用 | 3个 | ✅ 完整 |
| 现代特性 | 多个 | ✅ 完整 |
| 并发 | 8个 | ✅ 完整 |
| 消息队列 | 3个 | ✅ 完整 |
| gRPC | 2个 | ✅ 完整 |
| 其他 | 多个 | ✅ 完整 |

---

## 🌟 今日新增/更新

### 2025-12-03 更新

1. **✨ 新增**: `complete-integration/` - 完整集成示例
2. **✨ 新增**: `security/auth-example/` - 安全认证示例
3. **✨ 新增**: `observability/ebpf-monitoring/` - eBPF 监控示例
4. **✨ 更新**: `observability/docker-compose.yaml` - 升级到最新版本
5. **✨ 更新**: `observability/README.md` - 完善文档

---

## 📚 相关文档

### 核心文档
- [项目 README](../README.md)
- [架构状态](../README-ARCHITECTURE-STATUS.md)
- [最终报告](../FINAL-REPORT-2025-12-03.md)

### 技术文档
- [eBPF 实现](../pkg/observability/ebpf/README.md)
- [安全模块](../pkg/security/README.md)
- [测试框架](../test/README.md)

### 框架文档
- [框架文档](../docs/framework/)
- [架构文档](../docs/architecture/)

---

## 💡 使用建议

### 学习顺序

1. **第1步**: 运行 `complete-integration` 了解整体
2. **第2步**: 查看 `observability` 了解监控
3. **第3步**: 查看 `security` 了解安全
4. **第4步**: 查看 `framework-usage` 了解框架使用
5. **第5步**: 根据需求查看其他示例

### 复用建议

- ✅ 所有示例都可以直接复用
- ✅ 根据需求调整配置
- ✅ 集成到自己的项目

---

**状态**: ✅ 完整
**质量**: ⭐⭐⭐⭐⭐
**更新**: 持续维护

🎯 **从 complete-integration 开始！**
