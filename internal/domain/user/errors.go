package user

import "errors"

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidEmail 无效邮箱
	ErrInvalidEmail = errors.New("invalid email")
)
