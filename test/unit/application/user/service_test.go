package user

import (
	"context"
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

func TestService_CreateUser(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := appuser.NewService(mockRepo)

	req := appuser.CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Mock: 邮箱不存在
	mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, domain.ErrUserNotFound)
	// Mock: 创建成功
	mockRepo.On("Create", ctx, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == req.Email && u.Name == req.Name
	})).Return(nil)

	user, err := service.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Name, user.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_CreateUser_EmailExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := appuser.NewService(mockRepo)

	req := appuser.CreateUserRequest{
		Email: "existing@example.com",
		Name:  "Test User",
	}

	existingUser := domain.NewUser(req.Email, "Existing User")
	// Mock: 邮箱已存在
	mockRepo.On("FindByEmail", ctx, req.Email).Return(existingUser, nil)

	user, err := service.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}
