# Docker部署Go应用

> **简介**: Docker部署Go应用完整指南，包括Dockerfile编写、多阶段构建和最佳实践

> **版本**: Go 1.25.3, Docker 24+  
> **难度**: ⭐⭐⭐  
> **标签**: #Docker #部署 #容器化

---

## 📋 目录

- [1. 基础Dockerfile](#1-基础dockerfile)
  - [简单Dockerfile](#简单dockerfile)
- [2. 多阶段构建](#2-多阶段构建)
  - [标准多阶段构建](#标准多阶段构建)
  - [使用Scratch镜像](#使用scratch镜像)
- [3. 优化镜像大小](#3-优化镜像大小)
  - [利用构建缓存](#利用构建缓存)
  - [使用.dockerignore](#使用dockerignore)
  - [优化构建参数](#优化构建参数)
- [4. Docker Compose](#4-docker-compose)
  - [单服务](#单服务)
  - [完整微服务栈](#完整微服务栈)
- [5. 最佳实践](#5-最佳实践)
  - [1. 使用非root用户](#1-使用非root用户)
  - [2. 健康检查](#2-健康检查)
  - [3. 使用环境变量](#3-使用环境变量)
  - [4. 使用构建参数](#4-使用构建参数)
  - [5. 多架构构建](#5-多架构构建)
  - [6. 完整生产Dockerfile](#6-完整生产dockerfile)
- [🔗 相关资源](#-相关资源)

## 1. 基础Dockerfile

### 简单Dockerfile

```dockerfile
FROM golang:1.25.3

WORKDIR /app

# 复制源代码
COPY . .

# 下载依赖
RUN go mod download

# 构建应用
RUN go build -o main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
```

**构建和运行**:
```bash
# 构建镜像
docker build -t myapp:latest .

# 运行容器
docker run -p 8080:8080 myapp:latest
```

---

## 2. 多阶段构建

### 标准多阶段构建

```dockerfile
# 构建阶段
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest

# 安装CA证书（用于HTTPS请求）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从builder阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
```

**优势**:
- 最终镜像不包含Go编译器
- 镜像大小显著减小（从800MB+ → 10MB+）
- 更安全（攻击面小）

---

### 使用Scratch镜像

```dockerfile
# 构建阶段
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 静态链接构建
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o main .

# 运行阶段：使用scratch（最小镜像）
FROM scratch

# 复制CA证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 复制二进制文件
COPY --from=builder /app/main /main

# 暴露端口
EXPOSE 8080

# 运行应用
ENTRYPOINT ["/main"]
```

**优势**:
- 镜像大小最小（<10MB）
- 最安全
- 适合纯Go应用（无CGO依赖）

---

## 3. 优化镜像大小

### 利用构建缓存

```dockerfile
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# 先复制依赖文件（利用Docker缓存）
COPY go.mod go.sum ./
RUN go mod download

# 再复制源代码（代码变更不会重新下载依赖）
COPY . .

RUN CGO_ENABLED=0 go build -o main .

FROM alpine:latest
COPY --from=builder /app/main .
CMD ["./main"]
```

---

### 使用.dockerignore

```
# .dockerignore
.git
.gitignore
README.md
*.md
.env
.env.*
*.test
coverage.txt
.vscode
.idea
*.swp
*.log
tmp/
vendor/
```

---

### 优化构建参数

```dockerfile
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 优化构建标志
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \           # 去除调试信息和符号表
    -trimpath \                   # 去除文件路径信息
    -o main .

FROM scratch
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]
```

**标志说明**:
- `-ldflags="-w -s"`: 去除调试信息，减小25-30%大小
- `-trimpath`: 去除文件系统路径，提高可重复性
- `CGO_ENABLED=0`: 静态链接，无libc依赖

---

## 4. Docker Compose

### 单服务

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - ENV=production
      - DB_HOST=db
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=myapp
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  postgres_data:
```

**运行**:
```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止所有服务
docker-compose down
```

---

### 完整微服务栈

```yaml
version: '3.8'

services:
  # Go应用
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - REDIS_URL=redis:6379
      - DB_HOST=postgres
    depends_on:
      - postgres
      - redis
    networks:
      - app-network

  # PostgreSQL
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - app-network

  # Redis
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - app-network

  # Nginx反向代理
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    networks:
      - app-network

volumes:
  postgres_data:
  redis_data:

networks:
  app-network:
    driver: bridge
```

---

## 5. 最佳实践

### 1. 使用非root用户

```dockerfile
FROM alpine:latest

# 创建非root用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/main .

# 切换到非root用户
USER appuser

CMD ["./main"]
```

---

### 2. 健康检查

```dockerfile
FROM alpine:latest

COPY --from=builder /app/main .

# 添加健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]
```

**Go应用中实现健康检查端点**:
```go
func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func main() {
    http.HandleFunc("/health", healthCheck)
    http.ListenAndServe(":8080", nil)
}
```

---

### 3. 使用环境变量

```dockerfile
FROM alpine:latest

COPY --from=builder /app/main .

# 设置默认环境变量
ENV PORT=8080 \
    ENV=production \
    LOG_LEVEL=info

EXPOSE ${PORT}

CMD ["./main"]
```

**Go代码中读取环境变量**:
```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}

env := os.Getenv("ENV")
logLevel := os.Getenv("LOG_LEVEL")
```

---

### 4. 使用构建参数

```dockerfile
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 接收构建参数
ARG VERSION=dev
ARG BUILD_TIME

# 将参数传递给构建
RUN CGO_ENABLED=0 go build \
    -ldflags="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o main .

FROM alpine:latest
COPY --from=builder /app/main .
CMD ["./main"]
```

**构建时传递参数**:
```bash
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t myapp:1.0.0 .
```

**Go代码中使用**:
```go
var (
    Version   string
    BuildTime string
)

func main() {
    fmt.Printf("Version: %s\n", Version)
    fmt.Printf("Build Time: %s\n", BuildTime)
}
```

---

### 5. 多架构构建

```dockerfile
FROM --platform=$BUILDPLATFORM golang:1.25.3-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o main .

FROM alpine:latest
COPY --from=builder /app/main .
CMD ["./main"]
```

**构建多架构镜像**:
```bash
# 创建builder
docker buildx create --name multiarch --use

# 构建并推送
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t myapp:latest \
  --push .
```

---

### 6. 完整生产Dockerfile

```dockerfile
# 构建阶段
FROM golang:1.25.3-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git

WORKDIR /app

# 利用缓存下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建参数
ARG VERSION=dev
ARG BUILD_TIME

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -trimpath \
    -o main .

# 运行阶段
FROM alpine:latest

# 安装CA证书和tzdata
RUN apk --no-cache add ca-certificates tzdata

# 创建非root用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/main .

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 运行应用
CMD ["./main"]
```

---

## 🔗 相关资源

- [部署概览](./01-部署概览.md)
- [Kubernetes部署](./03-Kubernetes部署.md)
- [CI/CD流程](./04-CI-CD流程.md)

---

**最后更新**: 2025-10-28  
**Go版本**: 1.25.3  
**Docker版本**: 24+
