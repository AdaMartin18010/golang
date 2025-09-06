# Go语言现代化 - 完整测试体系 Dockerfile

# 多阶段构建
# 第一阶段：构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.Version=${VERSION:-1.0.0} -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.GitCommit=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" \
    -a -installsuffix cgo \
    -o testing-system .

# 第二阶段：运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/testing-system .

# 创建必要的目录
RUN mkdir -p /app/test-results /app/performance-results /app/profiles && \
    chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 设置环境变量
ENV GIN_MODE=release
ENV PORT=8080
ENV TEST_OUTPUT_DIR=/app/test-results
ENV PERF_OUTPUT_DIR=/app/performance-results
ENV PROFILE_OUTPUT_DIR=/app/profiles

# 设置入口点
ENTRYPOINT ["./testing-system"]

# 默认命令
CMD ["--demo"]
