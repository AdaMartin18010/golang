package registry

import (
	"context"
	"testing"
	"time"
)

func TestInMemoryRegistry_Register(t *testing.T) {
	registry := NewInMemoryRegistry()
	service := &Service{
		ID:      "service-1",
		Name:    "user-service",
		Address: "localhost",
		Port:    8080,
	}

	err := registry.Register(context.Background(), service)
	if err != nil {
		t.Fatalf("Failed to register service: %v", err)
	}

	// 验证服务已注册
	registered, err := registry.GetService(context.Background(), "service-1")
	if err != nil {
		t.Fatalf("Failed to get service: %v", err)
	}

	if registered.Name != "user-service" {
		t.Errorf("Expected service name 'user-service', got '%s'", registered.Name)
	}
}

func TestInMemoryRegistry_Deregister(t *testing.T) {
	registry := NewInMemoryRegistry()
	service := &Service{
		ID:   "service-1",
		Name: "user-service",
	}

	registry.Register(context.Background(), service)
	err := registry.Deregister(context.Background(), "service-1")
	if err != nil {
		t.Fatalf("Failed to deregister service: %v", err)
	}

	// 验证服务已注销
	_, err = registry.GetService(context.Background(), "service-1")
	if err != ErrServiceNotFound {
		t.Errorf("Expected ErrServiceNotFound, got %v", err)
	}
}

func TestInMemoryRegistry_ListServices(t *testing.T) {
	registry := NewInMemoryRegistry()

	services := []*Service{
		{ID: "service-1", Name: "user-service"},
		{ID: "service-2", Name: "user-service"},
		{ID: "service-3", Name: "order-service"},
	}

	for _, service := range services {
		registry.Register(context.Background(), service)
	}

	// 列出所有服务
	all, err := registry.ListServices(context.Background(), "")
	if err != nil {
		t.Fatalf("Failed to list services: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Expected 3 services, got %d", len(all))
	}

	// 列出特定名称的服务
	userServices, err := registry.ListServices(context.Background(), "user-service")
	if err != nil {
		t.Fatalf("Failed to list user services: %v", err)
	}
	if len(userServices) != 2 {
		t.Errorf("Expected 2 user services, got %d", len(userServices))
	}
}

func TestInMemoryRegistry_Watch(t *testing.T) {
	registry := NewInMemoryRegistry()

	ch, err := registry.Watch(context.Background(), "user-service")
	if err != nil {
		t.Fatalf("Failed to watch services: %v", err)
	}

	// 注册服务
	service := &Service{
		ID:   "service-1",
		Name: "user-service",
	}
	registry.Register(context.Background(), service)

	// 等待通知
	select {
	case services := <-ch:
		if len(services) != 1 {
			t.Errorf("Expected 1 service, got %d", len(services))
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for service update")
	}
}

func TestInMemoryRegistry_CleanupExpiredServices(t *testing.T) {
	registry := NewInMemoryRegistry()

	service := &Service{
		ID:      "service-1",
		Name:    "user-service",
		TTL:     100 * time.Millisecond,
		LastSeen: time.Now().Add(-200 * time.Millisecond),
	}

	registry.Register(context.Background(), service)
	registry.CleanupExpiredServices(context.Background(), 0)

	// 验证服务已被清理
	_, err := registry.GetService(context.Background(), "service-1")
	if err != ErrServiceNotFound {
		t.Errorf("Expected ErrServiceNotFound, got %v", err)
	}
}
