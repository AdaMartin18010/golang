// Package errors provides a unified error handling framework for the application.
//
// 统一错误处理框架提供了标准化的错误处理机制，包括：
// 1. 错误代码体系：统一的错误代码定义
// 2. 错误分类：客户端错误、服务端错误、业务错误
// 3. HTTP 状态码映射：自动映射到 HTTP 状态码
// 4. 详细信息支持：支持添加详细的错误信息
// 5. 追踪支持：支持添加追踪 ID
// 6. 可重试标记：标记错误是否可重试
//
// 设计原则：
// 1. 统一格式：所有错误使用相同的结构
// 2. 类型安全：使用类型化的错误代码
// 3. 可扩展性：支持添加详细信息和上下文
// 4. 可追踪性：支持错误链和追踪 ID
//
// 使用场景：
// - 应用层错误处理
// - HTTP 错误响应
// - 日志记录和错误追踪
// - 错误重试决策
//
// 示例：
//
//	// 创建错误
//	err := errors.NewNotFoundError("user", "123")
//
//	// 添加详细信息
//	err = err.WithDetails("field", "email").WithTraceID("trace-123")
//
//	// 检查错误类型
//	if appErr, ok := err.(*errors.AppError); ok {
//	    log.Printf("Error code: %s, HTTP status: %d", appErr.Code, appErr.HTTPStatusCode())
//	}
package errors

import (
	"fmt"
	"time"
)

// ErrorCode 是错误代码的类型定义。
//
// 功能说明：
// - 使用字符串类型表示错误代码
// - 提供类型安全的错误代码常量
// - 便于错误分类和处理
//
// 错误代码约定：
// - 使用大写字母和下划线（如 "NOT_FOUND"）
// - 错误代码应该是唯一的、有意义的
// - 错误代码应该与 HTTP 状态码对应
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

// ErrorCategory 是错误分类的类型定义。
//
// 功能说明：
// - 将错误分为客户端错误、服务端错误和业务错误
// - 用于错误处理和 HTTP 状态码映射
//
// 错误分类：
// - CategoryClient: 客户端错误（4xx），由客户端请求引起
// - CategoryServer: 服务端错误（5xx），由服务器内部问题引起
// - CategoryBusiness: 业务错误，业务逻辑相关的错误
type ErrorCategory string

const (
	CategoryClient   ErrorCategory = "CLIENT_ERROR"   // 客户端错误 (4xx)
	CategoryServer   ErrorCategory = "SERVER_ERROR"   // 服务端错误 (5xx)
	CategoryBusiness ErrorCategory = "BUSINESS_ERROR" // 业务错误
)

// HTTPStatus 返回错误分类对应的默认 HTTP 状态码。
//
// 功能说明：
// - 根据错误分类返回默认的 HTTP 状态码
// - 客户端错误返回 400
// - 服务端错误返回 500
// - 业务错误返回 400
//
// 返回：
// - int: HTTP 状态码
//
// 注意事项：
// - 这是默认状态码，可以在 AppError 中覆盖
// - 某些错误类型有特定的状态码（如 404、401、403 等）
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

