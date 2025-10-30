# Go最佳实践

Go开发最佳实践，涵盖测试、部署、工程化和可观测性。

---

## 📋 目录结构

### 核心模块

1. **[测试](./testing/README.md)** ⭐⭐⭐⭐⭐
   - 单元测试
   - 表格驱动测试
   - 集成测试
   - 性能测试
   - Mock与Stub

2. **[部署](./deployment/README.md)** ⭐⭐⭐⭐⭐
   - Docker部署
   - Kubernetes部署
   - CI/CD流程
   - 监控与日志
   - 滚动更新

3. **[工程化](./engineering/README.md)** ⭐⭐⭐⭐⭐
   - 代码规范
   - 项目结构
   - 完整测试体系
   - 文档编写
   - 版本管理

4. **[可观测性](./observability/README.md)** ⭐⭐⭐⭐⭐
   - 日志管理
   - 指标监控
   - 链路追踪
   - 告警管理

---

## 🎯 实践路径

### 测试工程师 (2-3周)

```text
单元测试 → Mock → 集成测试 → 性能测试 → 覆盖率
```

### DevOps工程师 (3-4周)

```text
Docker → Kubernetes → CI/CD → 监控 → 告警
```

### 工程化专家 (2-3周)

```text
代码规范 → 项目结构 → 文档 → 版本管理
```

---

## 🚀 快速开始

### 单元测试

```go
package calculator

import "testing"

func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -1, -2},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Docker部署

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

### Kubernetes部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: myapp:latest
        ports:
        - containerPort: 8080
```

---

## 📖 系统文档

- **[知识图谱](./00-知识图谱.md)**: 实践知识体系全景图
- **[对比矩阵](./00-对比矩阵.md)**: 实践方案对比
- **[概念定义体系](./00-概念定义体系.md)**: 核心概念详解

---

## 🛠️ 常用工具

### 测试工具

- testing - 标准测试库
- testify - 断言库
- gomock - Mock框架
- gocov - 覆盖率工具

### CI/CD工具

- GitHub Actions
- GitLab CI
- Jenkins
- Drone

### 监控工具

- Prometheus - 监控
- Grafana - 可视化
- ELK Stack - 日志
- Jaeger - 追踪

---

## 📚 推荐阅读顺序

1. **测试** → 单元测试 → Mock → 集成测试 → 性能测试
2. **部署** → Docker → Kubernetes → CI/CD
3. **工程化** → 规范 → 结构 → 文档
4. **可观测性** → 日志 → 监控 → 追踪

---

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3
