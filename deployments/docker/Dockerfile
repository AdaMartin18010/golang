# 多阶段构建 Dockerfile
# 用于构建和运行 Go 应用

# 阶段 1: 依赖下载
FROM golang:1.21-alpine AS deps
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 阶段 2: 代码构建
FROM golang:1.21-alpine AS builder
WORKDIR /app

# 从依赖阶段复制模块缓存
COPY --from=deps /go/pkg/mod /go/pkg/mod

# 复制源代码
COPY . .

# 构建应用
# -ldflags='-w -s': 去除调试信息，减小二进制大小
# -a -installsuffix cgo: 静态链接，不依赖 CGO
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app ./cmd/server

# 阶段 3: 运行环境
FROM alpine:latest

# 安装必要的工具
RUN apk --no-cache add ca-certificates tzdata wget

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/app .

# 更改文件所有者
RUN chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./app"]
