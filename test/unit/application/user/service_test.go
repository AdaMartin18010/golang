package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	appuser "github.com/yourusername/golang/internal/application/user"
	domain "github.com/yourusername/golang/internal/domain/user"
)

// MockRepository 模拟仓储
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

// TestService_CreateUser 测试创建用户（成功场景）
func TestService_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		req     appuser.CreateUserRequest
		setup   func(*MockRepository)
		wantErr error
	}{
		{
			name: "success",
			req: appuser.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			setup: func(m *MockRepository) {
				m.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, domain.ErrUserNotFound)
				m.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.Email == "test@example.com" && u.Name == "Test User"
				})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "email already exists",
			req: appuser.CreateUserRequest{
				Email: "existing@example.com",
				Name:  "Test User",
			},
			setup: func(m *MockRepository) {
				existingUser := domain.NewUser("existing@example.com", "Existing User")
				m.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			wantErr: appuser.ErrUserAlreadyExists,
		},
		{
			name: "empty email",
			req: appuser.CreateUserRequest{
				Email: "",
				Name:  "Test User",
			},
			setup:   func(m *MockRepository) {},
			wantErr: appuser.ErrInvalidInput,
		},
		{
			name: "empty name",
			req: appuser.CreateUserRequest{
				Email: "test@example.com",
				Name:  "",
			},
			setup:   func(m *MockRepository) {},
			wantErr: appuser.ErrInvalidInput,
		},
		{
			name: "invalid email format",
			req: appuser.CreateUserRequest{
				Email: "invalid-email",
				Name:  "Test User",
			},
			setup: func(m *MockRepository) {
				m.On("FindByEmail", mock.Anything, "invalid-email").Return(nil, domain.ErrUserNotFound)
			},
			wantErr: appuser.ErrInvalidInput,
		},
		{
			name: "repository error on FindByEmail",
			req: appuser.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			setup: func(m *MockRepository) {
				m.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("database error"))
			},
			wantErr: appuser.ErrInternal,
		},
		{
			name: "repository error on Create",
			req: appuser.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			setup: func(m *MockRepository) {
				m.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, domain.ErrUserNotFound)
				m.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			wantErr: appuser.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockRepository)
			tt.setup(mockRepo)
			service := appuser.NewService(mockRepo)

			resp, err := service.CreateUser(ctx, tt.req)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Email, resp.Email)
				assert.Equal(t, tt.req.Name, resp.Name)
				assert.NotEmpty(t, resp.ID)
				assert.NotEmpty(t, resp.CreatedAt)
				assert.NotEmpty(t, resp.UpdatedAt)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected error %v, got %v", tt.wantErr, err)
				assert.Nil(t, resp)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestService_GetUser 测试获取用户
func TestService_GetUser(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		setup   func(*MockRepository)
		wantErr error
	}{
		{
			name:   "success",
			userID: "user-123",
			setup: func(m *MockRepository) {
				user := domain.NewUser("test@example.com", "Test User")
				user.ID = "user-123"
				m.On("FindByID", mock.Anything, "user-123").Return(user, nil)
			},
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: "non-existent",
			setup: func(m *MockRepository) {
				m.On("FindByID", mock.Anything, "non-existent").Return(nil, domain.ErrUserNotFound)
			},
			wantErr: appuser.ErrUserNotFound,
		},
		{
			name:   "repository error",
			userID: "user-123",
			setup: func(m *MockRepository) {
				m.On("FindByID", mock.Anything, "user-123").Return(nil, errors.New("database error"))
			},
			wantErr: appuser.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockRepository)
			tt.setup(mockRepo)
			service := appuser.NewService(mockRepo)

			resp, err := service.GetUser(ctx, tt.userID)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.userID, resp.ID)
				assert.NotEmpty(t, resp.Email)
				assert.NotEmpty(t, resp.Name)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected error %v, got %v", tt.wantErr, err)
				assert.Nil(t, resp)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestService_ListUsers 测试列出用户
