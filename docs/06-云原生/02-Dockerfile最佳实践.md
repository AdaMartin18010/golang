# Dockerfile最佳实践

<!-- TOC START -->
- [Dockerfile最佳实践](#dockerfile最佳实践)
  - [📚 **理论分析**](#-理论分析)
  - [🛠️ **常见优化技巧**](#️-常见优化技巧)
  - [💻 **代码示例**](#-代码示例)
    - [**推荐Dockerfile模板**](#推荐dockerfile模板)
    - [**健康检查与非root用户**](#健康检查与非root用户)
  - [🎯 **最佳实践**](#-最佳实践)
  - [🔍 **常见问题**](#-常见问题)
  - [📚 **扩展阅读**](#-扩展阅读)
<!-- TOC END -->

## 📚 **理论分析**

- Dockerfile定义镜像构建流程，影响镜像体积、安全、可维护性。
- Go项目适合多阶段构建，产物小、依赖少。

## 🛠️ **常见优化技巧**

- 多阶段构建，分离编译与运行环境
- COPY/ADD顺序优化，减少缓存失效
- 只COPY必要文件，避免泄漏敏感信息
- 明确EXPOSE端口，CMD/ENTRYPOINT分离
- 使用alpine等精简基础镜像

## 💻 **代码示例**

### **推荐Dockerfile模板**

```dockerfile
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 8080
ENTRYPOINT ["./app"]

```

### **健康检查与非root用户**

```dockerfile
HEALTHCHECK CMD curl --fail http://localhost:8080/health || exit 1
RUN adduser -D appuser
USER appuser

```

## 🎯 **最佳实践**

- 镜像最小化，减少攻击面
- 不在生产镜像中保留源码/编译工具
- 健康检查保障服务可用性
- 使用非root用户运行
- 合理利用缓存加速构建

## 🔍 **常见问题**

- Q: 为什么要用多阶段构建？
  A: 分离编译与运行，减小镜像体积
- Q: 如何避免缓存失效？
  A: 先COPY go.mod/go.sum再COPY源码

## 📚 **扩展阅读**

- [Go Dockerfile最佳实践](https://geektutu.com/post/hpg-golang-dockerfile.html)
- [Dockerfile官方文档](https://docs.docker.com/engine/reference/builder/)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: Go 1.25.3+
