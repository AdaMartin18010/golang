# 🚀 Go 现代化架构项目

**版本**: Go 1.25.3
**架构**: Clean Architecture (标准4层)
**更新**: 2025-12-03
**评分**: 8.5/10 ⭐⭐⭐⭐⭐

## 🎯 项目定位

**Go Clean Architecture 企业级框架** - 集成最新最成熟的技术栈

### 核心特性

- ✅ **Clean Architecture** - 标准4层分层，依赖倒置
- ✅ **OTLP v1.38.0** - OpenTelemetry 完整可观测性
- ✅ **Cilium eBPF v0.20.0** - 真实的系统级监控
- ✅ **OAuth2/OIDC/RBAC/JWT** - 企业级安全
- ✅ **自我感知环境** - 容器/K8s/5大云厂商
- ✅ **完整测试框架** - testify + mocks
- ✅ **CI/CD** - GitHub Actions 完整流水线

### 技术栈

Chi, Ent, Wire, OpenTelemetry, Cilium eBPF, GraphQL, gRPC, Kafka, NATS, MQTT, PostgreSQL, Redis

**项目状态**: ✅ 核心架构完整 | 🔄 持续改进中 | 📊 [详细状态](./PROJECT-STATUS.md)

---

## 📁 项目结构

本项目遵循 [golang-standards/project-layout](https://github.com/golang-standards/project-layout) 标准布局：

```text
golang/
├── cmd/                    # 主程序入口
│   ├── server/            # HTTP 服务器
│   ├── grpc-server/       # gRPC 服务器
│   ├── graphql-server/    # GraphQL 服务器
│   ├── mqtt-client/       # MQTT 客户端
│   ├── cli/               # CLI 工具
│   └── temporal-worker/   # Temporal 工作流执行器
│
├── internal/               # 私有代码（Clean Architecture）
│   ├── domain/            # 领域层 - 核心业务逻辑
│   ├── application/       # 应用层 - 用例编排
│   ├── infrastructure/    # 基础设施层 - 技术实现
│   ├── interfaces/        # 接口层 - 外部接口适配
│   └── config/            # 配置管理
│
├── pkg/                    # 可被外部使用的公共库
│   ├── logger/            # 日志库
│   ├── errors/            # 错误处理
│   ├── validator/         # 验证器
│   ├── http3/             # HTTP/3 支持
│   ├── concurrency/       # 并发工具
│   └── observability/     # 可观测性工具
│
├── api/                    # API 定义文件
│   ├── openapi/           # OpenAPI/Swagger 定义
│   ├── graphql/           # GraphQL schema
│   └── asyncapi/          # AsyncAPI 定义
│
├── configs/                # 配置文件模板
│   └── config.yaml        # 默认配置
│
├── scripts/                # 构建、安装、分析等脚本
│   ├── build.sh
│   └── generate.sh
│
├── deployments/            # 部署配置和模板
│   ├── docker/            # Docker 配置
│   └── kubernetes/        # Kubernetes 配置
│
├── test/                   # 外部测试应用和测试数据
│   ├── unit/              # 单元测试
│   ├── integration/       # 集成测试
│   └── e2e/               # 端到端测试
│
├── docs/                   # 设计和用户文档
│   ├── architecture/      # 架构文档
│   ├── guides/            # 使用指南
│   ├── development/       # 开发文档
│   └── ...
│
├── examples/               # 应用程序和库的示例
│   ├── basic/             # 基础示例
│   ├── advanced/          # 高级示例
│   └── modern-features/   # 现代特性示例
│
├── tools/                  # 项目的支持工具
│   ├── codegen/           # 代码生成工具
│   └── formal-verifier/   # 形式化验证工具
│
├── migrations/             # 数据库迁移脚本
│   ├── postgres/          # PostgreSQL 迁移
│   └── ent/               # Ent 迁移
│
├── go.mod                  # Go 模块定义
├── go.sum                  # Go 模块校验和
├── go.work                 # Go 工作区
├── Makefile                # Make 构建脚本
├── Dockerfile              # Docker 镜像构建
└── README.md               # 项目说明
```

---

## 🚀 快速开始

### 方式 1: 框架快速开始（推荐）⭐

如果你是第一次使用框架，建议从框架快速开始指南开始：

```bash
# 1. 设置开发环境
make setup

# 2. 安装 Git hooks
make install-hooks

# 3. 运行示例
cd examples/framework-usage
go run main.go
```

📖 **详细指南**: [框架快速开始指南](docs/framework/05-快速开始指南.md) - 5 分钟快速上手

### 方式 2: 完整项目启动

```bash
# 1. 安装依赖
go mod tidy

# 2. 生成代码
make generate

# 3. 运行应用
make run

# 或使用 Docker Compose 启动所有服务
docker-compose -f deployments/docker/docker-compose.yml up -d
```

---

## 🏗️ Clean Architecture

### 分层说明

1. **Domain Layer** - 领域层：核心业务逻辑
2. **Application Layer** - 应用层：用例编排
3. **Infrastructure Layer** - 基础设施层：技术实现
4. **Interfaces Layer** - 接口层：外部接口适配

### 依赖方向

```text
Interfaces → Application → Domain
     ↓            ↓
Infrastructure → Domain
```

---

## 🛠️ 技术栈（2024最新版本）

### 核心框架

- **Go**: 1.25.3
- **Wire**: v0.6.0 (依赖注入)
- **Web框架**: Chi
- **ORM**: Ent
- **配置**: Viper
- **日志**: Slog (Go 1.21+)
- **依赖注入**: Wire
- **数据库**: PostgreSQL (pgx), SQLite3
- **可观测性**: OpenTelemetry (OTLP)
- **工作流编排**: Temporal
- **消息队列**: Kafka, MQTT, **NATS** ✅
- **API**: REST (OpenAPI), **gRPC** ✅, GraphQL, AsyncAPI

---

## 📚 文档

### 核心文档

- 📖 **[文档首页](docs/README.md)** - 完整文档导航和索引
- 🏗️ **[架构文档](docs/architecture/README.md)** - Clean Architecture、领域模型、工作流架构
- 📘 **[使用指南](docs/guides/)** - 开发、部署、测试指南
- 🔧 **[API 文档](api/README.md)** - REST、GraphQL、gRPC API 规范

### 快速导航

- 🎯 **[快速开始](docs/getting-started/quick-start-3min.md)** - 3分钟快速开始
- 📊 **[完整知识体系](docs/00-Go-1.25.3完整知识体系总览-2025.md)** - 系统化总领文档
- 📚 **[快速参考手册](docs/📚-Go-1.25.3快速参考手册-2025.md)** - 日常开发速查
- 🗺️ **[架构知识图谱](docs/architecture/00-知识图谱.md)** - 架构全景图
- 🔍 **[概念定义体系](docs/architecture/00-概念定义体系.md)** - 概念定义
- 📖 **[技术对比矩阵](docs/architecture/00-对比矩阵.md)** - 技术选型对比

### 项目结构

- 📋 **[项目结构重构计划](docs/00-项目结构重构计划.md)** - 项目结构优化方案
- 📊 **[项目文档索引](docs/00-项目文档索引.md)** - 完整文档索引
- 📋 **[文档结构规范](docs/00-项目文档结构规范.md)** - 文档格式规范

### 技术栈对标

- 🎉 **[技术栈实施完成最终总结](docs/00-技术栈实施完成最终总结.md)** - 实施完成最终总结 ⭐⭐⭐ 强烈推荐
- 🎉 **[技术栈实施完成README](docs/00-技术栈实施完成README.md)** - 实施完成快速导航 ⭐⭐⭐ 强烈推荐
- 🎉 **[技术栈实施完成最终报告](docs/00-技术栈实施完成最终报告.md)** - 最详细的完成报告 ⭐⭐⭐ 强烈推荐
- 🎉 **[技术栈实施全面完成报告](docs/00-技术栈实施全面完成报告.md)** - 实施完成总结报告 ⭐⭐ 推荐
- 🎯 **[项目重新定位与轻量级架构计划](docs/00-项目重新定位与轻量级架构计划.md)** - 项目重新定位和简化架构计划 ⭐ 推荐
- 🚀 **[技术栈实施快速开始](docs/00-技术栈实施快速开始.md)** - 快速开始实施指南 ⭐ 推荐
- 📋 **[技术栈实施总体规划](docs/00-技术栈实施总体规划.md)** - 4周实施总体规划和时间表 ⭐ 推荐
- 📝 **[技术栈实施详细方案](docs/00-技术栈实施详细方案.md)** - 详细的技术实施步骤和代码方案 ⭐ 推荐
- ✅ **[技术栈实施检查清单](docs/00-技术栈实施检查清单.md)** - 实施过程中的检查清单和验收标准 ⭐ 推荐
- 🎯 **[技术栈对标总结与建议](docs/00-技术栈对标总结与建议.md)** - 技术栈对标总结和优先级建议
- 📊 **[技术栈对标分析与改进计划](docs/00-技术栈对标分析与改进计划.md)** - 详细的技术栈分析和改进计划
- 💻 **[技术栈实施细节与代码建议](docs/00-技术栈实施细节与代码建议.md)** - 具体实施细节和代码实现建议

### 🔍 项目评价与改进计划

**综合评分**: 65/100 (C+) | **目标**: 90/100 (A)

- 📊 **[项目完整状态报告](docs/00-项目完整状态报告.md)** - 当前状态总览
- 📋 **[项目改进计划总览](docs/00-项目改进计划总览.md)** ⭐ **推荐阅读**
- 📖 **[执行摘要](docs/EXECUTIVE-SUMMARY.md)** - 总体评价和关键发现
- 📚 **[项目全面评价与改进计划](docs/CRITICAL-REVIEW-AND-IMPROVEMENT-PLAN.md)** - 详细评价
- 🗺️ **[改进路线图 - 可执行版本](docs/IMPROVEMENT-ROADMAP-EXECUTABLE.md)** - 详细路线图
- 📊 **[改进任务看板](docs/IMPROVEMENT-TASK-BOARD.md)** - 102 个任务清单
- 📚 **[改进文档索引](docs/README-IMPROVEMENT.md)** - 文档导航

**关键改进点**:

- 🔴 **安全性**: 50/100 → 90/100 (P0, 1-2个月)
- 🔴 **测试质量**: < 50% → > 80% (P0, 2-3个月)
- 🟡 **云原生**: 需要深度集成 Kubernetes (P1, 3-4个月)

---

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行单元测试
go test ./test/unit/...

# 测试覆盖率
go test -coverprofile=coverage.out ./...
```

---

## 📊 项目统计

- **核心模块**: 16个
- **中间件模块**: 7个
- **工具模块**: 46个
- **总计**: 64个模块
- **测试文件**: 64+个
- **文档文件**: 64+个
- **代码行数**: 30000+行

---

## 📝 License

MIT
