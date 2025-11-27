package user

import "errors"

// 用户领域错误定义
//
// 设计原则：
// 1. 错误应该表达业务含义
// 2. 错误应该提供足够的上下文信息
// 3. 错误应该可以被外部层处理

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidUserID 无效的用户 ID
	ErrInvalidUserID = errors.New("invalid user id")

	// ErrEmailRequired 邮箱必填
	ErrEmailRequired = errors.New("email is required")

	// ErrInvalidEmailFormat 无效的邮箱格式
	ErrInvalidEmailFormat = errors.New("invalid email format")

	// ErrEmailAlreadyExists 邮箱已存在
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrNameRequired 名称必填
	ErrNameRequired = errors.New("name is required")

	// ErrNameTooShort 名称太短
	ErrNameTooShort = errors.New("name is too short (minimum 2 characters)")

	// ErrNameTooLong 名称太长
	ErrNameTooLong = errors.New("name is too long (maximum 100 characters)")
)
