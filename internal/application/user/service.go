package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain "github.com/yourusername/golang/internal/domain/user"
)

// Service 应用服务接口
type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error)
	GetUser(ctx context.Context, id string) (*UserDTO, error)
	UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*UserDTO, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*UserDTO, error)
}

// service 应用服务实现
type service struct {
	repo domain.Repository
}

// NewService 创建应用服务
func NewService(repo domain.Repository) Service {
	return &service{repo: repo}
}

// CreateUser 创建用户
func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
	// 检查邮箱是否已存在
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// 创建领域实体
	user := domain.NewUser(req.Email, req.Name)
	user.ID = uuid.New().String()

	// 保存到仓储
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 转换为 DTO
	return ToDTO(user), nil
}

// GetUser 获取用户
func (s *service) GetUser(ctx context.Context, id string) (*UserDTO, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return ToDTO(user), nil
}

// UpdateUser 更新用户
func (s *service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*UserDTO, error) {
	// 获取现有用户
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		user.UpdateName(*req.Name)
	}
	if req.Email != nil {
		// 检查新邮箱是否已被使用
		if *req.Email != user.Email {
			existing, err := s.repo.FindByEmail(ctx, *req.Email)
			if err != nil && err != domain.ErrUserNotFound {
				return nil, fmt.Errorf("failed to check email: %w", err)
			}
			if existing != nil {
				return nil, fmt.Errorf("email already in use: %w", domain.ErrUserAlreadyExists)
			}
		}
		user.UpdateEmail(*req.Email)
	}

	// 保存更新
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return ToDTO(user), nil
}

// DeleteUser 删除用户
func (s *service) DeleteUser(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers 列出用户
func (s *service) ListUsers(ctx context.Context, limit, offset int) ([]*UserDTO, error) {
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