func TestService_ListUsers(t *testing.T) {
	tests := []struct {
		name    string
		req     appuser.ListUsersRequest
		setup   func(*MockRepository)
		wantErr error
	}{
		{
			name: "success",
			req: appuser.ListUsersRequest{
				Limit:  10,
				Offset: 0,
			},
			setup: func(m *MockRepository) {
				users := []*domain.User{
					domain.NewUser("user1@example.com", "User 1"),
					domain.NewUser("user2@example.com", "User 2"),
				}
				m.On("List", mock.Anything, 10, 0).Return(users, nil)
			},
			wantErr: nil,
		},
		{
			name: "empty list",
			req: appuser.ListUsersRequest{
				Limit:  10,
				Offset: 0,
			},
			setup: func(m *MockRepository) {
				m.On("List", mock.Anything, 10, 0).Return([]*domain.User{}, nil)
			},
			wantErr: nil,
		},
		{
			name: "default limit",
			req: appuser.ListUsersRequest{
				Limit:  0,
				Offset: 0,
			},
			setup: func(m *MockRepository) {
				m.On("List", mock.Anything, 10, 0).Return([]*domain.User{}, nil)
			},
			wantErr: nil,
		},
		{
			name: "limit exceeds max",
			req: appuser.ListUsersRequest{
				Limit:  200,
				Offset: 0,
			},
			setup: func(m *MockRepository) {
				m.On("List", mock.Anything, 100, 0).Return([]*domain.User{}, nil)
			},
			wantErr: nil,
		},
		{
			name: "repository error",
			req: appuser.ListUsersRequest{
				Limit:  10,
				Offset: 0,
			},
			setup: func(m *MockRepository) {
				m.On("List", mock.Anything, 10, 0).Return(nil, errors.New("database error"))
			},
			wantErr: appuser.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockRepository)
			tt.setup(mockRepo)
			service := appuser.NewService(mockRepo)

			resp, err := service.ListUsers(ctx, tt.req)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Users)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
				assert.Nil(t, resp)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestService_UpdateUser 测试更新用户
func TestService_UpdateUser(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		req     appuser.UpdateUserRequest
		setup   func(*MockRepository)
		wantErr error
	}{
		{
			name:   "success - update name",
			userID: "user-123",
			req: appuser.UpdateUserRequest{
				Name: stringPtr("New Name"),
			},
			setup: func(m *MockRepository) {
				user := domain.NewUser("test@example.com", "Old Name")
				user.ID = "user-123"
				m.On("FindByID", mock.Anything, "user-123").Return(user, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.ID == "user-123" && u.Name == "New Name"
				})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "success - update email",
			userID: "user-123",
			req: appuser.UpdateUserRequest{
				Email: stringPtr("new@example.com"),
			},
			setup: func(m *MockRepository) {
				user := domain.NewUser("old@example.com", "Test User")
				user.ID = "user-123"
				m.On("FindByID", mock.Anything, "user-123").Return(user, nil)
				m.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, domain.ErrUserNotFound)
				m.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.ID == "user-123" && u.Email == "new@example.com"
				})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: "non-existent",
			req: appuser.UpdateUserRequest{
				Name: stringPtr("New Name"),
			},
			setup: func(m *MockRepository) {
				m.On("FindByID", mock.Anything, "non-existent").Return(nil, domain.ErrUserNotFound)
			},
			wantErr: appuser.ErrUserNotFound,
		},
		{
			name:   "email already exists",
			userID: "user-123",
			req: appuser.UpdateUserRequest{
				Email: stringPtr("existing@example.com"),
			},
			setup: func(m *MockRepository) {
				user := domain.NewUser("old@example.com", "Test User")
				user.ID = "user-123"
				existingUser := domain.NewUser("existing@example.com", "Existing User")
				existingUser.ID = "user-456"
				m.On("FindByID", mock.Anything, "user-123").Return(user, nil)
				m.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			wantErr: appuser.ErrUserAlreadyExists,
		},
		{
			name:   "invalid input",
			userID: "user-123",
			req: appuser.UpdateUserRequest{
				Name: stringPtr("A"), // 名称太短
			},
			setup: func(m *MockRepository) {
				user := domain.NewUser("test@example.com", "Test User")
				user.ID = "user-123"
				m.On("FindByID", mock.Anything, "user-123").Return(user, nil)
			},
			wantErr: appuser.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockRepository)
			tt.setup(mockRepo)
			service := appuser.NewService(mockRepo)

			resp, err := service.UpdateUser(ctx, tt.userID, tt.req)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
				assert.Nil(t, resp)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestService_DeleteUser 测试删除用户
func TestService_DeleteUser(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		setup   func(*MockRepository)
		wantErr error
	}{
		{
			name:   "success",
			userID: "user-123",
			setup: func(m *MockRepository) {
				user := domain.NewUser("test@example.com", "Test User")
				user.ID = "user-123"
				m.On("FindByID", mock.Anything, "user-123").Return(user, nil)
				m.On("Delete", mock.Anything, "user-123").Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: "non-existent",
			setup: func(m *MockRepository) {
				m.On("FindByID", mock.Anything, "non-existent").Return(nil, domain.ErrUserNotFound)
			},
			wantErr: appuser.ErrUserNotFound,
		},
		{
			name:   "repository error on delete",
			userID: "user-123",
			setup: func(m *MockRepository) {
				user := domain.NewUser("test@example.com", "Test User")
				user.ID = "user-123"
				m.On("FindByID", mock.Anything, "user-123").Return(user, nil)
				m.On("Delete", mock.Anything, "user-123").Return(errors.New("database error"))
			},
			wantErr: appuser.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockRepository)
			tt.setup(mockRepo)
			service := appuser.NewService(mockRepo)

			err := service.DeleteUser(ctx, tt.userID)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestService_NewService 测试创建服务
func TestService_NewService(t *testing.T) {
	mockRepo := new(MockRepository)
	service := appuser.NewService(mockRepo)

	assert.NotNil(t, service)
}

// stringPtr 返回字符串指针（辅助函数）
func stringPtr(s string) *string {
	return &s
}
