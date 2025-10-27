# Docker部署

**难度**: 中级 | **预计阅读**: 20分钟 | **前置知识**: Docker基础

---

## 📖 概念介绍

Docker容器化是现代Go应用部署的标准方式，提供了环境一致性、快速部署和易于扩展的优势。

---

## 🎯 核心知识点

### 1. 基础Dockerfile

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
```

构建和运行：
```bash
docker build -t myapp:latest .
docker run -p 8080:8080 myapp:latest
```

---

### 2. 多阶段构建优化

```dockerfile
# 优化的多阶段构建
FROM golang:1.21-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git

WORKDIR /app

# 缓存依赖层
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译（静态链接，小体积）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags='-w -s' \
    -o main .

# 最小运行镜像
FROM scratch

# 只复制必需文件
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /main

EXPOSE 8080

ENTRYPOINT ["/main"]
```

**优势**:
- 镜像体积：从800MB → 10MB
- 安全性：最小攻击面
- 启动快：无额外进程

---

### 3. Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/myapp
      - REDIS_URL=redis://redis:6379
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_started
    restart: unless-stopped
    
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: myapp
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user"]
      interval: 10s
      timeout: 5s
      retries: 5
    
  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

运行：
```bash
docker-compose up -d
docker-compose logs -f app
docker-compose down
```

---

### 4. .dockerignore

```
# .dockerignore
.git
.gitignore
README.md
Dockerfile
docker-compose.yml
.env
*.md

# 测试文件
*_test.go
testdata/

# 构建产物
bin/
*.exe

# IDE
.vscode/
.idea/

# 依赖
vendor/
```

---

### 5. 健康检查

```dockerfile
FROM alpine:latest

COPY main /main

# 添加健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["/main"]
```

```go
// 健康检查端点
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
```

---

## 💡 最佳实践

### 1. 镜像体积优化

```dockerfile
# ✅ 好：使用alpine
FROM golang:1.21-alpine

# ✅ 好：使用scratch（最小）
FROM scratch

# ❌ 差：使用完整镜像
FROM golang:1.21  # 太大
```

### 2. 层缓存优化

```dockerfile
# ✅ 好：先复制依赖文件
COPY go.mod go.sum ./
RUN go mod download  # 缓存层

COPY . .  # 代码变化不影响依赖层
RUN go build

# ❌ 差：一起复制
COPY . .  # 每次都重新下载依赖
RUN go mod download
RUN go build
```

### 3. 安全实践

```dockerfile
# ✅ 创建非root用户
FROM alpine:latest

RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

USER appuser

COPY --chown=appuser:appgroup main /main

CMD ["/main"]
```

### 4. 构建参数

```dockerfile
ARG VERSION=dev
ARG BUILD_DATE

LABEL version="${VERSION}" \
      build_date="${BUILD_DATE}"

RUN go build -ldflags="-X main.Version=${VERSION}"
```

```bash
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t myapp:1.0.0 .
```

---

## 📚 相关资源

- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)

**下一步**: [03-Kubernetes部署](./03-Kubernetes部署.md)

---

**最后更新**: 2025-10-28

