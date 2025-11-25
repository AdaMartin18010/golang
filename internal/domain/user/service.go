package user

import (
	"context"
	"errors"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidEmail 无效邮箱
	ErrInvalidEmail = errors.New("invalid email")
)

// DomainService 领域服务接口
type DomainService interface {
	// ValidateEmail 验证邮箱格式
	ValidateEmail(email string) bool

	// IsEmailUnique 检查邮箱是否唯一
	IsEmailUnique(ctx context.Context, email string) (bool, error)
}

