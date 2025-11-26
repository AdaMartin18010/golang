package health

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestHealthChecker_Register(t *testing.T) {
	hc := NewHealthChecker()
	check := NewSimpleCheck("test", func(ctx context.Context) error {
		return nil
	})

	hc.Register(check)

	results := hc.Check(context.Background())
	if _, exists := results["test"]; !exists {
		t.Error("Expected check 'test' to be registered")
	}
}

func TestHealthChecker_Unregister(t *testing.T) {
	hc := NewHealthChecker()
	check := NewSimpleCheck("test", func(ctx context.Context) error {
		return nil
	})

	hc.Register(check)
	err := hc.Unregister("test")
	if err != nil {
		t.Fatalf("Failed to unregister check: %v", err)
	}

	results := hc.Check(context.Background())
	if _, exists := results["test"]; exists {
		t.Error("Expected check 'test' to be unregistered")
	}
}

func TestHealthChecker_OverallStatus(t *testing.T) {
	hc := NewHealthChecker()

	// 添加健康检查
	hc.Register(NewSimpleCheck("healthy", func(ctx context.Context) error {
		return nil
	}))

	status := hc.OverallStatus(context.Background())
	if status != StatusHealthy {
		t.Errorf("Expected status 'healthy', got '%s'", status)
	}

	// 添加不健康检查
	hc.Register(NewSimpleCheck("unhealthy", func(ctx context.Context) error {
		return errors.New("check failed")
	}))

	status = hc.OverallStatus(context.Background())
	if status != StatusUnhealthy {
		t.Errorf("Expected status 'unhealthy', got '%s'", status)
	}
}

func TestSimpleCheck_Check(t *testing.T) {
	check := NewSimpleCheck("test", func(ctx context.Context) error {
		return nil
	})

	result := check.Check(context.Background())
	if result.Status != StatusHealthy {
		t.Errorf("Expected status 'healthy', got '%s'", result.Status)
	}

	// 测试失败情况
	check = NewSimpleCheck("test", func(ctx context.Context) error {
		return errors.New("check failed")
	})

	result = check.Check(context.Background())
	if result.Status != StatusUnhealthy {
		t.Errorf("Expected status 'unhealthy', got '%s'", result.Status)
	}
}

func TestTimeoutCheck_Check(t *testing.T) {
	slowCheck := NewSimpleCheck("slow", func(ctx context.Context) error {
		time.Sleep(2 * time.Second)
		return nil
	})

	timeoutCheck := NewTimeoutCheck("timeout", 100*time.Millisecond, slowCheck)
	result := timeoutCheck.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("Expected status 'unhealthy' due to timeout, got '%s'", result.Status)
	}
}

func TestPeriodicCheck_Check(t *testing.T) {
	callCount := 0
	check := NewSimpleCheck("test", func(ctx context.Context) error {
		callCount++
		return nil
	})

	periodicCheck := NewPeriodicCheck("periodic", 1*time.Second, check)

	// 第一次调用应该执行检查
	result1 := periodicCheck.Check(context.Background())
	if result1.Status != StatusHealthy {
		t.Errorf("Expected status 'healthy', got '%s'", result1.Status)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// 立即再次调用应该使用缓存
	result2 := periodicCheck.Check(context.Background())
	if result2.Status != StatusHealthy {
		t.Errorf("Expected status 'healthy', got '%s'", result2.Status)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call (cached), got %d", callCount)
	}
}

func TestAggregateCheck_Check(t *testing.T) {
	healthyCheck := NewSimpleCheck("healthy", func(ctx context.Context) error {
		return nil
	})

	unhealthyCheck := NewSimpleCheck("unhealthy", func(ctx context.Context) error {
		return errors.New("check failed")
	})

	aggregateCheck := NewAggregateCheck("aggregate", healthyCheck, unhealthyCheck)
	result := aggregateCheck.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("Expected status 'unhealthy', got '%s'", result.Status)
	}
}
