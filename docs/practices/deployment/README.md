# Go部署实践

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Go部署实践](#go部署实践)
  - [� 目录](#-目录)
  - [📚 核心内容](#-核心内容)
  - [🚀 快速开始](#-快速开始)
    - [Dockerfile](#dockerfile)
    - [Kubernetes](#kubernetes)
  - [📖 系统文档](#-系统文档)

---

## 📚 核心内容

1. **[部署概览](./01-部署概览.md)** ⭐⭐⭐⭐⭐
2. **[Docker部署](./02-Docker部署.md)** ⭐⭐⭐⭐⭐
3. **[Kubernetes部署](./03-Kubernetes部署.md)** ⭐⭐⭐⭐⭐
4. **[CI/CD流程](./04-CI-CD流程.md)** ⭐⭐⭐⭐
5. **[监控与日志](./05-监控与日志.md)** ⭐⭐⭐⭐⭐
6. **[滚动更新](./06-滚动更新.md)** ⭐⭐⭐⭐
7. **[生产环境最佳实践](./07-生产环境最佳实践.md)** ⭐⭐⭐⭐⭐

---

## 🚀 快速开始

### Dockerfile

```dockerfile
FROM golang:1.25.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
COPY --from=builder /app/main .
CMD ["./main"]
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:latest
```

---

## 📖 系统文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3
