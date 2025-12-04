package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/golang/internal/domain/user"
)

// MockUserRepository 是用户仓储的mock实现
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

// TestNewService 测试创建服务
func TestNewService(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewService(repo)

	assert.NotNil(t, service)
}

// TestService_GetUser 测试获取用户
func TestService_GetUser(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 准备测试数据
	testUser := user.NewUser("test@example.com", "Test User")

	// 设置mock期望
	repo.On("FindByID", ctx, testUser.ID).Return(testUser, nil)

	// 执行测试
	result, err := service.GetUser(ctx, testUser.ID)

	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, result.ID)
	assert.Equal(t, testUser.Email, result.Email)

	// 验证mock被调用
	repo.AssertExpectations(t)
}

// TestService_GetUser_NotFound 测试获取不存在的用户
func TestService_GetUser_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 设置mock期望 - 返回错误
	repo.On("FindByID", ctx, "non-existent").Return(nil, errors.New("user not found"))

	// 执行测试
	result, err := service.GetUser(ctx, "non-existent")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

// TestService_CreateUser 测试创建用户
func TestService_CreateUser(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	email := "newuser@example.com"
	name := "New User"

	// 设置mock期望 - 邮箱不存在
	repo.On("FindByEmail", ctx, email).Return(nil, errors.New("not found"))
	// 设置mock期望 - 保存成功
	repo.On("Save", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	// 执行测试
	result, err := service.CreateUser(ctx, email, name)

	// 验证结果
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, email, result.Email)
	assert.Equal(t, name, result.Name)

	repo.AssertExpectations(t)
}

// TestService_CreateUser_DuplicateEmail 测试创建重复邮箱用户
func TestService_CreateUser_DuplicateEmail(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	email := "existing@example.com"
	existingUser := user.NewUser(email, "Existing User")

	// 设置mock期望 - 邮箱已存在
	repo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	// 执行测试
	result, err := service.CreateUser(ctx, email, "New User")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")

	repo.AssertExpectations(t)
}

// TestService_UpdateUser 测试更新用户
func TestService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	testUser := user.NewUser("test@example.com", "Old Name")

	// 设置mock期望
	repo.On("FindByID", ctx, testUser.ID).Return(testUser, nil)
	repo.On("Update", ctx, testUser).Return(nil)

	// 执行测试
	err := service.UpdateUserName(ctx, testUser.ID, "New Name")

	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, "New Name", testUser.Name)

	repo.AssertExpectations(t)
}

// TestService_DeleteUser 测试删除用户
func TestService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	userID := "user-to-delete"

	// 设置mock期望
	repo.On("Delete", ctx, userID).Return(nil)

	// 执行测试
	err := service.DeleteUser(ctx, userID)

	// 验证结果
	require.NoError(t, err)

	repo.AssertExpectations(t)
}

// TestService_ListUsers 测试列出用户
func TestService_ListUsers(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	users := []*user.User{
		user.NewUser("user1@example.com", "User 1"),
		user.NewUser("user2@example.com", "User 2"),
		user.NewUser("user3@example.com", "User 3"),
	}

	// 设置mock期望
	repo.On("List", ctx, 10, 0).Return(users, nil)

	// 执行测试
	result, err := service.ListUsers(ctx, 10, 0)

	// 验证结果
	require.NoError(t, err)
	assert.Len(t, result, 3)

	repo.AssertExpectations(t)
}

// BenchmarkService_GetUser 性能测试 - 获取用户
func BenchmarkService_GetUser(b *testing.B) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	testUser := user.NewUser("test@example.com", "Test User")
	repo.On("FindByID", ctx, testUser.ID).Return(testUser, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetUser(ctx, testUser.ID)
	}
}
