# 开发环境搭建指南

> **版本**: v1.0  
> **更新日期**: 2026-04-02  
> **预计时间**: 30 分钟

---

## 📋 前置要求

| 工具 | 版本 | 说明 |
|------|------|------|
| Go | 1.26+ | 编程语言 |
| Docker | 20.10+ | 容器化 |
| Docker Compose | 2.0+ | 本地开发 |
| Make | - | 构建工具 |
| Git | 2.30+ | 版本控制 |

---

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/yourusername/golang.git
cd golang
```

### 2. 安装依赖

```bash
# 下载 Go 依赖
go mod download

# 安装工具依赖
make install-tools
```

### 3. 启动基础设施

```bash
# 启动 PostgreSQL, Redis, NATS
docker-compose -f docker-compose.dev.yml up -d
```

### 4. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env，设置本地开发配置
```

### 5. 数据库迁移

```bash
make migrate-up
```

### 6. 生成代码

```bash
# 生成 Ent 代码
make generate-ent

# 生成 gRPC 代码
make generate-grpc

# 生成 Wire 代码
make generate-wire
```

### 7. 运行应用

```bash
# 方式 1: 直接运行
go run cmd/server/main.go

# 方式 2: 使用 Air 热重载
air

# 方式 3: Docker
docker-compose up app
```

### 8. 验证运行

```bash
# 健康检查
curl http://localhost:8080/health

# 预期响应
{"status":"healthy"}
```

---

## 🧪 运行测试

```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 查看覆盖率
make coverage
```

---

## 📝 开发工作流

### 创建新功能

```bash
# 1. 创建分支
git checkout -b feature/my-feature

# 2. 开发代码
# ...

# 3. 运行测试
make test

# 4. 提交代码
git add .
git commit -m "feat: add my feature"

# 5. 推送分支
git push origin feature/my-feature

# 6. 创建 Pull Request
```

### 代码规范

```bash
# 格式化代码
make fmt

# 运行 Linter
make lint

# 自动修复
make lint-fix
```

---

## 🔧 常用命令

| 命令 | 说明 |
|------|------|
| `make dev` | 启动开发服务器 |
| `make build` | 构建二进制文件 |
| `make test` | 运行测试 |
| `make clean` | 清理构建产物 |
| `make generate` | 生成所有代码 |

---

## 🆘 常见问题

### Q: 数据库连接失败
**A**: 检查 `DB_HOST` 是否为 `localhost`，Docker 网络是否正确

### Q: 端口被占用
**A**: 修改 `.env` 中的端口配置，或停止占用端口的进程

### Q: 代码生成失败
**A**: 确保安装了所有工具: `make install-tools`

---

## 📚 下一步

- [API 文档](../api/README.md)
- [架构设计](../architecture/clean-architecture.md)
- [部署指南](../deployment/README.md)

---

*最后更新: 2026-04-02*
