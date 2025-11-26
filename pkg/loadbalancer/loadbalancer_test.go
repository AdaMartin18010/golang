package loadbalancer

import (
	"context"
	"testing"

	"github.com/yourusername/golang/pkg/registry"
)

func TestRoundRobin_Select(t *testing.T) {
	lb := NewRoundRobin()
	services := []*registry.Service{
		{ID: "service-1", Name: "test"},
		{ID: "service-2", Name: "test"},
		{ID: "service-3", Name: "test"},
	}

	// 测试轮询
	selected1, _ := lb.Select(context.Background(), services)
	selected2, _ := lb.Select(context.Background(), services)
	selected3, _ := lb.Select(context.Background(), services)
	selected4, _ := lb.Select(context.Background(), services)

	if selected1.ID != "service-1" {
		t.Errorf("Expected service-1, got %s", selected1.ID)
	}
	if selected2.ID != "service-2" {
		t.Errorf("Expected service-2, got %s", selected2.ID)
	}
	if selected3.ID != "service-3" {
		t.Errorf("Expected service-3, got %s", selected3.ID)
	}
	if selected4.ID != "service-1" {
		t.Errorf("Expected service-1 (wrapped), got %s", selected4.ID)
	}
}

func TestRandom_Select(t *testing.T) {
	lb := NewRandom()
	services := []*registry.Service{
		{ID: "service-1", Name: "test"},
		{ID: "service-2", Name: "test"},
	}

	// 测试随机选择（至少应该能选择到服务）
	selected, err := lb.Select(context.Background(), services)
	if err != nil {
		t.Fatalf("Failed to select service: %v", err)
	}
	if selected == nil {
		t.Error("Expected a service, got nil")
	}
}

func TestLeastConnections_Select(t *testing.T) {
	lb := NewLeastConnections()
	services := []*registry.Service{
		{ID: "service-1", Name: "test"},
		{ID: "service-2", Name: "test"},
	}

	// 第一次选择应该选择service-1
	selected1, _ := lb.Select(context.Background(), services)
	if selected1.ID != "service-1" {
		t.Errorf("Expected service-1, got %s", selected1.ID)
	}

	// 第二次选择应该选择service-2（因为service-1连接数更多）
	selected2, _ := lb.Select(context.Background(), services)
	if selected2.ID != "service-2" {
		t.Errorf("Expected service-2, got %s", selected2.ID)
	}

	// 释放service-1的连接
	lb.Release("service-1")

	// 再次选择应该选择service-1（因为连接数更少）
	selected3, _ := lb.Select(context.Background(), services)
	if selected3.ID != "service-1" {
		t.Errorf("Expected service-1, got %s", selected3.ID)
	}
}

func TestServiceSelector_Select(t *testing.T) {
	reg := registry.NewInMemoryRegistry()
	lb := NewRoundRobin()
	selector := NewServiceSelector(reg, lb, "test-service")

	// 注册服务
	services := []*registry.Service{
		{ID: "service-1", Name: "test-service"},
		{ID: "service-2", Name: "test-service"},
	}

	for _, service := range services {
		reg.Register(context.Background(), service)
	}

	// 选择服务
	selected, err := selector.Select(context.Background())
	if err != nil {
		t.Fatalf("Failed to select service: %v", err)
	}

	if selected == nil {
		t.Error("Expected a service, got nil")
	}

	if selected.Name != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", selected.Name)
	}
}

func TestLoadBalancer_NoServices(t *testing.T) {
	lb := NewRoundRobin()
	_, err := lb.Select(context.Background(), []*registry.Service{})
	if err != ErrNoServices {
		t.Errorf("Expected ErrNoServices, got %v", err)
	}
}
