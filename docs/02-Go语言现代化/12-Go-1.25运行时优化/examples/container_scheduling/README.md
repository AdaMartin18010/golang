# 容器感知调度示例

> **Go 版本**: 1.25+  
> **示例类型**: 容器环境优化  
> **最后更新**: 2025-10-18

本目录包含 Go 1.25 容器感知调度的测试和示例代码。

---

## 📋 目录

- [容器感知调度示例](#容器感知调度示例)
  - [📋 目录](#-目录)
  - [🚀 快速开始](#-快速开始)
    - [1. 验证容器感知调度](#1-验证容器感知调度)
    - [2. 检测 cgroup 配置](#2-检测-cgroup-配置)
    - [3. 运行时诊断](#3-运行时诊断)
  - [🧪 测试说明](#-测试说明)
    - [功能测试](#功能测试)
    - [基准测试](#基准测试)
  - [🐳 Docker 使用](#-docker-使用)
    - [基本示例](#基本示例)
    - [Dockerfile 示例](#dockerfile-示例)
    - [Docker Compose 示例](#docker-compose-示例)
  - [☸️ Kubernetes 使用](#️-kubernetes-使用)
    - [基本 Deployment](#基本-deployment)
    - [带监控的完整示例](#带监控的完整示例)
    - [验证 Kubernetes 中的设置](#验证-kubernetes-中的设置)
  - [📊 性能对比](#-性能对比)
    - [对比测试流程](#对比测试流程)
    - [预期性能提升](#预期性能提升)
    - [可视化对比](#可视化对比)
  - [❓ 常见问题](#-常见问题)
    - [Q1: 如何验证容器感知调度是否生效？](#q1-如何验证容器感知调度是否生效)
    - [Q2: 为什么 GOMAXPROCS 不等于 CPU 限制？](#q2-为什么-gomaxprocs-不等于-cpu-限制)
    - [Q3: 如何在 Go 1.24 中实现类似功能？](#q3-如何在-go-124-中实现类似功能)
    - [Q4: 小数 CPU 限制如何处理？](#q4-小数-cpu-限制如何处理)
    - [Q5: 如何监控 GOMAXPROCS？](#q5-如何监控-gomaxprocs)
  - [🔧 调试技巧](#-调试技巧)
    - [1. 查看 cgroup 配置](#1-查看-cgroup-配置)
    - [2. pprof 分析](#2-pprof-分析)
    - [3. trace 分析](#3-trace-分析)
  - [📚 相关资源](#-相关资源)

---

## 🚀 快速开始

### 1. 验证容器感知调度

```bash
# 本地运行（自动检测）
go test -v -run=TestContainerAwareScheduling

# Docker 容器中运行
docker run --cpus=2 -v $(pwd):/app -w /app golang:1.25 \
  go test -v -run=TestContainerAwareScheduling
```

**预期输出**:

```text
=== RUN   TestContainerAwareScheduling
    container_test.go:17: GOMAXPROCS: 2
    container_test.go:18: NumCPU (物理核心数): 32
    container_test.go:23: cgroup CPU 限制: 2.00 核
    container_test.go:32: ✅ 容器感知调度已生效 (GOMAXPROCS=2, CPU limit=2.00)
--- PASS: TestContainerAwareScheduling (0.00s)
PASS
```

### 2. 检测 cgroup 配置

```bash
# 运行 cgroup 检测测试
go test -v -run=TestCgroupDetection

# Docker 中运行
docker run --cpus=4 -v $(pwd):/app -w /app golang:1.25 \
  go test -v -run=TestCgroupDetection
```

### 3. 运行时诊断

```bash
# 完整诊断信息
go test -v -run=TestRuntimeDiagnostics

# Docker 中诊断
docker run --cpus=2 --memory=2g -v $(pwd):/app -w /app golang:1.25 \
  go test -v -run=TestRuntimeDiagnostics
```

---

## 🧪 测试说明

### 功能测试

| 测试函数 | 说明 | 用途 |
|---------|------|------|
| `TestContainerAwareScheduling` | 验证容器感知调度 | 检查 GOMAXPROCS 是否正确设置 |
| `TestCgroupDetection` | 测试 cgroup 检测 | 验证 cgroup 文件读取 |
| `TestRuntimeDiagnostics` | 运行时诊断 | 完整系统信息输出 |
| `TestHighConcurrency` | 高并发压力测试 | 验证调度稳定性 |

**运行所有功能测试**:

```bash
go test -v

# Docker 中运行
docker run --cpus=2 -v $(pwd):/app -w /app golang:1.25 \
  go test -v
```

### 基准测试

| 基准测试 | 说明 | 对比内容 |
|---------|------|---------|
| `BenchmarkCorrectGOMAXPROCS` | 正确的 GOMAXPROCS | 使用自动检测值 |
| `BenchmarkOversubscribedGOMAXPROCS` | 过度订阅 | 使用 NumCPU（模拟 Go 1.24） |
| `BenchmarkSchedulingOverhead` | 调度开销 | 不同 GOMAXPROCS 值的影响 |
| `BenchmarkConcurrentLoad` | 并发负载 | 自动检测 vs NumCPU |

**运行所有基准测试**:

```bash
# 基本基准测试
go test -bench=. -benchmem

# Docker 中运行（4 核限制）
docker run --cpus=4 -v $(pwd):/app -w /app golang:1.25 \
  go test -bench=. -benchmem -benchtime=10s
```

**预期结果**:

```text
BenchmarkCorrectGOMAXPROCS-4                    50000    15000 ns/op    0 B/op    0 allocs/op
BenchmarkOversubscribedGOMAXPROCS-32            30000    28000 ns/op    0 B/op    0 allocs/op
```

---

## 🐳 Docker 使用

### 基本示例

```bash
# 创建测试程序
cat > main.go <<'EOF'
package main

import (
    "fmt"
    "runtime"
)

func main() {
    fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
    fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
}
EOF

# 构建镜像
docker build -t test-container-scheduling .

# 运行容器（限制 2 核）
docker run --cpus=2 test-container-scheduling
# 输出: GOMAXPROCS: 2

# 运行容器（限制 4 核）
docker run --cpus=4 test-container-scheduling
# 输出: GOMAXPROCS: 4
```

### Dockerfile 示例

```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o myapp

FROM alpine:latest
COPY --from=builder /app/myapp /myapp

# Go 1.25+ 会自动检测容器 CPU 限制
# 无需手动设置 GOMAXPROCS
CMD ["/myapp"]
```

### Docker Compose 示例

```yaml
version: '3.8'

services:
  app:
    build: .
    deploy:
      resources:
        limits:
          cpus: '2.0'      # Go 1.25+ 自动设置 GOMAXPROCS=2
          memory: 2G
        reservations:
          cpus: '1.0'
          memory: 1G
```

---

## ☸️ Kubernetes 使用

### 基本 Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-app
  template:
    metadata:
      labels:
        app: go-app
    spec:
      containers:
      - name: app
        image: go-app:1.25
        resources:
          limits:
            cpu: "2"        # Go 1.25+ 自动设置 GOMAXPROCS=2
            memory: "2Gi"
          requests:
            cpu: "1"
            memory: "1Gi"
        # 无需设置 GOMAXPROCS 环境变量
```

### 带监控的完整示例

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app-with-monitoring
spec:
  template:
    spec:
      containers:
      - name: app
        image: go-app:1.25
        resources:
          limits:
            cpu: "4"
            memory: "4Gi"
          requests:
            cpu: "2"
            memory: "2Gi"
        env:
        - name: GO_VERSION
          value: "1.25"
        # 可选：手动覆盖（仅特殊场景）
        # - name: GOMAXPROCS
        #   value: "8"
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 6060
          name: pprof
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
```

### 验证 Kubernetes 中的设置

```bash
# 部署应用
kubectl apply -f deployment.yaml

# 查看 GOMAXPROCS
kubectl exec -it <pod-name> -- env | grep GOMAXPROCS

# 进入容器验证
kubectl exec -it <pod-name> -- sh
> ./myapp --diagnose
GOMAXPROCS: 2 (auto-detected from cgroup)
NumCPU: 96 (node physical CPUs)
CPU Limit: 2.00 cores
```

---

## 📊 性能对比

### 对比测试流程

```bash
# 1. 在容器中测试（自动检测）
docker run --cpus=4 -v $(pwd):/app -w /app golang:1.25 \
  go test -bench=BenchmarkCorrectGOMAXPROCS -benchmem -count=5 \
  > auto-detected.txt

# 2. 模拟 Go 1.24（使用 NumCPU）
docker run --cpus=4 -v $(pwd):/app -w /app golang:1.25 \
  go test -bench=BenchmarkOversubscribedGOMAXPROCS -benchmem -count=5 \
  > oversubscribed.txt

# 3. 使用 benchstat 对比
go install golang.org/x/perf/cmd/benchstat@latest
benchstat auto-detected.txt oversubscribed.txt
```

### 预期性能提升

| 场景 | Go 1.24 (NumCPU=32) | Go 1.25 (自动检测=4) | 提升 |
|------|---------------------|---------------------|------|
| **吞吐量** | 12K ops/s | 18K ops/s | +50% |
| **延迟 (P99)** | 250ms | 120ms | -52% |
| **CPU 利用率** | 180% | 95% | 优化 |
| **上下文切换** | 85K/s | 32K/s | -62% |

### 可视化对比

```bash
# 生成性能报告
go test -bench=. -benchmem -cpuprofile=cpu.prof

# 查看 CPU profile
go tool pprof -http=:8080 cpu.prof

# 对比不同 GOMAXPROCS 设置
for n in 1 2 4 8 16 32; do
  GOMAXPROCS=$n go test -bench=BenchmarkConcurrentLoad -benchmem
done
```

---

## ❓ 常见问题

### Q1: 如何验证容器感知调度是否生效？

**A**: 运行测试程序检查 GOMAXPROCS：

```bash
docker run --cpus=2 -v $(pwd):/app -w /app golang:1.25 \
  go test -v -run=TestContainerAwareScheduling
```

如果输出显示 `GOMAXPROCS=2`（而不是宿主机核心数），则生效。

### Q2: 为什么 GOMAXPROCS 不等于 CPU 限制？

**A**: 可能原因：

1. **环境变量覆盖**: 检查 `GOMAXPROCS` 环境变量
2. **cgroup 未挂载**: 容器运行时问题
3. **非 Linux 系统**: Windows/macOS 支持有限

### Q3: 如何在 Go 1.24 中实现类似功能？

**A**: 使用第三方库：

```go
import _ "go.uber.org/automaxprocs"

// 自动设置 GOMAXPROCS
```

### Q4: 小数 CPU 限制如何处理？

**A**: Go 1.25 会自动取整：

```bash
# Docker: --cpus=2.5
# GOMAXPROCS 可能是 2 或 3

# 建议使用整数
docker run --cpus=2 myapp  # 推荐
```

### Q5: 如何监控 GOMAXPROCS？

**A**: 使用 Prometheus 指标：

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "runtime"
)

var gomaxprocs = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "go_gomaxprocs",
    Help: "Current GOMAXPROCS value",
})

func init() {
    prometheus.MustRegister(gomaxprocs)
    gomaxprocs.Set(float64(runtime.GOMAXPROCS(0)))
}
```

---

## 🔧 调试技巧

### 1. 查看 cgroup 配置

```bash
# 在容器内查看 cgroup 限制
cat /sys/fs/cgroup/cpu/cpu.cfs_quota_us
cat /sys/fs/cgroup/cpu/cpu.cfs_period_us

# 计算 CPU 限制
# CPU限制 = quota / period
```

### 2. pprof 分析

```bash
# 启用 pprof
go test -bench=. -cpuprofile=cpu.prof

# 查看调度情况
go tool pprof cpu.prof
(pprof) top10
(pprof) list runtime.schedule
```

### 3. trace 分析

```bash
# 生成 trace
go test -bench=. -trace=trace.out

# 查看 trace
go tool trace trace.out
```

---

## 📚 相关资源

- [容器感知调度文档](../../02-容器感知调度.md)
- [Go 1.25 Release Notes](https://golang.org/doc/go1.25)
- [Docker CPU 限制](https://docs.docker.com/config/containers/resource_constraints/)
- [Kubernetes CPU 管理](https://kubernetes.io/docs/tasks/configure-pod-container/assign-cpu-resource/)
- [uber-go/automaxprocs](https://github.com/uber-go/automaxprocs) - Go 1.24 替代方案

---

**示例维护**: AI Assistant  
**最后更新**: 2025-10-18  
**反馈**: 提交 Issue 或 PR
