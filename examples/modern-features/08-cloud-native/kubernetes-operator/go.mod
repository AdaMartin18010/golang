module kubernetes-operator-example

go 1.25

// 注意：这个示例需要Kubernetes相关的依赖
// 这些依赖体积较大，在实际使用前需要运行：
// go mod download

require (
	github.com/prometheus/client_golang v1.19.0
	k8s.io/api v0.29.0
	k8s.io/apimachinery v0.29.0
	k8s.io/client-go v0.29.0
	sigs.k8s.io/controller-runtime v0.17.0
)

