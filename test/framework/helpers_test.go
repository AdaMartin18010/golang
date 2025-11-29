package framework

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestTestContext(t *testing.T) {
	tc := NewTestContext(t)
	defer tc.DeferCleanup()

	// 测试添加清理函数
	cleaned := false
	tc.AddCleanup(func() {
		cleaned = true
	})

	// 执行清理
	tc.CleanupAll()

	if !cleaned {
		t.Error("Cleanup function should be called")
	}
}

func TestTestContext_Assertions(t *testing.T) {
	tc := NewTestContext(t)
	defer tc.DeferCleanup()

	// 测试断言
	tc.AssertTrue(true, "should be true")
	tc.AssertFalse(false, "should be false")
	tc.AssertEqual(1, 1, "should be equal")
	tc.AssertNotNil(tc, "should not be nil")
}

func TestRetryHelper(t *testing.T) {
	helper := NewRetryHelper(3, 100*time.Millisecond)

	attempts := 0
	err := helper.Retry(func() error {
		attempts++
		if attempts < 2 {
			return fmt.Errorf("attempt %d failed", attempts)
		}
		return nil
	})

	if err != nil {
		t.Errorf("Retry should succeed, got error: %v", err)
	}

	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestRetryHelper_MaxAttempts(t *testing.T) {
	helper := NewRetryHelper(3, 10*time.Millisecond)

	attempts := 0
	err := helper.Retry(func() error {
		attempts++
		return fmt.Errorf("always fail")
	})

	if err == nil {
		t.Error("Retry should fail after max attempts")
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestEnvironmentHelper(t *testing.T) {
	helper := NewEnvironmentHelper()
	defer helper.Restore()

	// 设置环境变量
	helper.SetEnv("TEST_VAR", "test_value")

	// 验证设置成功
	if os.Getenv("TEST_VAR") != "test_value" {
		t.Error("Environment variable should be set")
	}

	// 恢复
	helper.Restore()

	// 验证已恢复
	if os.Getenv("TEST_VAR") != "" {
		t.Error("Environment variable should be restored")
	}
}

func TestTestDataHelper(t *testing.T) {
	helper := NewTestDataHelper()

	// 设置数据
	helper.Set("string", "value")
	helper.Set("int", 42)

	// 获取数据
	str, ok := helper.GetString("string")
	if !ok || str != "value" {
		t.Errorf("Expected 'value', got '%s'", str)
	}

	i, ok := helper.GetInt("int")
	if !ok || i != 42 {
		t.Errorf("Expected 42, got %d", i)
	}
}

func TestMockHelper(t *testing.T) {
	helper := NewMockHelper()

	// 注册 Mock
	mock := "mock_value"
	helper.RegisterMock("test", mock)

	// 获取 Mock
	value, ok := helper.GetMock("test")
	if !ok {
		t.Error("Mock should be found")
	}

	if value != mock {
		t.Errorf("Expected '%s', got '%v'", mock, value)
	}
}
