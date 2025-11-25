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
- **消息队列**: Kafka, MQTT
- **API**: REST, gRPC, GraphQL

---

## 📚 文档

- [架构文档](docs/architecture/)
- [开发指南](docs/guides/)
- [API 文档](docs/api/)

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
