package user

import (
	"context"
	"errors"

	"github.com/yourusername/golang/internal/domain/user"
)

// Service 是用户应用服务，负责协调领域对象完成业务用例。
//
// 设计原理：
// 1. 应用服务负责用例编排，不包含业务逻辑
// 2. 业务逻辑应该在领域层（实体、领域服务）
// 3. 应用服务协调领域对象完成业务用例
// 4. 应用服务负责 DTO 转换
//
// 架构位置：
// - Application Layer (internal/application/user/)
// - 依赖 Domain Layer (internal/domain/user/)
//
// 职责：
// 1. 用例编排：协调领域对象完成业务用例
// 2. 事务管理：应用服务方法通常是一个事务边界
// 3. DTO 转换：领域对象和 DTO 之间的转换
// 4. 错误处理：将领域错误转换为应用错误
//
// 使用场景：
// 1. HTTP Handler 调用应用服务
// 2. gRPC Handler 调用应用服务
// 3. GraphQL Resolver 调用应用服务
// 4. 工作流 Activity 调用应用服务
type Service struct {
	userRepo user.Repository
}

// NewService 创建新的用户应用服务。
//
// 功能说明：
// - 接收用户仓储接口
// - 创建并返回配置好的服务实例
//
// 参数：
// - userRepo: 用户仓储接口
//
// 返回：
// - *Service: 配置好的服务实例
//
// 使用示例：
//
//	userRepo := infrastructure.NewUserRepository(client)
//	userService := user.NewService(userRepo)
func NewService(userRepo user.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

// CreateUserRequest 创建用户请求。
//
// 字段说明：
// - Email: 用户邮箱（必填）
// - Name: 用户名称（必填）
type CreateUserRequest struct {
	Email string
	Name  string
}

// UserResponse 用户响应。
//
// 字段说明：
// - ID: 用户唯一标识
// - Email: 用户邮箱
// - Name: 用户名称
// - CreatedAt: 创建时间
// - UpdatedAt: 更新时间
type UserResponse struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
	UpdatedAt string
}

// CreateUser 创建用户。
//
// 功能说明：
// 1. 验证请求参数
// 2. 检查邮箱是否已存在
// 3. 创建用户实体
// 4. 验证用户实体
// 5. 保存用户
// 6. 转换为响应 DTO
//
// 参数：
// - ctx: 上下文
// - req: 创建用户请求
//
// 返回：
// - *UserResponse: 用户响应
// - error: 创建失败时返回错误
//
// 错误处理：
// - ErrUserAlreadyExists: 用户已存在
// - ErrInvalidInput: 无效的输入
// - ErrInternal: 内部错误
//
// 使用示例：
//
//	req := user.CreateUserRequest{
//	    Email: "test@example.com",
//	    Name:  "Test User",
//	}
//	resp, err := service.CreateUser(ctx, req)
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	// 1. 验证请求参数
	if req.Email == "" {
		return nil, ErrInvalidInput
	}
	if req.Name == "" {
		return nil, ErrInvalidInput
	}

	// 2. 检查邮箱是否已存在
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return nil, ErrInternal
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 3. 创建用户实体
	domainUser := user.NewUser(req.Email, req.Name)

	// 4. 验证用户实体
	if err := domainUser.IsValid(); err != nil {
		return nil, ErrInvalidInput
	}

	// 5. 保存用户
	if err := s.userRepo.Create(ctx, domainUser); err != nil {
		return nil, ErrInternal
	}

	// 6. 转换为响应 DTO
	return toUserResponse(domainUser), nil
}

// GetUser 获取用户。
//
// 功能说明：
// 1. 根据 ID 查找用户
// 2. 转换为响应 DTO
//
// 参数：
// - ctx: 上下文
// - id: 用户 ID
//
// 返回：
// - *UserResponse: 用户响应
// - error: 获取失败时返回错误
//
// 错误处理：
// - ErrUserNotFound: 用户不存在
// - ErrInternal: 内部错误
func (s *Service) GetUser(ctx context.Context, id string) (*UserResponse, error) {
	domainUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}

	return toUserResponse(domainUser), nil
}

