package registry

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	// ErrServiceNotFound 服务未找到
	ErrServiceNotFound = errors.New("service not found")
	// ErrServiceAlreadyExists 服务已存在
	ErrServiceAlreadyExists = errors.New("service already exists")
)

// Service 服务实例
type Service struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Tags     []string          `json:"tags,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	TTL      time.Duration     `json:"ttl,omitempty"`
	LastSeen time.Time         `json:"last_seen"`
}

// Registry 服务注册中心接口
type Registry interface {
	// Register 注册服务
	Register(ctx context.Context, service *Service) error
	// Deregister 注销服务
	Deregister(ctx context.Context, serviceID string) error
	// GetService 获取服务
	GetService(ctx context.Context, serviceID string) (*Service, error)
	// ListServices 列出所有服务
	ListServices(ctx context.Context, name string) ([]*Service, error)
	// Watch 监听服务变化
	Watch(ctx context.Context, name string) (<-chan []*Service, error)
	// Health 健康检查
	Health(ctx context.Context) error
}

// InMemoryRegistry 内存服务注册中心
type InMemoryRegistry struct {
	services map[string]*Service
	watchers map[string][]chan []*Service
	mu       sync.RWMutex
}

// NewInMemoryRegistry 创建内存服务注册中心
func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		services: make(map[string]*Service),
		watchers: make(map[string][]chan []*Service),
	}
}

// Register 注册服务
func (r *InMemoryRegistry) Register(ctx context.Context, service *Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if service.ID == "" {
		return errors.New("service ID is required")
	}
	if service.Name == "" {
		return errors.New("service name is required")
	}

	service.LastSeen = time.Now()
	r.services[service.ID] = service

	// 通知监听者
	r.notifyWatchers(service.Name)

	return nil
}

// Deregister 注销服务
func (r *InMemoryRegistry) Deregister(ctx context.Context, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	service, exists := r.services[serviceID]
	if !exists {
		return ErrServiceNotFound
	}

	delete(r.services, serviceID)

	// 通知监听者
	r.notifyWatchers(service.Name)

	return nil
}

// GetService 获取服务
func (r *InMemoryRegistry) GetService(ctx context.Context, serviceID string) (*Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, exists := r.services[serviceID]
	if !exists {
		return nil, ErrServiceNotFound
	}

	return service, nil
}

// ListServices 列出所有服务
func (r *InMemoryRegistry) ListServices(ctx context.Context, name string) ([]*Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var services []*Service
	for _, service := range r.services {
		if name == "" || service.Name == name {
			services = append(services, service)
		}
	}

	return services, nil
}

// Watch 监听服务变化
func (r *InMemoryRegistry) Watch(ctx context.Context, name string) (<-chan []*Service, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ch := make(chan []*Service, 10)
	r.watchers[name] = append(r.watchers[name], ch)

	// 立即发送当前服务列表
	go func() {
		services, _ := r.ListServices(ctx, name)
		select {
		case ch <- services:
		case <-ctx.Done():
		}
	}()

	return ch, nil
}

// Health 健康检查
func (r *InMemoryRegistry) Health(ctx context.Context) error {
	return nil
}

// notifyWatchers 通知监听者
func (r *InMemoryRegistry) notifyWatchers(name string) {
	watchers := r.watchers[name]
	services, _ := r.ListServices(context.Background(), name)

	for _, ch := range watchers {
		select {
		case ch <- services:
		default:
			// 如果通道已满，跳过
		}
	}
}

// CleanupExpiredServices 清理过期服务
func (r *InMemoryRegistry) CleanupExpiredServices(ctx context.Context, maxAge time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for id, service := range r.services {
		if service.TTL > 0 && now.Sub(service.LastSeen) > service.TTL {
			delete(r.services, id)
			r.notifyWatchers(service.Name)
		} else if maxAge > 0 && now.Sub(service.LastSeen) > maxAge {
			delete(r.services, id)
			r.notifyWatchers(service.Name)
		}
	}
}
