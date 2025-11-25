# 🚀 Go 现代化架构项目

**版本**: Go 1.25.3
**架构**: Clean Architecture
**技术栈**: Chi, Ent, Viper, Slog, Wire, OpenTelemetry, GraphQL, gRPC, MQTT, Kafka, PostgreSQL

---

## 📁 项目结构

```text
golang/
├── cmd/                    # 主程序入口
│   ├── server/            # HTTP 服务器
│   ├── grpc-server/       # gRPC 服务器
│   ├── graphql-server/    # GraphQL 服务器
│   ├── mqtt-client/       # MQTT 客户端
│   └── cli/               # CLI 工具
│
├── internal/               # Clean Architecture
│   ├── domain/            # 领域层
│   ├── application/       # 应用层
│   ├── infrastructure/    # 基础设施层
│   ├── interfaces/        # 接口层
│   └── config/            # 配置管理
│
├── pkg/                   # 公共库
├── api/                   # API 定义
├── configs/               # 配置文件
├── scripts/               # 脚本
├── test/                  # 测试
└── docs/                  # 文档
```

---

## 🚀 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 生成代码

```bash
# Ent 代码
go generate ./internal/infrastructure/database/ent/...

# gRPC 代码
go generate ./scripts/generate/

# Wire 代码
go generate ./scripts/wire/
```

### 3. 运行应用

```bash
# HTTP 服务器
go run ./cmd/server
# 或使用 Makefile
make run

# Temporal Worker（工作流执行器）
go run ./cmd/temporal-worker
# 或使用 Makefile
make run-worker

# 使用 Docker Compose 启动所有服务（包括 Temporal）
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

## 🛠️ 技术栈

- **Web框架**: Chi
- **ORM**: Ent
- **配置**: Viper
- **日志**: Slog (Go 1.21+)
- **依赖注入**: Wire
- **数据库**: PostgreSQL (pgx)
- **可观测性**: OpenTelemetry (OTLP)
- **工作流编排**: Temporal
- **消息队列**: Kafka, MQTT
- **API**: REST, gRPC, GraphQL

---

## 📚 文档

### 核心文档

- 🏗️ **[架构文档](docs/architecture/README.md)** - Clean Architecture、领域模型、工作流架构
- 📖 **[使用指南](docs/guides/)** - 开发、部署、测试指南
- 🔧 **[API 文档](docs/api/)** - REST、GraphQL、gRPC API 规范

### 导航文档

- 📊 **[项目文档索引](docs/00-项目文档索引.md)** - 完整文档索引
- 📋 **[文档结构规范](docs/00-项目文档结构规范.md)** - 文档格式规范
- 🗺️ **[架构知识图谱](docs/architecture/00-知识图谱.md)** - 架构全景图
- 🔍 **[概念定义体系](docs/architecture/00-概念定义体系.md)** - 概念定义
- 📖 **[技术对比矩阵](docs/architecture/00-对比矩阵.md)** - 技术选型对比

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

## 📝 License

MIT
