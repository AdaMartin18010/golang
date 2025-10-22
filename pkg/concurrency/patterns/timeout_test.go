package patterns

import (
	"errors"
	"testing"
	"time"
)

func TestWithTimeoutFunc(t *testing.T) {
	// 成功案例
	result, err := WithTimeoutFunc(100*time.Millisecond, func() (interface{}, error) {
		return "success", nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got %v", result)
	}

	// 超时案例
	_, err = WithTimeoutFunc(10*time.Millisecond, func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return "should timeout", nil
	})
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestWithDeadlineFunc(t *testing.T) {
	// 成功案例
	deadline := time.Now().Add(100 * time.Millisecond)
	result, err := WithDeadlineFunc(deadline, func() (interface{}, error) {
		return "success", nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got %v", result)
	}

	// 已过期的deadline
	pastDeadline := time.Now().Add(-1 * time.Second)
	_, err = WithDeadlineFunc(pastDeadline, func() (interface{}, error) {
		return nil, nil
	})
	if err == nil {
		t.Error("Expected deadline passed error")
	}
}

func TestTimeoutRetry(t *testing.T) {
	attempts := 0

	// 第2次尝试成功
	err := TimeoutRetry(3, 100*time.Millisecond, func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error after retry, got %v", err)
	}
	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestTimeoutRetryAllFailed(t *testing.T) {
	attempts := 0

	err := TimeoutRetry(3, 100*time.Millisecond, func() error {
		attempts++
		return errors.New("persistent error")
	})

	if err == nil {
		t.Error("Expected error after all retries failed")
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestBatchWithTimeout(t *testing.T) {
	tasks := []func() error{
		func() error { time.Sleep(10 * time.Millisecond); return nil },
		func() error { time.Sleep(10 * time.Millisecond); return nil },
		func() error { time.Sleep(10 * time.Millisecond); return errors.New("task failed") },
		func() error { time.Sleep(200 * time.Millisecond); return nil }, // 应该超时
	}

	results := BatchWithTimeout(tasks, 50*time.Millisecond)

	if len(results) != 4 {
		t.Errorf("Expected 4 results, got %d", len(results))
	}

	// 前两个应该成功
	if results[0] != nil {
		t.Errorf("Task 0 should succeed, got %v", results[0])
	}
	if results[1] != nil {
		t.Errorf("Task 1 should succeed, got %v", results[1])
	}

	// 第3个应该失败
	if results[2] == nil {
		t.Error("Task 2 should fail")
	}

	// 第4个应该超时
	if results[3] == nil {
		t.Error("Task 3 should timeout")
	}
}

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)

	// 初始状态应该是closed
	if cb.GetState() != "closed" {
		t.Errorf("Initial state should be 'closed', got %s", cb.GetState())
	}

	// 执行成功的任务
	err := cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Successful task should not error, got %v", err)
	}

	// 执行3次失败的任务
	for i := 0; i < 3; i++ {
		cb.Execute(func() error {
			return errors.New("task failed")
		})
	}

	// 断路器应该打开
	if cb.GetState() != "open" {
		t.Errorf("Circuit breaker should be 'open' after max failures, got %s", cb.GetState())
	}

	// 尝试执行任务应该被拒绝
	err = cb.Execute(func() error {
		return nil
	})
	if err == nil {
		t.Error("Execute should be rejected when circuit breaker is open")
	}

	// 等待重置
	time.Sleep(150 * time.Millisecond)

	// 状态应该变为half-open
	if cb.GetState() != "half-open" {
		t.Errorf("Circuit breaker should be 'half-open' after timeout, got %s", cb.GetState())
	}

	// 执行成功的任务
	err = cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Task should succeed in half-open state, got %v", err)
	}

	// 状态应该变回closed
	if cb.GetState() != "closed" {
		t.Errorf("Circuit breaker should be 'closed' after successful execution, got %s", cb.GetState())
	}
}

func BenchmarkWithTimeoutFunc(b *testing.B) {
	task := func() (interface{}, error) {
		return "result", nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithTimeoutFunc(100*time.Millisecond, task)
	}
}

func BenchmarkCircuitBreaker(b *testing.B) {
	cb := NewCircuitBreaker(100, 1*time.Second)
	task := func() error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cb.Execute(task)
	}
}
