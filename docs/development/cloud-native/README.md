# Go云原生开发

Go云原生开发完整指南，涵盖Docker、Kubernetes、服务网格和云原生实践。

---

## 📚 核心内容

1. **Docker容器化**
2. **Kubernetes部署**
3. **服务网格 (Service Mesh)**
4. **配置管理**
5. **CI/CD流程**
6. **云原生最佳实践**

---

## 🚀 Docker示例

```dockerfile
FROM golang:1.25.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

---

## 📖 系统文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**最后更新**: 2025-10-28
