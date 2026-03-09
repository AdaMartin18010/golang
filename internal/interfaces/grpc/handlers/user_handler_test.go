// Package handlers provides tests for gRPC user handlers.
package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/golang/internal/application/user"
	domainuser "github.com/yourusername/golang/internal/domain/user"
	userpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/userpb"
)

// MockUserService is a mock for user.Service
type MockUserAppService struct {
	mock.Mock
}

func (m *MockUserAppService) GetUser(ctx context.Context, id string) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *MockUserAppService) CreateUser(ctx context.Context, email, name string) (*domainuser.User, error) {
	args := m.Called(ctx, email, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *MockUserAppService) UpdateUserName(ctx context.Context, id, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func (m *MockUserAppService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserAppService) ListUsers(ctx context.Context, limit, offset int) ([]*domainuser.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainuser.User), args.Error(1)
}

createTestUser := func(id, email, name string) *domainuser.User {
	return &domainuser.User{
		ID:        id,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestNewUserHandler(t *testing.T) {
	mockService := new(MockUserAppService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	handler := NewUserHandler((*user.Service)(mockService), logger)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
	assert.NotNil(t, handler.logger)
}

func TestNewUserHandler_NilLogger(t *testing.T) {
	mockService := new(MockUserAppService)

	handler := NewUserHandler((*user.Service)(mockService), nil)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.logger)
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	domainUser := createTestUser("123", "test@example.com", "Test User")
	mockService.On("GetUser", mock.Anything, "123").Return(domainUser, nil)

	req := &userpb.GetUserRequest{Id: "123"}
	resp, err := handler.GetUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "123", resp.User.Id)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Equal(t, "Test User", resp.User.Name)
	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_EmptyID(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	req := &userpb.GetUserRequest{Id: ""}
	resp, err := handler.GetUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "user ID is required")
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("GetUser", mock.Anything, "123").Return(nil, domainuser.ErrUserNotFound)

	req := &userpb.GetUserRequest{Id: "123"}
	resp, err := handler.GetUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "user not found")
	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_InternalError(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("GetUser", mock.Anything, "123").Return(nil, errors.New("database error"))

	req := &userpb.GetUserRequest{Id: "123"}
	resp, err := handler.GetUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	domainUser := createTestUser("123", "test@example.com", "Test User")
	mockService.On("CreateUser", mock.Anything, "test@example.com", "Test User").Return(domainUser, nil)

	req := &userpb.CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}
	resp, err := handler.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "123", resp.User.Id)
	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_EmptyEmail(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	req := &userpb.CreateUserRequest{
		Email: "",
		Name:  "Test User",
	}
	resp, err := handler.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUserHandler_CreateUser_EmptyName(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	req := &userpb.CreateUserRequest{
		Email: "test@example.com",
		Name:  "",
	}
	resp, err := handler.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUserHandler_CreateUser_InvalidEmail(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("CreateUser", mock.Anything, "invalid", "Test").Return(nil, domainuser.ErrInvalidEmailFormat)

	req := &userpb.CreateUserRequest{
		Email: "invalid",
		Name:  "Test",
	}
	resp, err := handler.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_NameTooShort(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("CreateUser", mock.Anything, "test@example.com", "A").Return(nil, domainuser.ErrNameTooShort)

	req := &userpb.CreateUserRequest{
		Email: "test@example.com",
		Name:  "A",
	}
	resp, err := handler.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_AlreadyExists(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("CreateUser", mock.Anything, "test@example.com", "Test User").Return(nil, domainuser.ErrUserAlreadyExists)

	req := &userpb.CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}
	resp, err := handler.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_InternalError(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("CreateUser", mock.Anything, "test@example.com", "Test User").Return(nil, errors.New("database error"))

	req := &userpb.CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}
	resp, err := handler.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	domainUser := createTestUser("123", "test@example.com", "Original Name")
	mockService.On("GetUser", mock.Anything, "123").Return(domainUser, nil)
	mockService.On("UpdateUserName", mock.Anything, "123", "Updated Name").Return(nil)

	req := &userpb.UpdateUserRequest{
		Id:   "123",
		Name: "Updated Name",
	}
	resp, err := handler.UpdateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "123", resp.User.Id)
	assert.Equal(t, "Updated Name", resp.User.Name)
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_EmptyID(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	req := &userpb.UpdateUserRequest{
		Id:   "",
		Name: "Updated Name",
	}
	resp, err := handler.UpdateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUserHandler_UpdateUser_UserNotFound(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("GetUser", mock.Anything, "123").Return(nil, domainuser.ErrUserNotFound)

	req := &userpb.UpdateUserRequest{
		Id:   "123",
		Name: "Updated Name",
	}
	resp, err := handler.UpdateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_UpdateNameError(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	domainUser := createTestUser("123", "test@example.com", "Original Name")
	mockService.On("GetUser", mock.Anything, "123").Return(domainUser, nil)
	mockService.On("UpdateUserName", mock.Anything, "123", "A").Return(domainuser.ErrNameTooShort)

	req := &userpb.UpdateUserRequest{
		Id:   "123",
		Name: "A",
	}
	resp, err := handler.UpdateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_InternalError(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("GetUser", mock.Anything, "123").Return(nil, errors.New("database error"))

	req := &userpb.UpdateUserRequest{
		Id:   "123",
		Name: "Updated Name",
	}
	resp, err := handler.UpdateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_UpdateEmail(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	domainUser := createTestUser("123", "old@example.com", "Test User")
	mockService.On("GetUser", mock.Anything, "123").Return(domainUser, nil)

	req := &userpb.UpdateUserRequest{
		Id:    "123",
		Email: "new@example.com",
	}
	resp, err := handler.UpdateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// Email should be updated locally
	assert.Equal(t, "new@example.com", resp.User.Email)
	mockService.AssertExpectations(t)
}

func TestUserHandler_DeleteUser_Success(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("DeleteUser", mock.Anything, "123").Return(nil)

	req := &userpb.DeleteUserRequest{Id: "123"}
	resp, err := handler.DeleteUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	mockService.AssertExpectations(t)
}

func TestUserHandler_DeleteUser_EmptyID(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	req := &userpb.DeleteUserRequest{Id: ""}
	resp, err := handler.DeleteUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUserHandler_DeleteUser_NotFound(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("DeleteUser", mock.Anything, "123").Return(domainuser.ErrUserNotFound)

	req := &userpb.DeleteUserRequest{Id: "123"}
	resp, err := handler.DeleteUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	mockService.AssertExpectations(t)
}

func TestUserHandler_DeleteUser_InternalError(t *testing.T) {
	mockService := new(MockUserAppService)
	handler := NewUserHandler((*user.Service)(mockService), nil)

	mockService.On("DeleteUser", mock.Anything, "123").Return(errors.New("database error"))

	req := &userpb.DeleteUserRequest{Id: "123"}
	resp, err := handler.DeleteUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	mockService.AssertExpectations(t)
}

func TestToProtoUser(t *testing.T) {
	domainUser := createTestUser("123", "test@example.com", "Test User")

	protoUser := toProtoUser(domainUser)

	assert.NotNil(t, protoUser)
	assert.Equal(t, "123", protoUser.Id)
	assert.Equal(t, "test@example.com", protoUser.Email)
	assert.Equal(t, "Test User", protoUser.Name)
	assert.NotNil(t, protoUser.CreatedAt)
	assert.NotNil(t, protoUser.UpdatedAt)
}

func TestToProtoUser_Nil(t *testing.T) {
	protoUser := toProtoUser(nil)

	assert.Nil(t, protoUser)
}
