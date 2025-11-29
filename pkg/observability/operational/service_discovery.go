package operational

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
	Register(ctx context.Context, service ServiceInfo) error
	Deregister(ctx context.Context, serviceID string) error
	HealthCheck(ctx context.Context, serviceID string) error
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Address     string            `json:"address"`
	Port        int               `json:"port"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	HealthCheck HealthCheckConfig `json:"health_check"`
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Path     string        `json:"path"`
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
}

// ServiceRegistry 服务注册表（简化实现）
type ServiceRegistry struct {
	services map[string]ServiceInfo
	mu       sync.RWMutex
}

// NewServiceRegistry 创建服务注册表
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]ServiceInfo),
	}
}

// Register 注册服务
func (sr *ServiceRegistry) Register(ctx context.Context, service ServiceInfo) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if service.ID == "" {
		return fmt.Errorf("service ID is required")
	}

	sr.services[service.ID] = service
	return nil
}

// Deregister 注销服务
func (sr *ServiceRegistry) Deregister(ctx context.Context, serviceID string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	delete(sr.services, serviceID)
	return nil
}

// GetService 获取服务信息
func (sr *ServiceRegistry) GetService(serviceID string) (ServiceInfo, bool) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	service, exists := sr.services[serviceID]
	return service, exists
}

// ListServices 列出所有服务
func (sr *ServiceRegistry) ListServices() []ServiceInfo {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	services := make([]ServiceInfo, 0, len(sr.services))
	for _, service := range sr.services {
		services = append(services, service)
	}
	return services
}
