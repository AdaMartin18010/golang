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
	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/internal/interfaces/http/chi/handlers"
)

// TestUserHandler_ListUsers 测试列出用户
func TestUserHandler_ListUsers(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "success",
			queryParams: "limit=10&offset=0",
			mockSetup: func(m *MockUserService) {
				m.On("ListUsers", mock.Anything, appuser.ListUsersRequest{
					Limit:  10,
					Offset: 0,
				}).Return(&appuser.ListUsersResponse{
					Users: []*appuser.UserResponse{
						{
							ID:        "user-1",
							Email:     "user1@example.com",
							Name:      "User 1",
							CreatedAt: "2025-01-01T00:00:00Z",
							UpdatedAt: "2025-01-01T00:00:00Z",
						},
					},
					Total:  1,
					Limit:  10,
					Offset: 0,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "default pagination",
			queryParams: "",
			mockSetup: func(m *MockUserService) {
				m.On("ListUsers", mock.Anything, appuser.ListUsersRequest{
					Limit:  10,
					Offset: 0,
				}).Return(&appuser.ListUsersResponse{
					Users: []*appuser.UserResponse{},
					Total:  0,
					Limit:  10,
					Offset: 0,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.mockSetup(mockService)
			handler := handlers.NewUserHandler(mockService)

			url := "/users"
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			handler.ListUsers(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if !tt.expectedError {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, "success", response["message"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestUserHandler_UpdateUser 测试更新用户
func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "success - update name",
			userID: "user-123",
			requestBody: map[string]string{
				"name": "New Name",
			},
			mockSetup: func(m *MockUserService) {
				m.On("UpdateUser", mock.Anything, "user-123", mock.MatchedBy(func(req appuser.UpdateUserRequest) bool {
					return req.Name != nil && *req.Name == "New Name"
				})).Return(&appuser.UserResponse{
					ID:        "user-123",
					Email:     "test@example.com",
					Name:      "New Name",
					CreatedAt: "2025-01-01T00:00:00Z",
					UpdatedAt: "2025-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "user not found",
			userID: "non-existent",
			requestBody: map[string]string{
				"name": "New Name",
			},
			mockSetup: func(m *MockUserService) {
				m.On("UpdateUser", mock.Anything, "non-existent", mock.Anything).Return(nil, appuser.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:   "empty id",
			userID: "",
			requestBody: map[string]string{
				"name": "New Name",
			},
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.mockSetup(mockService)
			handler := handlers.NewUserHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.userID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.UpdateUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

// TestUserHandler_DeleteUser 测试删除用户
func TestUserHandler_DeleteUser(t *testing.T) {
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
				m.On("DeleteUser", mock.Anything, "user-123").Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name:   "user not found",
			userID: "non-existent",
			mockSetup: func(m *MockUserService) {
				m.On("DeleteUser", mock.Anything, "non-existent").Return(appuser.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:           "empty id",
			userID:         "",
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.mockSetup(mockService)
			handler := handlers.NewUserHandler(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.DeleteUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
