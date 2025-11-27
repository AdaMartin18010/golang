package control

import (
	"testing"
	"time"
)

func TestFeatureController_EnableDisable(t *testing.T) {
	controller := NewFeatureController()

	controller.Register("feature-a", "Feature A", false, nil)

	if controller.IsEnabled("feature-a") {
		t.Error("Feature should be disabled initially")
	}

	if err := controller.Enable("feature-a"); err != nil {
		t.Fatalf("Failed to enable feature: %v", err)
	}

	if !controller.IsEnabled("feature-a") {
		t.Error("Feature should be enabled")
	}

	if err := controller.Disable("feature-a"); err != nil {
		t.Fatalf("Failed to disable feature: %v", err)
	}

	if controller.IsEnabled("feature-a") {
		t.Error("Feature should be disabled")
	}
}

func TestFeatureController_Config(t *testing.T) {
	controller := NewFeatureController()

	config := map[string]interface{}{
		"max_requests": 100,
	}
	controller.Register("feature-a", "Feature A", true, config)

	retrieved, err := controller.GetConfig("feature-a")
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if retrieved == nil {
		t.Error("Expected non-nil config")
	}

	// 更新配置
	newConfig := map[string]interface{}{
		"max_requests": 200,
	}
	if err := controller.SetConfig("feature-a", newConfig); err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}
}

func TestFeatureController_Watch(t *testing.T) {
	controller := NewFeatureController()

	controller.Register("feature-a", "Feature A", true, nil)

	called := false
	controller.Watch("feature-a", func(config interface{}) {
		called = true
	})

	// 更新配置应该触发回调
	controller.SetConfig("feature-a", map[string]interface{}{"test": true})

	// 等待异步回调
	time.Sleep(100 * time.Millisecond)

	if !called {
		t.Error("Watch callback should be called")
	}
}

func TestRateController_Allow(t *testing.T) {
	controller := NewRateController()

	controller.SetRateLimit("api", 5.0, time.Second)

	// 应该允许前 5 次
	allowed := 0
	for i := 0; i < 10; i++ {
		if controller.Allow("api") {
			allowed++
		}
	}

	if allowed < 5 {
		t.Errorf("Expected at least 5 allowed, got %d", allowed)
	}
}

func TestCircuitController(t *testing.T) {
	controller := NewCircuitController()

	controller.RegisterCircuit("external-api", 3, 2, 1*time.Second)

	// 初始状态应该是关闭的
	if controller.IsOpen("external-api") {
		t.Error("Circuit should be closed initially")
	}

	// 记录失败，应该打开熔断器
	controller.RecordFailure("external-api")
	controller.RecordFailure("external-api")
	controller.RecordFailure("external-api")

	if !controller.IsOpen("external-api") {
		t.Error("Circuit should be open after failures")
	}

	// 等待超时后应该进入半开状态
	time.Sleep(1100 * time.Millisecond)

	// 记录成功，应该关闭熔断器
	controller.RecordSuccess("external-api")
	controller.RecordSuccess("external-api")

	if controller.IsOpen("external-api") {
		t.Error("Circuit should be closed after successes")
	}
}

