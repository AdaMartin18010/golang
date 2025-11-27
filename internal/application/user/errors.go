package user

import "errors"

// 应用服务错误定义
//
// 设计原理：
// 1. 应用层错误应该表达应用层面的含义
// 2. 应用层错误可以包装领域错误
// 3. 应用层错误应该可以被接口层处理

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidInput 无效的输入
	ErrInvalidInput = errors.New("invalid input")

	// ErrInternal 内部错误
	ErrInternal = errors.New("internal error")
)
