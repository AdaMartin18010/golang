# Go 1.25 容器感知调度示例 Dockerfile

# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 复制代码
COPY . .

# 构建应用
# Go 1.25+ 会自动启用容器感知调度，无需额外配置
RUN go build -o container-aware-app main.go
RUN go test -c -o container-test container_test.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/container-aware-app .
COPY --from=builder /app/container-test .

# Go 1.25+ 应用会自动检测容器的 CPU 限制
# 无需设置 GOMAXPROCS 环境变量
# ENV GOMAXPROCS=2  # 不推荐，会覆盖自动检测

# 暴露端口（如果需要）
# EXPOSE 8080

# 运行应用
CMD ["./container-aware-app"]

# 使用方法:
# docker build -t go-container-aware .
# docker run --cpus=2 go-container-aware
# docker run --cpus=4 go-container-aware
# docker run --cpus=8 go-container-aware

