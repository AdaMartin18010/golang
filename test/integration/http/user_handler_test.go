package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	domainuser "github.com/yourusername/golang/internal/domain/user"
	"github.com/yourusername/golang/internal/interfaces/http/chi/handlers"
)

// MockUserService 模拟用户应用服务
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

// TestUserHandler_CreateUser 测试创建用户
func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "success",
			requestBody: map[string]string{
				"email": "test@example.com",
				"name":  "Test User",
			},
			mockSetup: func(m *MockUserService) {
				user := domainuser.NewUser("test@example.com", "Test User")
				m.On("CreateUser", mock.Anything, "test@example.com", "Test User").Return(user, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:        "invalid json",
			requestBody: "invalid json",
			mockSetup:   func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "empty email",
			requestBody: map[string]string{
				"email": "",
				"name":  "Test User",
			},
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "empty name",
			requestBody: map[string]string{
				"email": "test@example.com",
				"name":  "",
			},
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置 Mock
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			// 创建处理器
			handler := handlers.NewUserHandler(mockService)

			// 创建请求
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 执行请求
			handler.CreateUser(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

// TestUserHandler_GetUser 测试获取用户
func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "success",
			userID: "user-123",
			mockSetup: func(m *MockUserService) {
				user := domainuser.NewUser("test@example.com", "Test User")
				user.ID = "user-123"
				m.On("GetUser", mock.Anything, "user-123").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "user not found",
			userID: "non-existent",
			mockSetup: func(m *MockUserService) {
				m.On("GetUser", mock.Anything, "non-existent").Return(nil, domainuser.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置 Mock
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			// 创建处理器
			handler := handlers.NewUserHandler(mockService)

			// 创建请求（使用 Chi 路由）
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// 执行请求
			handler.GetUser(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}
