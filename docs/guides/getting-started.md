# 快速开始指南

## 前置要求

- Go 1.25.3+
- PostgreSQL 15+
- Docker (可选)

## 安装步骤

### 1. 克隆项目

```bash
git clone <repository-url>
cd golang
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境

复制配置文件：

```bash
cp configs/.env.example configs/.env
```

编辑 `configs/.env` 设置数据库连接等信息。

### 4. 启动数据库

使用 Docker Compose：

```bash
cd deployments/docker
docker-compose up -d postgres
```

### 5. 运行应用

```bash
go run cmd/server/main.go
```

应用将在 `http://localhost:8080` 启动。

## 验证

访问健康检查端点：

```bash
curl http://localhost:8080/health
```

应该返回 `OK`。

## 下一步

- 查看 [API 文档](../api/)
- 阅读 [架构文档](../architecture/)
- 了解 [开发指南](development.md)
