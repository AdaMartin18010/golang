package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/yourusername/golang/internal/domain/user"
	"github.com/yourusername/golang/test/mocks"
)

// ServiceTestSuite 用户服务测试套件
type ServiceTestSuite struct {
	suite.Suite
	ctx      context.Context
	repo     *mocks.MockRepository[user.User]
	service  *Service
}

// SetupTest 每个测试前执行
func (s *ServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = &mocks.MockRepository[user.User]{}
	s.service = &Service{
		repo: s.repo,
	}
}

// TearDownTest 每个测试后执行
func (s *ServiceTestSuite) TearDownTest() {
	s.repo.AssertExpectations(s.T())
}

// TestCreateUser 测试创建用户
func (s *ServiceTestSuite) TestCreateUser() {
	// 准备测试数据
	testUser := &user.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// 设置 mock 期望
	s.repo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Return(nil)

	// 执行测试
	err := s.service.CreateUser(s.ctx, testUser)

	// 验证结果
	s.NoError(err)
}

// TestCreateUser_ValidationError 测试创建用户 - 验证错误
func (s *ServiceTestSuite) TestCreateUser_ValidationError() {
	// 无效的用户（空邮箱）
	invalidUser := &user.User{
		Email: "",
		Name:  "Test User",
	}

	// 不应该调用 repository
	// s.repo.On(...) 不设置期望

	// 执行测试
	err := s.service.CreateUser(s.ctx, invalidUser)

	// 验证结果
	s.Error(err)
	s.Contains(err.Error(), "email")
}

// TestCreateUser_RepositoryError 测试创建用户 - 仓储错误
func (s *ServiceTestSuite) TestCreateUser_RepositoryError() {
	testUser := &user.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// 设置 mock 返回错误
	expectedErr := errors.New("database error")
	s.repo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Return(expectedErr)

	// 执行测试
	err := s.service.CreateUser(s.ctx, testUser)

	// 验证结果
	s.Error(err)
	s.Equal(expectedErr, err)
}

// TestFindUserByID 测试根据ID查找用户
func (s *ServiceTestSuite) TestFindUserByID() {
	expectedUser := &user.User{
		ID:    "123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// 设置 mock 期望
	s.repo.On("FindByID", s.ctx, "123").Return(expectedUser, nil)

	// 执行测试
	result, err := s.service.FindUserByID(s.ctx, "123")

	// 验证结果
	s.NoError(err)
	s.Equal(expectedUser, result)
}

// TestFindUserByID_NotFound 测试查找用户 - 未找到
func (s *ServiceTestSuite) TestFindUserByID_NotFound() {
	// 设置 mock 返回 nil
	s.repo.On("FindByID", s.ctx, "999").Return(nil, nil)

	// 执行测试
	result, err := s.service.FindUserByID(s.ctx, "999")

	// 验证结果
	s.NoError(err)
	s.Nil(result)
}

// 运行测试套件
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

// 表格驱动测试示例
func TestService_UpdateUser_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		updates map[string]interface{}
		setup   func(*mocks.MockRepository[user.User])
		wantErr bool
	}{
		{
			name:   "successful update",
			userID: "123",
			updates: map[string]interface{}{
				"name": "Updated Name",
			},
			setup: func(repo *mocks.MockRepository[user.User]) {
				existingUser := &user.User{
					ID:    "123",
					Email: "test@example.com",
					Name:  "Old Name",
				}
				repo.On("FindByID", mock.Anything, "123").Return(existingUser, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: "999",
			updates: map[string]interface{}{
				"name": "Updated Name",
			},
			setup: func(repo *mocks.MockRepository[user.User]) {
				repo.On("FindByID", mock.Anything, "999").Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 mock
			repo := &mocks.MockRepository[user.User]{}
			if tt.setup != nil {
				tt.setup(repo)
			}

			service := &Service{repo: repo}

			// 执行测试
			err := service.UpdateUser(context.Background(), tt.userID, tt.updates)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// 验证 mock 调用
			repo.AssertExpectations(t)
		})
	}
}