// AppError 是应用程序错误的统一结构。
//
// 功能说明：
// - 提供统一的错误格式
// - 包含错误代码、消息、分类等信息
// - 支持错误链和详细信息
//
// 字段说明：
// - Code: 错误代码（如 "NOT_FOUND"、"INVALID_INPUT"）
// - Message: 错误消息（人类可读的错误描述）
// - Cause: 底层错误（错误链，不序列化到 JSON）
// - Category: 错误分类（客户端、服务端、业务）
// - HTTPStatus: HTTP 状态码（如 404、400、500）
// - Details: 详细信息（键值对，可选）
// - Timestamp: 错误发生时间
// - Retryable: 是否可重试
// - TraceID: 追踪 ID（用于分布式追踪，可选）
//
// 使用示例：
//
//	err := &AppError{
//	    Code:       ErrCodeNotFound,
//	    Message:    "User not found",
//	    Category:   CategoryClient,
//	    HTTPStatus: 404,
//	    Timestamp:  time.Now(),
//	    Retryable:  false,
//	}
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Cause      error                  `json:"-"` // 不序列化到 JSON
	Category   ErrorCategory          `json:"category"`
	HTTPStatus int                    `json:"http_status"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Retryable  bool                   `json:"retryable"`
	TraceID    string                 `json:"trace_id,omitempty"`
}

// Error 实现 error 接口，返回错误的字符串表示。
//
// 功能说明：
// - 返回格式化的错误字符串
// - 如果存在底层错误，包含错误链信息
//
// 返回：
// - string: 错误字符串表示
//
// 格式：
// - 无底层错误："{Code}: {Message}"
// - 有底层错误："{Code}: {Message}: {Cause}"
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 返回底层错误，支持错误链。
//
// 功能说明：
// - 实现 errors.Unwrap 接口
// - 用于错误链的遍历和处理
//
// 返回：
// - error: 底层错误，如果不存在则返回 nil
//
// 使用示例：
//
//	err := errors.NewInternalError("Database error", dbErr)
//	cause := errors.Unwrap(err) // 返回 dbErr
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithDetails 添加详细信息到错误中。
//
// 功能说明：
// - 支持链式调用
// - 可以添加多个详细信息
//
// 参数：
// - key: 信息键
// - value: 信息值（可以是任意类型）
//
// 返回：
// - *AppError: 返回自身，支持链式调用
//
// 使用示例：
//
//	err := errors.NewValidationError("Validation failed", nil).
//	    WithDetails("field", "email").
//	    WithDetails("reason", "invalid format")
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithTraceID 添加追踪 ID 到错误中。
//
// 功能说明：
// - 用于分布式追踪
// - 支持链式调用
//
// 参数：
// - traceID: 追踪 ID（通常来自 OpenTelemetry 或其他追踪系统）
//
// 返回：
// - *AppError: 返回自身，支持链式调用
//
// 使用示例：
//
//	err := errors.NewInternalError("Database error", dbErr).
//	    WithTraceID("trace-123-456")
func (e *AppError) WithTraceID(traceID string) *AppError {
	e.TraceID = traceID
	return e
}

// IsRetryable 判断错误是否可重试。
//
// 功能说明：
// - 用于错误重试决策
// - 某些错误（如超时、限流）可以重试
//
// 返回：
// - bool: 如果错误可重试返回 true，否则返回 false
//
// 使用示例：
//
//	if err.IsRetryable() {
//	    // 执行重试逻辑
//	}
func (e *AppError) IsRetryable() bool {
	return e.Retryable
}

// HTTPStatusCode 返回错误的 HTTP 状态码。
//
// 功能说明：
// - 优先返回显式设置的 HTTPStatus
// - 如果没有设置，则根据错误分类返回默认状态码
//
// 返回：
// - int: HTTP 状态码
//
// 使用示例：
//
//	statusCode := err.HTTPStatusCode()
//	w.WriteHeader(statusCode)
func (e *AppError) HTTPStatusCode() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}
	return e.Category.HTTPStatus()
}

// NewNotFoundError 创建资源不存在错误。
//
// 功能说明：
// - 用于资源不存在的情况（如用户、订单等）
// - HTTP 状态码：404
// - 不可重试
//
// 参数：
// - resource: 资源类型（如 "user"、"order"）
// - id: 资源 ID
//
// 返回：
// - *AppError: 配置好的错误实例
//
// 使用示例：
//
//	err := errors.NewNotFoundError("user", "123")
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

// NewInvalidInputError 创建无效输入错误。
//
// 功能说明：
// - 用于输入参数无效的情况
// - HTTP 状态码：400
// - 不可重试
//
// 参数：
// - message: 错误消息
//
// 返回：
// - *AppError: 配置好的错误实例
//
// 使用示例：
//
//	err := errors.NewInvalidInputError("Email format is invalid")
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

// NewValidationError 创建验证错误。
//
// 功能说明：
// - 用于数据验证失败的情况
// - HTTP 状态码：400
// - 不可重试
// - 支持添加详细的验证错误信息
//
// 参数：
// - message: 错误消息
// - details: 详细的验证错误信息（可选，可以为 nil）
//
// 返回：
// - *AppError: 配置好的错误实例
//
// 使用示例：
//
//	err := errors.NewValidationError("Validation failed", map[string]interface{}{
//	    "email": "invalid format",
//	    "name":  "required",
//	})
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

// NewInternalError 创建内部错误。
//
// 功能说明：
// - 用于服务器内部错误（如数据库错误、系统错误等）
// - HTTP 状态码：500
// - 不可重试
// - 支持错误链（包含底层错误）
//
// 参数：
// - message: 错误消息
// - cause: 底层错误（可选，可以为 nil）
//
// 返回：
// - *AppError: 配置好的错误实例
//
// 使用示例：
//
//	err := errors.NewInternalError("Database connection failed", dbErr)
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

// NewUnauthorizedError 创建未授权错误。
//
// 功能说明：
// - 用于未授权访问的情况（如未登录、token 无效等）
// - HTTP 状态码：401
// - 不可重试
//
// 参数：
// - message: 错误消息
//
// 返回：
// - *AppError: 配置好的错误实例
//
// 使用示例：
//
//	err := errors.NewUnauthorizedError("Invalid or expired token")
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

// NewForbiddenError 创建禁止访问错误。
//
// 功能说明：
// - 用于权限不足的情况（如无权限访问资源）
// - HTTP 状态码：403
// - 不可重试
//
// 参数：
// - message: 错误消息
//
// 返回：
// - *AppError: 配置好的错误实例
//
// 使用示例：
//
//	err := errors.NewForbiddenError("You don't have permission to access this resource")
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

// NewConflictError 创建资源冲突错误。
//
// 功能说明：
// - 用于资源冲突的情况（如重复创建、状态冲突等）
// - HTTP 状态码：409
// - 不可重试
//
// 参数：
// - message: 错误消息
//
// 返回：
// - *AppError: 配置好的错误实例
//
// 使用示例：
//
//	err := errors.NewConflictError("User with this email already exists")
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

// NewTimeoutError 创建超时错误。
//
// 功能说明：
// - 用于操作超时的情况
// - HTTP 状态码：504
// - 可重试
//
// 参数：
// - message: 错误消息
//
// 返回：
// - *AppError: 配置好的错误实例
//
// 使用示例：
//
//	err := errors.NewTimeoutError("Request timeout after 30 seconds")
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

// NewRateLimitError 创建限流错误。
//
// 功能说明：
// - 用于请求频率超过限制的情况
// - HTTP 状态码：429
// - 可重试（建议延迟后重试）
//
// 参数：
// - message: 错误消息
//
// 返回：
// - *AppError: 配置好的错误实例
//
// 使用示例：
//
//	err := errors.NewRateLimitError("Rate limit exceeded. Please try again later")
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

// NewServiceUnavailableError 创建服务不可用错误。
//
// 功能说明：
// - 用于服务暂时不可用的情况（如维护中、过载等）
// - HTTP 状态码：503
// - 可重试
//
// 参数：
// - message: 错误消息
//
// 返回：
// - *AppError: 配置好的错误实例
//
// 使用示例：
//
//	err := errors.NewServiceUnavailableError("Service is temporarily unavailable")
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

// FromDomainError 从领域错误转换为应用错误。
//
// 功能说明：
// - 将领域层的错误转换为应用层的 AppError
// - 如果已经是 AppError，直接返回
// - 否则包装为内部错误
//
// 参数：
// - err: 领域错误（可以是任意 error 类型）
//
// 返回：
// - *AppError: 应用错误，如果输入为 nil 则返回 nil
//
// 使用场景：
// - 在应用层处理领域层返回的错误
// - 将领域错误转换为适合返回给客户端的错误格式
//
// 使用示例：
//
//	domainErr := domainUser.Validate()
//	if domainErr != nil {
//	    return nil, errors.FromDomainError(domainErr)
//	}
//
// 注意事项：
// - 这是一个通用转换函数，具体实现需要根据实际的领域错误类型来调整
// - 可以根据领域错误的类型或消息进行更精确的转换
// - 建议在应用层定义更具体的转换函数
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
	// 例如：
	//   - 检查错误类型并映射到相应的错误代码
	//   - 解析错误消息并提取关键信息
	//   - 根据错误来源（数据库、外部服务等）设置不同的错误分类
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
