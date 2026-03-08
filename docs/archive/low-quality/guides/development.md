# 开发指南

## 开发环境设置

### 1. 安装依赖

```bash
make deps
# 或
go mod download
go mod tidy
```

### 2. 生成代码

```bash
make generate
# 或
go generate ./...
```

### 3. 运行应用

```bash
make run
# 或
go run ./cmd/server
```

## 代码规范

### 1. 格式化代码

```bash
make fmt
# 或
go fmt ./...
```

### 2. 代码检查

```bash
make lint
# 或
golangci-lint run
```

### 3. 运行测试

```bash
make test
# 或
go test ./...
```

## 开发流程

### 1. 添加新功能

1. 在 Domain 层定义实体和接口
2. 在 Application 层实现用例
3. 在 Infrastructure 层实现技术细节
4. 在 Interfaces 层添加 HTTP/gRPC 接口

### 2. 添加新 API

1. 更新 `api/openapi/openapi.yaml`
2. 在 `internal/interfaces/http/chi/handlers/` 添加处理器
3. 在 `internal/interfaces/http/chi/router.go` 注册路由

### 3. 数据库迁移

1. 更新 Ent Schema
2. 运行 `go generate` 生成代码
3. 应用会自动运行迁移

## 调试

### 1. 查看日志

应用使用结构化日志（JSON 格式），可以通过环境变量设置日志级别：

```bash
export LOG_LEVEL=debug
go run ./cmd/server
```

### 2. OpenTelemetry 追踪

确保 OpenTelemetry Collector 正在运行：

```bash
docker-compose -f deployments/docker/docker-compose.yml up -d otel-collector
```

## 常见问题

### 1. 依赖问题

如果遇到依赖问题，运行：

```bash
go mod tidy
go mod download
```

### 2. Ent 代码生成失败

确保已安装 Ent CLI：

```bash
go install entgo.io/ent/cmd/ent@latest
```

### 3. Wire 代码生成失败

确保已安装 Wire：

```bash
go install github.com/google/wire/cmd/wire@latest
```
