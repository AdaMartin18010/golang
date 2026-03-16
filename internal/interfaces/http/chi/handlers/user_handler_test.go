// Package handlers provides tests for user HTTP handlers.
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	appuser "github.com/yourusername/golang/internal/app/user"
	domainuser "github.com/yourusername/golang/internal/domain/user"
)

// MockUserService is a mock implementation of UserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, email, name string) (*domainuser.User, error) {
	args := m.Called(ctx, email, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, id string) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *MockUserService) ListUsers(ctx context.Context, limit, offset int) ([]*domainuser.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainuser.User), args.Error(1)
}

func (m *MockUserService) UpdateUserName(ctx context.Context, id, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func createTestUser(id, email, name string) *domainuser.User {
	return &domainuser.User{
		ID:        id,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestNewUserHandler(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func TestCreateUserRequestStruct(t *testing.T) {
	req := CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "Test User", req.Name)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	user := createTestUser("123", "test@example.com", "Test User")
	mockService.On("CreateUser", mock.Anything, "test@example.com", "Test User").Return(user, nil)

	reqBody := `{"email": "test@example.com", "name": "Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_InvalidJSON(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	reqBody := `{"invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_CreateUser_EmptyEmail(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	reqBody := `{"email": "", "name": "Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_CreateUser_EmptyName(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	reqBody := `{"email": "test@example.com", "name": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_CreateUser_UserAlreadyExists(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("CreateUser", mock.Anything, "test@example.com", "Test User").Return(nil, appuser.ErrUserAlreadyExists)

	reqBody := `{"email": "test@example.com", "name": "Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_InternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("CreateUser", mock.Anything, "test@example.com", "Test User").Return(nil, errors.New("database error"))

	reqBody := `{"email": "test@example.com", "name": "Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	user := createTestUser("123", "test@example.com", "Test User")
	mockService.On("GetUser", mock.Anything, "123").Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.GetUser(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_EmptyID(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "")
	rec := httptest.NewRecorder()

	handler.GetUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("GetUser", mock.Anything, "123").Return(nil, domainuser.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.GetUser(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_InternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("GetUser", mock.Anything, "123").Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.GetUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_ListUsers_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	users := []*domainuser.User{
		createTestUser("1", "user1@example.com", "User 1"),
		createTestUser("2", "user2@example.com", "User 2"),
	}
	mockService.On("ListUsers", mock.Anything, 10, 0).Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()

	handler.ListUsers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response APIResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Message)
	mockService.AssertExpectations(t)
}

func TestUserHandler_ListUsers_WithPagination(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	users := []*domainuser.User{
		createTestUser("3", "user3@example.com", "User 3"),
	}
	mockService.On("ListUsers", mock.Anything, 5, 10).Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=5&offset=10", nil)
	rec := httptest.NewRecorder()

	handler.ListUsers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_ListUsers_InvalidPagination(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	users := []*domainuser.User{}
	mockService.On("ListUsers", mock.Anything, 10, 0).Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=invalid&offset=invalid", nil)
	rec := httptest.NewRecorder()

	handler.ListUsers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_ListUsers_InternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("ListUsers", mock.Anything, 10, 0).Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()

	handler.ListUsers(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	user := createTestUser("123", "test@example.com", "Updated Name")
	mockService.On("UpdateUserName", mock.Anything, "123", "Updated Name").Return(nil)
	mockService.On("GetUser", mock.Anything, "123").Return(user, nil)

	reqBody := `{"name": "Updated Name"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/123", bytes.NewBufferString(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.UpdateUser(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_EmptyID(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	reqBody := `{"name": "Updated Name"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/", bytes.NewBufferString(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "")
	rec := httptest.NewRecorder()

	handler.UpdateUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_UpdateUser_InvalidJSON(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	reqBody := `{"invalid json`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/123", bytes.NewBufferString(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.UpdateUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_UpdateUser_EmptyName(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	reqBody := `{"name": ""}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/123", bytes.NewBufferString(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.UpdateUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_UpdateUser_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("UpdateUserName", mock.Anything, "123", "Updated Name").Return(domainuser.ErrUserNotFound)

	reqBody := `{"name": "Updated Name"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/123", bytes.NewBufferString(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.UpdateUser(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_GetUserError(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("UpdateUserName", mock.Anything, "123", "Updated Name").Return(nil)
	mockService.On("GetUser", mock.Anything, "123").Return(nil, errors.New("database error"))

	reqBody := `{"name": "Updated Name"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/123", bytes.NewBufferString(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.UpdateUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_DeleteUser_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("DeleteUser", mock.Anything, "123").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.DeleteUser(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_DeleteUser_EmptyID(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "")
	rec := httptest.NewRecorder()

	handler.DeleteUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_DeleteUser_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("DeleteUser", mock.Anything, "123").Return(domainuser.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.DeleteUser(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_DeleteUser_InternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("DeleteUser", mock.Anything, "123").Return(errors.New("database error"))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "123")
	rec := httptest.NewRecorder()

	handler.DeleteUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateUserRequestStruct(t *testing.T) {
	req := UpdateUserRequest{
		Name: "Updated Name",
	}

	assert.Equal(t, "Updated Name", req.Name)
}

