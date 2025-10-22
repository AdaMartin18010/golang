package core

import (
	"errors"
	"testing"
)

// TestNewError 测试创建新错误
func TestNewError(t *testing.T) {
	err := NewError(ErrorCodeInternal, "test error")

	if err.Code != ErrorCodeInternal {
		t.Errorf("Expected code %s, got %s", ErrorCodeInternal, err.Code)
	}
	if err.Message != "test error" {
		t.Errorf("Expected message 'test error', got '%s'", err.Message)
	}
	if err.Retryable {
		t.Error("Expected non-retryable error")
	}
}

// TestNewRetryableError 测试创建可重试错误
func TestNewRetryableError(t *testing.T) {
	err := NewRetryableError(ErrorCodeTimeout, "timeout error")

	if !err.Retryable {
		t.Error("Expected retryable error")
	}
	if err.Code != ErrorCodeTimeout {
		t.Errorf("Expected code %s, got %s", ErrorCodeTimeout, err.Code)
	}
}

// TestWrapError 测试包装错误
func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapError(originalErr, ErrorCodeProcessingFailed, "processing failed")

	if wrappedErr == nil {
		t.Fatal("Expected non-nil wrapped error")
	}
	if wrappedErr.Code != ErrorCodeProcessingFailed {
		t.Errorf("Expected code %s, got %s", ErrorCodeProcessingFailed, wrappedErr.Code)
	}
	if wrappedErr.Cause != originalErr {
		t.Error("Expected cause to be set to original error")
	}
}

// TestWrapErrorNil 测试包装nil错误
func TestWrapErrorNil(t *testing.T) {
	wrappedErr := WrapError(nil, ErrorCodeInternal, "test")

	if wrappedErr != nil {
		t.Error("Expected nil when wrapping nil error")
	}
}

// TestWrapAgentError 测试包装AgentError
func TestWrapAgentError(t *testing.T) {
	innerErr := NewError(ErrorCodeInvalidInput, "inner error")
	outerErr := WrapError(innerErr, ErrorCodeProcessingFailed, "outer error")

	if outerErr.Cause != innerErr {
		t.Error("Expected cause to be inner error")
	}
}

// TestAgentErrorError 测试Error()方法
func TestAgentErrorError(t *testing.T) {
	// 简单错误
	err1 := NewError(ErrorCodeInternal, "simple error")
	if err1.Error() == "" {
		t.Error("Error message should not be empty")
	}

	// 带详情的错误
	err2 := NewError(ErrorCodeInternal, "error with details").WithDetails("more info")
	if err2.Error() == "" {
		t.Error("Error message should not be empty")
	}

	// 带cause的错误
	cause := errors.New("cause error")
	err3 := WrapError(cause, ErrorCodeInternal, "wrapper")
	if err3.Error() == "" {
		t.Error("Error message should not be empty")
	}
}

// TestAgentErrorUnwrap 测试Unwrap()方法
func TestAgentErrorUnwrap(t *testing.T) {
	cause := errors.New("cause")
	err := WrapError(cause, ErrorCodeInternal, "wrapper")

	unwrapped := errors.Unwrap(err)
	if unwrapped != cause {
		t.Error("Unwrap should return the cause")
	}
}

// TestAgentErrorIs 测试errors.Is()
func TestAgentErrorIs(t *testing.T) {
	err1 := NewError(ErrorCodeTimeout, "timeout")
	err2 := NewError(ErrorCodeTimeout, "another timeout")
	err3 := NewError(ErrorCodeInternal, "internal")

	if !errors.Is(err1, err2) {
		t.Error("Errors with same code should match")
	}
	if errors.Is(err1, err3) {
		t.Error("Errors with different codes should not match")
	}
}

// TestAgentErrorWithContext 测试添加上下文
func TestAgentErrorWithContext(t *testing.T) {
	err := NewError(ErrorCodeInternal, "test")
	err.WithContext("user_id", "12345")
	err.WithContext("request_id", "req-001")

	if len(err.Context) != 2 {
		t.Errorf("Expected 2 context items, got %d", len(err.Context))
	}
	if err.Context["user_id"] != "12345" {
		t.Error("Context not set correctly")
	}
}

