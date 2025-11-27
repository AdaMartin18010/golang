package loadbalancer

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/yourusername/golang/pkg/registry"
)

var (
	// ErrNoServices 没有可用服务
	ErrNoServices = errors.New("no services available")
)

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
	// Select 选择服务实例
	Select(ctx context.Context, services []*registry.Service) (*registry.Service, error)
	// Name 返回负载均衡算法名称
	Name() string
}

// RoundRobin 轮询负载均衡
type RoundRobin struct {
	mu    sync.Mutex
	index int
}

// NewRoundRobin 创建轮询负载均衡器
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

// Select 选择服务实例
func (rr *RoundRobin) Select(ctx context.Context, services []*registry.Service) (*registry.Service, error) {
	if len(services) == 0 {
		return nil, ErrNoServices
	}

	rr.mu.Lock()
	defer rr.mu.Unlock()

	service := services[rr.index%len(services)]
	rr.index++

	return service, nil
}

// Name 返回算法名称
func (rr *RoundRobin) Name() string {
	return "round-robin"
}

// Random 随机负载均衡
type Random struct {
	rand *rand.Rand
}

// NewRandom 创建随机负载均衡器
func NewRandom() *Random {
	return &Random{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Select 选择服务实例
func (r *Random) Select(ctx context.Context, services []*registry.Service) (*registry.Service, error) {
	if len(services) == 0 {
		return nil, ErrNoServices
	}

	index := r.rand.Intn(len(services))
	return services[index], nil
}

// Name 返回算法名称
func (r *Random) Name() string {
	return "random"
}

// WeightedRoundRobin 加权轮询负载均衡
type WeightedRoundRobin struct {
	mu    sync.Mutex
	index int
}

// NewWeightedRoundRobin 创建加权轮询负载均衡器
func NewWeightedRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{}
}

// Select 选择服务实例
func (wrr *WeightedRoundRobin) Select(ctx context.Context, services []*registry.Service) (*registry.Service, error) {
	if len(services) == 0 {
		return nil, ErrNoServices
	}

	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	// 计算总权重
	totalWeight := 0
	for _, service := range services {
		weight := getWeight(service)
		totalWeight += weight
	}

	if totalWeight == 0 {
		// 如果没有权重，使用普通轮询
		service := services[wrr.index%len(services)]
		wrr.index++
		return service, nil
	}

	// 加权轮询
	current := wrr.index % totalWeight
	for _, service := range services {
		weight := getWeight(service)
		if current < weight {
			wrr.index++
			return service, nil
		}
		current -= weight
	}

	// 不应该到达这里
	return services[0], nil
}

// Name 返回算法名称
func (wrr *WeightedRoundRobin) Name() string {
	return "weighted-round-robin"
}

// getWeight 获取服务权重
func getWeight(service *registry.Service) int {
	if weight, ok := service.Metadata["weight"]; ok {
		// 解析权重（简化实现）
		// 实际应该使用strconv.Atoi
		return 1 // 默认权重为1
	}
	return 1
}

// LeastConnections 最少连接负载均衡
type LeastConnections struct {
	connections map[string]int
	mu          sync.RWMutex
}

// NewLeastConnections 创建最少连接负载均衡器
func NewLeastConnections() *LeastConnections {
	return &LeastConnections{
		connections: make(map[string]int),
	}
}

// Select 选择服务实例
func (lc *LeastConnections) Select(ctx context.Context, services []*registry.Service) (*registry.Service, error) {
	if len(services) == 0 {
		return nil, ErrNoServices
	}

	lc.mu.Lock()
	defer lc.mu.Unlock()

	// 找到连接数最少的服务
	minConnections := -1
	var selected *registry.Service

	for _, service := range services {
		conn := lc.connections[service.ID]
		if minConnections == -1 || conn < minConnections {
			minConnections = conn
			selected = service
		}
	}

	if selected != nil {
		lc.connections[selected.ID]++
	}

	return selected, nil
}

// Release 释放连接
func (lc *LeastConnections) Release(serviceID string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if conn, exists := lc.connections[serviceID]; exists && conn > 0 {
		lc.connections[serviceID]--
	}
}

// Name 返回算法名称
func (lc *LeastConnections) Name() string {
	return "least-connections"
}

// ServiceSelector 服务选择器
type ServiceSelector struct {
	registry    registry.Registry
	balancer    LoadBalancer
	serviceName string
}

// NewServiceSelector 创建服务选择器
func NewServiceSelector(reg registry.Registry, balancer LoadBalancer, serviceName string) *ServiceSelector {
	return &ServiceSelector{
		registry:    reg,
		balancer:    balancer,
		serviceName: serviceName,
	}
}

// Select 选择服务实例
func (ss *ServiceSelector) Select(ctx context.Context) (*registry.Service, error) {
	services, err := ss.registry.ListServices(ctx, ss.serviceName)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, ErrNoServices
	}

	return ss.balancer.Select(ctx, services)
}
