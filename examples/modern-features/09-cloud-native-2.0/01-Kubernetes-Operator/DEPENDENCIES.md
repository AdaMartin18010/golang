# Kubernetes Operator 高级示例 - 依赖说明

## ⚠️ 重要提示

这个示例需要 Kubernetes 相关的外部依赖库。这些依赖体积较大（约100MB+），仅供学习和参考使用。

## 📦 依赖列表

- `k8s.io/api` - Kubernetes API 类型定义（包含apps/v1, autoscaling/v2等）
- `k8s.io/apimachinery` - Kubernetes API 机制
- `k8s.io/client-go` - Kubernetes 客户端库
- `k8s.io/utils` - Kubernetes 工具函数
- `sigs.k8s.io/controller-runtime` - Controller Runtime 框架
- `github.com/prometheus/client_golang` - Prometheus 客户端

## 🚀 使用方法

### 选项1: 查看代码（推荐）

直接查看代码了解高级 Kubernetes Operator 的实现模式，包括：
- 完整的应用生命周期管理
- 自动扩缩容（HPA）
- 资源管理和调度
- 高级健康检查
- 存储和网络配置

### 选项2: 运行示例

如果需要实际运行这个示例：

```bash
cd examples/modern-features/09-cloud-native-2.0/01-Kubernetes-Operator

# 下载依赖（这将下载约100MB+的依赖）
go mod download

# 编译
go build ./...

# 运行测试
go test ./...
```

## 📚 学习重点

这个高级示例展示了：

1. **完整的应用规范** - 包括副本、资源、健康检查等
2. **自动扩缩容** - HPA（Horizontal Pod Autoscaler）集成
3. **存储管理** - PVC（Persistent Volume Claim）管理
4. **网络配置** - Service和LoadBalancer配置
5. **安全上下文** - SecurityContext和RBAC
6. **资源管理器** - 统一的资源创建和更新逻辑
7. **事件记录** - 详细的事件记录和监控

## 💡 提示

- 这个示例是**演示性质**的，展示了生产级 Operator 的核心概念
- 实际生产环境需要更多的错误处理、重试逻辑和边界情况处理
- 建议先熟悉基础 Operator 示例（08-cloud-native）后再学习这个高级示例

## 🔗 相关资源

- [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)
- [Operator Best Practices](https://sdk.operatorframework.io/docs/best-practices/)
- [Kubebuilder Book](https://book.kubebuilder.io/)

