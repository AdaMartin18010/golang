# 🚀 Go 现代化架构项目

**版本**: Go 1.25.3
**架构**: Clean Architecture
**技术栈**: Chi, Echo, Ent, Viper, Slog, Wire, OpenTelemetry, GraphQL, gRPC, MQTT, Kafka, PostgreSQL

---

## 📁 项目结构

```
golang/
├── cmd/                    # 主程序入口
├── internal/               # 私有应用代码（Clean Architecture）
│   ├── domain/            # 领域层
│   ├── application/       # 应用层
│   ├── infrastructure/    # 基础设施层
│   ├── interfaces/        # 接口层
│   └── config/            # 配置管理
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
go mod download
```

### 2. 配置环境

复制配置文件：

```bash
cp configs/.env.example configs/.env
```

编辑 `configs/.env` 设置数据库等配置。

### 3. 运行服务器

```bash
go run cmd/server/main.go
```

---

## 🏗️ Clean Architecture

### 分层说明

1. **Domain Layer** - 领域层：核心业务逻辑，不依赖任何外部框架
2. **Application Layer** - 应用层：用例编排，协调领域对象
3. **Infrastructure Layer** - 基础设施层：技术实现细节（数据库、消息队列等）
4. **Interfaces Layer** - 接口层：外部接口适配（HTTP、gRPC、GraphQL等）

### 依赖方向

```
Interfaces → Application → Domain
     ↓            ↓
Infrastructure → Domain
```

---

## 🛠️ 技术栈

- **Web框架**: Chi, Echo
- **ORM**: Ent
- **配置**: Viper
- **日志**: Slog (Go 1.21+)
- **依赖注入**: Wire
- **数据库**: PostgreSQL (pgx)
- **可观测性**: OpenTelemetry (OTLP)
- **API**: OpenAPI, AsyncAPI, GraphQL, gRPC
- **消息队列**: Kafka, MQTT

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

# 运行集成测试
go test ./test/integration/...
```

---

## 📝 License

MIT
