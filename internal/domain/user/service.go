package user

import "context"

// DomainService 领域服务接口
type DomainService interface {
	// ValidateEmail 验证邮箱格式
	ValidateEmail(email string) bool

	// IsEmailUnique 检查邮箱是否唯一
	IsEmailUnique(ctx context.Context, email string) (bool, error)
}