// TestAgentErrorWithDetails 测试添加详情
func TestAgentErrorWithDetails(t *testing.T) {
	err := NewError(ErrorCodeInternal, "test").WithDetails("detailed information")

	if err.Details != "detailed information" {
		t.Errorf("Expected details 'detailed information', got '%s'", err.Details)
	}
}

// TestIsRetryable 测试判断可重试
func TestIsRetryable(t *testing.T) {
	retryable := NewRetryableError(ErrorCodeTimeout, "timeout")
	nonRetryable := NewError(ErrorCodeInternal, "internal")

	if !IsRetryable(retryable) {
		t.Error("Expected error to be retryable")
	}
	if IsRetryable(nonRetryable) {
		t.Error("Expected error to be non-retryable")
	}

	// 测试普通error
	regularErr := errors.New("regular error")
	if IsRetryable(regularErr) {
		t.Error("Regular error should not be retryable")
	}
}

// TestGetErrorCode 测试获取错误代码
func TestGetErrorCode(t *testing.T) {
	err := NewError(ErrorCodeTimeout, "timeout")
	code := GetErrorCode(err)

	if code != ErrorCodeTimeout {
		t.Errorf("Expected code %s, got %s", ErrorCodeTimeout, code)
	}

	// 测试普通error
	regularErr := errors.New("regular error")
	code = GetErrorCode(regularErr)
	if code != ErrorCodeInternal {
		t.Errorf("Regular error should return %s, got %s", ErrorCodeInternal, code)
	}
}

// TestGetErrorContext 测试获取错误上下文
func TestGetErrorContext(t *testing.T) {
	err := NewError(ErrorCodeInternal, "test").WithContext("key", "value")
	ctx := GetErrorContext(err)

	if ctx == nil {
		t.Fatal("Expected non-nil context")
	}
	if ctx["key"] != "value" {
		t.Error("Context not retrieved correctly")
	}

	// 测试普通error
	regularErr := errors.New("regular error")
	ctx = GetErrorContext(regularErr)
	if ctx != nil {
		t.Error("Regular error should have nil context")
	}
}

// TestPredefinedErrors 测试预定义错误
func TestPredefinedErrors(t *testing.T) {
	errors := []*AgentError{
		ErrAgentNotRunning,
		ErrAgentAlreadyRunning,
		ErrInvalidInput,
		ErrTimeout,
		ErrPluginNotFound,
	}

	for _, err := range errors {
		if err == nil {
			t.Error("Predefined error should not be nil")
		}
		if err.Code == "" {
			t.Error("Predefined error should have a code")
		}
		if err.Message == "" {
			t.Error("Predefined error should have a message")
		}
	}
}

// TestErrorCodes 测试所有错误代码
func TestErrorCodes(t *testing.T) {
	codes := []ErrorCode{
		ErrorCodeInternal,
		ErrorCodeTimeout,
		ErrorCodeCancelled,
		ErrorCodeInvalidState,
		ErrorCodeResourceExhausted,
		ErrorCodeInvalidConfig,
		ErrorCodeMissingConfig,
		ErrorCodePluginNotFound,
		ErrorCodePluginFailed,
		ErrorCodeProcessingFailed,
		ErrorCodeInvalidInput,
		ErrorCodeInvalidOutput,
		ErrorCodeDecisionFailed,
		ErrorCodeNoDecision,
		ErrorCodeLearningFailed,
		ErrorCodeInvalidExperience,
	}

	for _, code := range codes {
		if string(code) == "" {
			t.Error("Error code should not be empty")
		}
	}
}

