package errors

import (
	"fmt"
	"time"
)

// ErrorCode 错误代码
type ErrorCode string

const (
	// ErrCodeNotFound 资源不存在
	ErrCodeNotFound ErrorCode = "NOT_FOUND"
	// ErrCodeInvalidInput 无效输入
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	// ErrCodeInternal 内部错误
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
	// ErrCodeUnauthorized 未授权
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	// ErrCodeForbidden 禁止访问
	ErrCodeForbidden ErrorCode = "FORBIDDEN"
	// ErrCodeConflict 资源冲突
	ErrCodeConflict ErrorCode = "CONFLICT"
	// ErrCodeValidation 验证失败
	ErrCodeValidation ErrorCode = "VALIDATION_ERROR"
	// ErrCodeTimeout 超时
	ErrCodeTimeout ErrorCode = "TIMEOUT"
	// ErrCodeRateLimit 限流
	ErrCodeRateLimit ErrorCode = "RATE_LIMIT_EXCEEDED"
	// ErrCodeServiceUnavailable 服务不可用
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// ErrorCategory 错误分类
type ErrorCategory string

const (
	CategoryClient   ErrorCategory = "CLIENT_ERROR"   // 客户端错误 (4xx)
	CategoryServer   ErrorCategory = "SERVER_ERROR"   // 服务端错误 (5xx)
	CategoryBusiness ErrorCategory = "BUSINESS_ERROR" // 业务错误
)

// HTTPStatus 返回对应的HTTP状态码
func (c ErrorCategory) HTTPStatus() int {
	switch c {
	case CategoryClient:
		return 400
	case CategoryServer:
		return 500
	case CategoryBusiness:
		return 400
	default:
		return 500
	}
}

// AppError 应用错误
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Cause      error                  `json:"-"`
	Category   ErrorCategory          `json:"category"`
	HTTPStatus int                    `json:"http_status"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Retryable  bool                   `json:"retryable"`
	TraceID    string                 `json:"trace_id,omitempty"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 返回底层错误
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithDetails 添加详细信息
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithTraceID 添加追踪ID
func (e *AppError) WithTraceID(traceID string) *AppError {
	e.TraceID = traceID
	return e
}

// IsRetryable 判断是否可重试
func (e *AppError) IsRetryable() bool {
	return e.Retryable
}

// HTTPStatusCode 返回HTTP状态码
func (e *AppError) HTTPStatusCode() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}
	return e.Category.HTTPStatus()
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(resource string, id string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s with id %s not found", resource, id),
		Category:   CategoryClient,
		HTTPStatus: 404,
		Timestamp:  time.Now(),
		Retryable:  false,
	}
}

// NewInvalidInputError 创建无效输入错误
func NewInvalidInputError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidInput,
		Message:    message,
		Category:   CategoryClient,
		HTTPStatus: 400,
		Timestamp:  time.Now(),
		Retryable:  false,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		Category:   CategoryClient,
		HTTPStatus: 400,
		Details:    details,
		Timestamp:  time.Now(),
		Retryable:  false,
	}
}

// NewInternalError 创建内部错误
func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    message,
		Cause:      cause,
		Category:   CategoryServer,
		HTTPStatus: 500,
		Timestamp:  time.Now(),
		Retryable:  false,
	}
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		Category:   CategoryClient,
		HTTPStatus: 401,
		Timestamp:  time.Now(),
		Retryable:  false,
	}
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		Category:   CategoryClient,
		HTTPStatus: 403,
		Timestamp:  time.Now(),
		Retryable:  false,
	}
}

// NewConflictError 创建资源冲突错误
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		Category:   CategoryClient,
		HTTPStatus: 409,
		Timestamp:  time.Now(),
		Retryable:  false,
	}
}

// NewTimeoutError 创建超时错误
func NewTimeoutError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeTimeout,
		Message:    message,
		Category:   CategoryServer,
		HTTPStatus: 504,
		Timestamp:  time.Now(),
		Retryable:  true,
	}
}

// NewRateLimitError 创建限流错误
func NewRateLimitError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeRateLimit,
		Message:    message,
		Category:   CategoryClient,
		HTTPStatus: 429,
		Timestamp:  time.Now(),
		Retryable:  true,
	}
}

// NewServiceUnavailableError 创建服务不可用错误
func NewServiceUnavailableError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeServiceUnavailable,
		Message:    message,
		Category:   CategoryServer,
		HTTPStatus: 503,
		Timestamp:  time.Now(),
		Retryable:  true,
	}
}

// FromDomainError 从领域错误转换为应用错误
// 这是一个通用转换函数，具体实现需要根据实际的领域错误类型来调整
func FromDomainError(err error) *AppError {
	if err == nil {
		return nil
	}

	// 如果已经是 AppError，直接返回
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	// 根据错误消息或类型进行转换
	// 这里可以根据实际的领域错误类型进行扩展
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    "Internal error",
		Cause:      err,
		Category:   CategoryServer,
		HTTPStatus: 500,
		Timestamp:  time.Now(),
		Retryable:  false,
	}
}