// UpdateUserRequest 更新用户请求。
//
// 字段说明：
// - Name: 用户名称（可选）
// - Email: 用户邮箱（可选）
type UpdateUserRequest struct {
	Name  *string
	Email *string
}

// ListUsersRequest 列出用户请求。
//
// 字段说明：
// - Limit: 每页数量（默认 10，最大 100）
// - Offset: 偏移量（默认 0）
type ListUsersRequest struct {
	Limit  int
	Offset int
}

// ListUsersResponse 列出用户响应。
//
// 字段说明：
// - Users: 用户列表
// - Total: 总数量
// - Limit: 每页数量
// - Offset: 偏移量
type ListUsersResponse struct {
	Users  []*UserResponse
	Total  int
	Limit  int
	Offset int
}

// ListUsers 列出用户。
//
// 功能说明：
// 1. 验证分页参数
// 2. 调用仓储获取用户列表
// 3. 转换为响应 DTO
//
// 参数：
// - ctx: 上下文
// - req: 列出用户请求
//
// 返回：
// - *ListUsersResponse: 用户列表响应
// - error: 获取失败时返回错误
//
// 错误处理：
// - ErrInternal: 内部错误
func (s *Service) ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error) {
	// 1. 验证分页参数
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// 2. 调用仓储获取用户列表
	domainUsers, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, ErrInternal
	}

	// 3. 转换为响应 DTO
	users := make([]*UserResponse, 0, len(domainUsers))
	for _, u := range domainUsers {
		users = append(users, toUserResponse(u))
	}

	return &ListUsersResponse{
		Users:  users,
		Total:  len(users), // TODO: 需要仓储提供 Count 方法获取总数
		Limit:  limit,
		Offset: offset,
	}, nil
}

// UpdateUser 更新用户。
//
// 功能说明：
// 1. 根据 ID 查找用户
// 2. 更新用户信息
// 3. 验证用户实体
// 4. 保存更新
// 5. 转换为响应 DTO
//
// 参数：
// - ctx: 上下文
// - id: 用户 ID
// - req: 更新用户请求
//
// 返回：
// - *UserResponse: 用户响应
// - error: 更新失败时返回错误
//
// 错误处理：
// - ErrUserNotFound: 用户不存在
// - ErrInvalidInput: 无效的输入
// - ErrInternal: 内部错误
func (s *Service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*UserResponse, error) {
	// 1. 根据 ID 查找用户
	domainUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}

	// 2. 更新用户信息
	if req.Name != nil {
		domainUser.UpdateName(*req.Name)
	}
	if req.Email != nil {
		// 检查邮箱是否已被其他用户使用
		existingUser, err := s.userRepo.FindByEmail(ctx, *req.Email)
		if err != nil && !errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrInternal
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrUserAlreadyExists
		}
		domainUser.UpdateEmail(*req.Email)
	}

	// 3. 验证用户实体
	if err := domainUser.IsValid(); err != nil {
		return nil, ErrInvalidInput
	}

	// 4. 保存更新
	if err := s.userRepo.Update(ctx, domainUser); err != nil {
		return nil, ErrInternal
	}

	// 5. 转换为响应 DTO
	return toUserResponse(domainUser), nil
}

// DeleteUser 删除用户。
//
// 功能说明：
// 1. 根据 ID 查找用户（验证存在）
// 2. 删除用户
//
// 参数：
// - ctx: 上下文
// - id: 用户 ID
//
// 返回：
// - error: 删除失败时返回错误
//
// 错误处理：
// - ErrUserNotFound: 用户不存在
// - ErrInternal: 内部错误
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	// 1. 验证用户存在
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return ErrInternal
	}

	// 2. 删除用户
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return ErrInternal
	}

	return nil
}

// toUserResponse 将领域用户实体转换为响应 DTO。
//
// 功能说明：
// - 将领域对象转换为应用层 DTO
// - 格式化时间戳
func toUserResponse(u *user.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
