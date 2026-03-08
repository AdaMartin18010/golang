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

func (m *MockRepository) Save(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
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
		email   string
		userName string
		setup   func(*MockRepository)
		wantErr error
	}{
		{
			name:     "success",
			email:    "test@example.com",
			userName: "Test User",
			setup: func(m *MockRepository) {
				m.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, domain.ErrUserNotFound)
				m.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:     "email already exists",
			email:    "existing@example.com",
			userName: "Test User",
			setup: func(m *MockRepository) {
				existingUser := domain.NewUser("existing@example.com", "Existing User")
				m.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			wantErr: errors.New("user with email existing@example.com already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockRepository)
			tt.setup(mockRepo)
			service := appuser.NewService(mockRepo)

			user, err := service.CreateUser(ctx, tt.email, tt.userName)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.userName, user.Name)
			} else {
				assert.Error(t, err)
				assert.Nil(t, user)
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
			wantErr: domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockRepository)
			tt.setup(mockRepo)
			service := appuser.NewService(mockRepo)

			user, err := service.GetUser(ctx, tt.userID)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userID, user.ID)
			} else {
				assert.Error(t, err)
				assert.Nil(t, user)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestService_ListUsers 测试列出用户
func TestService_ListUsers(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		offset  int
		setup   func(*MockRepository)
		wantErr error
	}{
		{
			name:   "success",
			limit:  10,
			offset: 0,
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
			name:   "empty list",
			limit:  10,
			offset: 0,
			setup: func(m *MockRepository) {
				m.On("List", mock.Anything, 10, 0).Return([]*domain.User{}, nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockRepository)
			tt.setup(mockRepo)
			service := appuser.NewService(mockRepo)

			users, err := service.ListUsers(ctx, tt.limit, tt.offset)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, users)
			} else {
				assert.Error(t, err)
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
				m.On("Delete", mock.Anything, "user-123").Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: "non-existent",
			setup: func(m *MockRepository) {
				m.On("Delete", mock.Anything, "non-existent").Return(domain.ErrUserNotFound)
			},
			wantErr: domain.ErrUserNotFound,
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
