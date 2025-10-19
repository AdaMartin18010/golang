package domain

import "fmt"

// BusinessError 定义业务错误类型
type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error 实现error接口
func (e *BusinessError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewBusinessError 创建新的业务错误
func NewBusinessError(code, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// 预定义的错误类型
var (
	ErrUserNotFound      = NewBusinessError("USER_NOT_FOUND", "user not found")
	ErrUserAlreadyExists = NewBusinessError("USER_ALREADY_EXISTS", "user already exists")
	ErrInvalidInput      = NewBusinessError("INVALID_INPUT", "invalid input data")
	ErrDatabaseError     = NewBusinessError("DATABASE_ERROR", "database operation failed")
)