// TestErrorTracker 测试错误跟踪器
func TestErrorTracker(t *testing.T) {
	tracker := NewErrorTracker()

	// 跟踪几个错误
	err1 := NewError(ErrorCodeTimeout, "timeout 1")
	err2 := NewError(ErrorCodeTimeout, "timeout 2")
	err3 := NewRetryableError(ErrorCodeProcessingFailed, "processing failed")

	tracker.Track(err1)
	tracker.Track(err2)
	tracker.Track(err3)

	stats := tracker.GetStats()

	if stats.TotalErrors != 3 {
		t.Errorf("Expected 3 total errors, got %d", stats.TotalErrors)
	}
	if stats.ErrorsByCode[ErrorCodeTimeout] != 2 {
		t.Errorf("Expected 2 timeout errors, got %d", stats.ErrorsByCode[ErrorCodeTimeout])
	}
	if stats.RetryableCount != 1 {
		t.Errorf("Expected 1 retryable error, got %d", stats.RetryableCount)
	}
}

// TestErrorTrackerNil 测试跟踪nil错误
func TestErrorTrackerNil(t *testing.T) {
	tracker := NewErrorTracker()
	tracker.Track(nil)

	stats := tracker.GetStats()
	if stats.TotalErrors != 0 {
		t.Error("Tracking nil should not increment counter")
	}
}

// TestErrorTrackerRegularError 测试跟踪普通错误
func TestErrorTrackerRegularError(t *testing.T) {
	tracker := NewErrorTracker()
	regularErr := errors.New("regular error")
	tracker.Track(regularErr)

	stats := tracker.GetStats()
	if stats.TotalErrors != 1 {
		t.Errorf("Expected 1 error, got %d", stats.TotalErrors)
	}
	if stats.ErrorsByCode[ErrorCodeInternal] != 1 {
		t.Error("Regular error should be counted as internal")
	}
}

// TestErrorTrackerReset 测试重置统计
func TestErrorTrackerReset(t *testing.T) {
	tracker := NewErrorTracker()

	tracker.Track(NewError(ErrorCodeInternal, "error"))
	tracker.Reset()

	stats := tracker.GetStats()
	if stats.TotalErrors != 0 {
		t.Error("Stats should be reset")
	}
	if len(stats.ErrorsByCode) != 0 {
		t.Error("ErrorsByCode should be empty after reset")
	}
}

// TestErrorTrackerLastError 测试最后一个错误
func TestErrorTrackerLastError(t *testing.T) {
	tracker := NewErrorTracker()

	err1 := NewError(ErrorCodeInternal, "first")
	err2 := NewError(ErrorCodeTimeout, "second")

	tracker.Track(err1)
	tracker.Track(err2)

	stats := tracker.GetStats()
	if stats.LastError == nil {
		t.Fatal("LastError should not be nil")
	}
	if stats.LastError.Code != ErrorCodeTimeout {
		t.Error("LastError should be the most recent error")
	}
}

// BenchmarkNewError 基准测试：创建错误
func BenchmarkNewError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewError(ErrorCodeInternal, "test error")
	}
}

// BenchmarkWrapError 基准测试：包装错误
func BenchmarkWrapError(b *testing.B) {
	err := errors.New("test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WrapError(err, ErrorCodeInternal, "wrapper")
	}
}

// BenchmarkErrorTracker 基准测试：错误跟踪
func BenchmarkErrorTracker(b *testing.B) {
	tracker := NewErrorTracker()
	err := NewError(ErrorCodeInternal, "test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracker.Track(err)
	}
}

// TestAgentErrorChaining 测试错误链
func TestAgentErrorChaining(t *testing.T) {
	err1 := errors.New("base error")
	err2 := WrapError(err1, ErrorCodeInternal, "layer 2")
	err3 := WrapError(err2, ErrorCodeProcessingFailed, "layer 3")

	// 验证可以unwrap到原始错误
	if !errors.Is(err3, err2) {
		t.Error("Should be able to find err2 in chain")
	}

	// 验证错误消息包含层级信息
	if err3.Error() == "" {
		t.Error("Error message should not be empty")
	}
}
