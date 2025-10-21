# Go与Kubernetes入门

> **简介**: Kubernetes核心概念与Go应用部署实战，掌握Pod、Deployment、Service等基础资源
> **版本**: Go 1.23+ / Kubernetes 1.28+  
> **难度**: ⭐⭐⭐⭐  
> **标签**: #Kubernetes #K8s #容器编排 #云原生

<!-- TOC START -->
- [Go与Kubernetes入门](#go与kubernetes入门)
  - [📚 **理论分析**](#-理论分析)
  - [🛠️ **核心概念**](#️-核心概念)
  - [💻 **部署流程与YAML示例**](#-部署流程与yaml示例)
    - [**Deployment示例**](#deployment示例)
    - [**Service示例**](#service示例)
  - [🎯 **最佳实践**](#-最佳实践)
  - [🔍 **常见问题**](#-常见问题)
  - [📚 **扩展阅读**](#-扩展阅读)
<!-- TOC END -->

## 📚 **理论分析**

- Kubernetes（K8s）是主流容器编排平台，实现服务自动部署、扩缩容、健康检查等。
- Go服务与K8s天然兼容，易于云原生部署。

## 🛠️ **核心概念**

- Pod：最小部署单元，封装一个或多个容器
- Service：服务发现与负载均衡
- Deployment：声明式部署与滚动升级
- ConfigMap/Secret：配置与密钥管理

## 💻 **部署流程与YAML示例**

### **Deployment示例**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-demo
spec:
  replicas: 2
  selector:
    matchLabels:
      app: go-demo
  template:
    metadata:
      labels:
        app: go-demo
    spec:
      containers:
      - name: go-demo
        image: go-demo:latest
        ports:
        - containerPort: 8080

```

### **Service示例**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: go-demo-svc
spec:
  selector:
    app: go-demo
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP

```

## 🎯 **最佳实践**

- 镜像小型化，启动健康检查
- 配置与密钥分离，使用ConfigMap/Secret
- 资源限制（CPU/内存）合理配置
- 滚动升级与回滚策略

## 🔍 **常见问题**

- Q: Go服务如何暴露外部访问？
  A: 配置Service为NodePort或Ingress
- Q: 如何调试K8s中的Go服务？
  A: 查看Pod日志，kubectl exec进入容器

## 📚 **扩展阅读**

- [Go与Kubernetes实战](https://geektutu.com/post/hpg-golang-k8s.html)
- [Kubernetes官方文档](https://kubernetes.io/zh/docs/)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: Go 1.25.3+
