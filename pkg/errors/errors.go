package errors

import (
	"fmt"
	domain "github.com/yourusername/golang/internal/domain/user"
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
)

// AppError 应用错误
type AppError struct {
	Code    ErrorCode
	Message string
	Cause   error
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

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(resource string, id string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s with id %s not found", resource, id),
	}
}

// NewInvalidInputError 创建无效输入错误
func NewInvalidInputError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeInvalidInput,
		Message: message,
	}
}

// NewInternalError 创建内部错误
func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Code:    ErrCodeInternal,
		Message: message,
		Cause:   cause,
	}
}

// NewConflictError 创建资源冲突错误
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeConflict,
		Message: message,
	}
}

// FromDomainError 从领域错误转换为应用错误
func FromDomainError(err error) *AppError {
	switch err {
	case domain.ErrUserNotFound:
		return &AppError{
			Code:    ErrCodeNotFound,
			Message: "User not found",
			Cause:   err,
		}
	case domain.ErrUserAlreadyExists:
		return &AppError{
			Code:    ErrCodeConflict,
			Message: "User already exists",
			Cause:   err,
		}
	default:
		return &AppError{
			Code:    ErrCodeInternal,
			Message: "Internal error",
			Cause:   err,
		}
	}
}
