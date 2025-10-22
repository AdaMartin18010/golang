package core

import (
	"errors"
	"fmt"
	"time"
)

// =============================================================================
// 增强的错误处理 - Enhanced Error Handling
// =============================================================================

// ErrorCode 错误代码
type ErrorCode string

const (
	// 系统错误
	ErrorCodeInternal          ErrorCode = "INTERNAL_ERROR"
	ErrorCodeTimeout           ErrorCode = "TIMEOUT"
	ErrorCodeCancelled         ErrorCode = "CANCELLED"
	ErrorCodeInvalidState      ErrorCode = "INVALID_STATE"
	ErrorCodeResourceExhausted ErrorCode = "RESOURCE_EXHAUSTED"

	// 配置错误
	ErrorCodeInvalidConfig ErrorCode = "INVALID_CONFIG"
	ErrorCodeMissingConfig ErrorCode = "MISSING_CONFIG"

	// 插件错误
	ErrorCodePluginNotFound ErrorCode = "PLUGIN_NOT_FOUND"
	ErrorCodePluginFailed   ErrorCode = "PLUGIN_FAILED"

	// 处理错误
	ErrorCodeProcessingFailed ErrorCode = "PROCESSING_FAILED"
	ErrorCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrorCodeInvalidOutput    ErrorCode = "INVALID_OUTPUT"

	// 决策错误
	ErrorCodeDecisionFailed ErrorCode = "DECISION_FAILED"
	ErrorCodeNoDecision     ErrorCode = "NO_DECISION"

	// 学习错误
	ErrorCodeLearningFailed    ErrorCode = "LEARNING_FAILED"
	ErrorCodeInvalidExperience ErrorCode = "INVALID_EXPERIENCE"
)

// AgentError Agent错误类型
type AgentError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	Cause     error                  `json:"cause,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Retryable bool                   `json:"retryable"`
}

// Error 实现error接口
func (e *AgentError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (caused by: %v)", e.Code, e.Message, e.Details, e.Cause)
	}
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 实现errors.Unwrap接口
func (e *AgentError) Unwrap() error {
	return e.Cause
}

// Is 实现errors.Is接口
func (e *AgentError) Is(target error) bool {
	t, ok := target.(*AgentError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// WithContext 添加上下文信息
func (e *AgentError) WithContext(key string, value interface{}) *AgentError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithDetails 添加详细信息
func (e *AgentError) WithDetails(details string) *AgentError {
	e.Details = details
	return e
}

// NewError 创建新错误
func NewError(code ErrorCode, message string) *AgentError {
	return &AgentError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Retryable: false,
	}
}

// NewRetryableError 创建可重试的错误
func NewRetryableError(code ErrorCode, message string) *AgentError {
	return &AgentError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Retryable: true,
	}
}

// WrapError 包装错误
func WrapError(err error, code ErrorCode, message string) *AgentError {
	if err == nil {
		return nil
	}

	// 如果已经是AgentError，添加新的层级
	if agentErr, ok := err.(*AgentError); ok {
		return &AgentError{
			Code:      code,
			Message:   message,
			Cause:     agentErr,
			Timestamp: time.Now(),
			Retryable: agentErr.Retryable,
		}
	}

	return &AgentError{
		Code:      code,
		Message:   message,
		Cause:     err,
		Timestamp: time.Now(),
		Retryable: false,
	}
}

// =============================================================================
// 预定义错误
// =============================================================================

// ErrAgentNotRunning Agent未运行错误
var ErrAgentNotRunning = NewError(ErrorCodeInvalidState, "agent is not running")

// ErrAgentAlreadyRunning Agent已运行错误
var ErrAgentAlreadyRunning = NewError(ErrorCodeInvalidState, "agent is already running")

// ErrInvalidInput 无效输入错误
var ErrInvalidInput = NewError(ErrorCodeInvalidInput, "invalid input data")

// ErrTimeout 超时错误
var ErrTimeout = NewRetryableError(ErrorCodeTimeout, "operation timed out")

// ErrPluginNotFound 插件未找到错误
var ErrPluginNotFound = NewError(ErrorCodePluginNotFound, "plugin not found")

// =============================================================================
// 错误处理工具函数
// =============================================================================

// IsRetryable 判断错误是否可重试
func IsRetryable(err error) bool {
	var agentErr *AgentError
	if errors.As(err, &agentErr) {
		return agentErr.Retryable
	}
	return false
}

// GetErrorCode 获取错误代码
func GetErrorCode(err error) ErrorCode {
	var agentErr *AgentError
	if errors.As(err, &agentErr) {
		return agentErr.Code
	}
	return ErrorCodeInternal
}

// GetErrorContext 获取错误上下文
func GetErrorContext(err error) map[string]interface{} {
	var agentErr *AgentError
	if errors.As(err, &agentErr) {
		return agentErr.Context
	}
	return nil
}

// =============================================================================
// 错误统计
// =============================================================================

// ErrorStats 错误统计
type ErrorStats struct {
	TotalErrors    int64
	ErrorsByCode   map[ErrorCode]int64
	RetryableCount int64
	LastError      *AgentError
	LastErrorTime  time.Time
}

// ErrorTracker 错误跟踪器
type ErrorTracker struct {
	stats ErrorStats
}

// NewErrorTracker 创建错误跟踪器
func NewErrorTracker() *ErrorTracker {
	return &ErrorTracker{
		stats: ErrorStats{
			ErrorsByCode: make(map[ErrorCode]int64),
		},
	}
}

// Track 跟踪错误
func (et *ErrorTracker) Track(err error) {
	if err == nil {
		return
	}

	et.stats.TotalErrors++
	et.stats.LastErrorTime = time.Now()

	var agentErr *AgentError
	if errors.As(err, &agentErr) {
		et.stats.ErrorsByCode[agentErr.Code]++
		et.stats.LastError = agentErr
		if agentErr.Retryable {
			et.stats.RetryableCount++
		}
	} else {
		et.stats.ErrorsByCode[ErrorCodeInternal]++
	}
}

// GetStats 获取统计信息
func (et *ErrorTracker) GetStats() ErrorStats {
	return et.stats
}

// Reset 重置统计
func (et *ErrorTracker) Reset() {
	et.stats = ErrorStats{
		ErrorsByCode: make(map[ErrorCode]int64),
	}
}
