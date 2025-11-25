# 多阶段构建
# Stage 1: 构建
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git make

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Stage 2: 运行
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .

# 复制配置文件
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./server"]
