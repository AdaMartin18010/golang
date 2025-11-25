package user

import (
	"context"
	"fmt"

	domain "github.com/yourusername/golang/internal/domain/user"
)

// Service 用户应用服务
type Service struct {
	repo domain.Repository
}

// NewService 创建用户服务
func NewService(repo domain.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateUser 创建用户用例
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
	// 检查邮箱是否已存在
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("email already exists: %w", domain.ErrUserAlreadyExists)
	}

	// 创建领域实体
	user := domain.NewUser(req.Email, req.Name)

	// 保存到仓储
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 返回DTO
	return ToDTO(user), nil
}

// GetUser 获取用户用例
func (s *Service) GetUser(ctx context.Context, id string) (*UserDTO, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return ToDTO(user), nil
}

// UpdateUser 更新用户用例
func (s *Service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*UserDTO, error) {
	// 查找用户
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	// 更新字段
	if req.Name != nil {
		user.UpdateName(*req.Name)
	}
	if req.Email != nil {
		// 检查新邮箱是否已被使用
		existing, err := s.repo.FindByEmail(ctx, *req.Email)
		if err == nil && existing != nil && existing.ID != id {
			return nil, fmt.Errorf("email already exists: %w", domain.ErrUserAlreadyExists)
		}
		user.UpdateEmail(*req.Email)
	}

	// 保存更新
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return ToDTO(user), nil
}

// DeleteUser 删除用户用例
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers 列出用户用例
func (s *Service) ListUsers(ctx context.Context, limit, offset int) ([]*UserDTO, error) {
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	dtos := make([]*UserDTO, len(users))
	for i, u := range users {
		dtos[i] = ToDTO(u)
	}

	return dtos, nil
}
