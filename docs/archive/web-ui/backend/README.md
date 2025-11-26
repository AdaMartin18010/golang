# Go Formal Verification - Web UI Backend

Web UI后端API服务器，为Go形式化验证框架提供RESTful API和WebSocket实时通信支持。

## 🎯 功能特性

### RESTful API

**分析API** (`/api/v1/analysis`):

- `POST /cfg` - 控制流图分析
- `POST /concurrency` - 并发安全分析
- `POST /types` - 类型系统分析
- `GET /history` - 分析历史记录

**并发模式API** (`/api/v1/patterns`):

- `GET /` - 列出所有并发模式
- `GET /:name` - 获取特定模式详情
- `POST /generate` - 生成模式代码

**项目管理API** (`/api/v1/projects`):

- `GET /` - 列出所有项目
- `POST /` - 创建新项目
- `GET /:id` - 获取项目详情
- `DELETE /:id` - 删除项目

### WebSocket

- `/ws` - WebSocket端点
- 实时分析进度推送
- 实时结果更新
- 双向通信支持

### 健康检查

- `/health` - 服务健康状态

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Git

### 安装依赖

```bash
cd web-ui/backend
go mod download
```

### 运行服务器

```bash
# 开发模式
go run cmd/server/main.go

# 编译后运行
go build -o server cmd/server/main.go
./server
```

服务器将在 `http://localhost:8080` 启动。

### 测试API

```bash
# 健康检查
curl http://localhost:8080/health

# 列出并发模式
curl http://localhost:8080/api/v1/patterns

# 分析CFG
curl -X POST http://localhost:8080/api/v1/analysis/cfg \
  -H "Content-Type: application/json" \
  -d '{"code":"package main\n\nfunc main() {\n\tx := 0\n\tif x < 10 {\n\t\tx++\n\t}\n}"}'

# 生成并发模式
curl -X POST http://localhost:8080/api/v1/patterns/generate \
  -H "Content-Type: application/json" \
  -d '{"pattern":"worker-pool","parameters":{"workers":"10","bufferSize":"100"}}'
```

### WebSocket测试

使用 `wscat` (需要先安装: `npm install -g wscat`):

```bash
wscat -c ws://localhost:8080/ws
```

## 📁 项目结构

```text
backend/
├── cmd/
│   └── server/
│       └── main.go           # 服务器入口
├── internal/
│   ├── api/                  # REST API处理器
│   │   ├── analysis.go       # 分析相关API
│   │   ├── patterns.go       # 并发模式API
│   │   └── projects.go       # 项目管理API
│   ├── ws/                   # WebSocket处理器
│   │   └── handler.go
│   └── service/              # 业务逻辑服务
├── go.mod
└── README.md
```

## 🔧 配置

### 环境变量

- `PORT` - 服务器端口 (默认: 8080)
- `GIN_MODE` - Gin模式 (debug/release)
- `CORS_ORIGINS` - 允许的CORS来源

### 示例

```bash
export PORT=3000
export GIN_MODE=release
export CORS_ORIGINS="http://localhost:5173,https://yourdomain.com"
```

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📊 API响应格式

### 成功响应

```json
{
  "success": true,
  "data": { ... },
  "time": "2025-10-23T12:00:00Z"
}
```

### 错误响应

```json
{
  "success": false,
  "error": "Error message",
  "time": "2025-10-23T12:00:00Z"
}
```

## 🔌 与Formal Verifier集成

后端API将调用 `formal-verifier` 工具进行实际分析：

```go
// 示例集成代码
import "github.com/your-org/go-formal-verification/tools/formal-verifier/pkg/analyzer"

result, err := analyzer.AnalyzeCFG(code)
if err != nil {
    // 处理错误
}
// 返回结果
```

## 🚧 开发状态

当前状态：**Alpha**

- [x] 基础架构搭建
- [x] REST API骨架
- [x] WebSocket支持
- [x] CORS配置
- [ ] 与formal-verifier集成
- [ ] 数据持久化
- [ ] 认证授权
- [ ] 性能优化

## 📝 TODO

### 高优先级

- [ ] 集成 formal-verifier 实际分析功能
- [ ] 实现数据库持久化 (SQLite)
- [ ] 添加用户认证 (JWT)
- [ ] 实现分析任务队列

### 中优先级

- [ ] 添加缓存层 (Redis)
- [ ] 实现速率限制
- [ ] 添加请求日志
- [ ] 性能监控

### 低优先级

- [ ] GraphQL API支持
- [ ] gRPC支持
- [ ] 多租户支持
- [ ] 集群部署

## 🤝 贡献

欢迎贡献！请查看主项目的 [CONTRIBUTING.md](../../CONTRIBUTING.md)

## 📄 许可证

MIT License - 详见 [LICENSE](../../LICENSE)

## 📞 联系方式

- Issues: <https://github.com/your-org/go-formal-verification/issues>
- Email: <support@go-formal-verification.org>

---

**Go Formal Verification Framework - Web UI Backend**  
*理论驱动，工程落地，持续创新！* 🚀
