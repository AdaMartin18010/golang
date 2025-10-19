# Kubernetes Operator 示例 - 依赖说明

## ⚠️ 重要提示

这个示例需要 Kubernetes 相关的外部依赖库。这些依赖体积较大（约100MB+），仅供学习和参考使用。

## 📦 依赖列表

- `k8s.io/api` - Kubernetes API 类型定义
- `k8s.io/apimachinery` - Kubernetes API 机制
- `k8s.io/client-go` - Kubernetes 客户端库
- `sigs.k8s.io/controller-runtime` - Controller Runtime 框架
- `github.com/prometheus/client_golang` - Prometheus 客户端

## 🚀 使用方法

### 选项1: 查看代码（推荐）

直接查看代码了解 Kubernetes Operator 的实现模式和最佳实践。

### 选项2: 运行示例

如果需要实际运行这个示例：

```bash
cd examples/modern-features/08-cloud-native/kubernetes-operator

# 下载依赖（这将下载约100MB+的依赖）
go mod download

# 编译
go build ./...

# 运行测试
go test ./...
```

## 📚 学习重点

这个示例展示了：

1. **CRD（自定义资源定义）** 的设计
2. **Controller 模式** 的实现
3. **Reconcile 循环** 的编写
4. **事件记录** 和 **指标收集**
5. **资源管理** 的最佳实践

## 💡 提示

- 这个示例是**演示性质**的，展示了 Operator 的核心概念
- 生产环境的 Operator 需要更多的错误处理和边界情况处理
- 建议先学习 Kubernetes 基础知识后再深入 Operator 开发

## 🔗 相关资源

- [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)
- [Kubebuilder](https://book.kubebuilder.io/)

