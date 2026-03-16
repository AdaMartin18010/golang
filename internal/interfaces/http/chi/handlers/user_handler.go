// Package handlers provides HTTP handlers for user-related operations.
//
// HTTP 处理器负责：
// 1. 接收 HTTP 请求
// 2. 解析请求参数
// 3. 调用应用层服务
// 4. 格式化 HTTP 响应
//
// 设计原则：
// 1. 协议适配：将 HTTP 协议转换为应用层接口
// 2. 参数验证：验证 HTTP 请求参数
// 3. 错误处理：将应用层错误映射为 HTTP 状态码
// 4. 响应格式化：统一 API 响应格式
//
// 架构位置：
// - 位置：Interfaces Layer (internal/interfaces/http/chi/handlers/)
// - 职责：HTTP 协议适配、请求处理、响应格式化
// - 依赖：Application Layer（调用应用服务）
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	appuser "github.com/yourusername/golang/internal/app/user"
	domainuser "github.com/yourusername/golang/internal/domain/user"
	apperrors "github.com/yourusername/golang/pkg/errors"
)

// UserService 用户应用服务接口（用于依赖注入和测试）。
type UserService interface {
	CreateUser(ctx context.Context, email, name string) (*domainuser.User, error)
	GetUser(ctx context.Context, id string) (*domainuser.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*domainuser.User, error)
	UpdateUserName(ctx context.Context, id, name string) error
	DeleteUser(ctx context.Context, id string) error
}

// UserHandler 用户 HTTP 处理器。
type UserHandler struct {
	service UserService
}

// NewUserHandler 创建用户 HTTP 处理器。
func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// CreateUserRequest 创建用户请求。
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// CreateUser 创建用户。
// POST /api/v1/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, apperrors.NewInvalidInputError("invalid request body"))
		return
	}

	if req.Email == "" {
		Error(w, http.StatusBadRequest, apperrors.NewInvalidInputError("email is required"))
		return
	}
	if req.Name == "" {
		Error(w, http.StatusBadRequest, apperrors.NewInvalidInputError("name is required"))
		return
	}

	user, err := h.service.CreateUser(ctx, req.Email, req.Name)
	if err != nil {
		if errors.Is(err, appuser.ErrUserAlreadyExists) {
			Error(w, http.StatusConflict, apperrors.NewConflictError("user already exists"))
			return
		}
		Error(w, http.StatusInternalServerError, apperrors.NewInternalError("failed to create user", err))
		return
	}

	Success(w, http.StatusCreated, user)
}

// GetUser 获取用户。
// GET /api/v1/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		Error(w, http.StatusBadRequest, apperrors.NewInvalidInputError("user id is required"))
		return
	}

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, domainuser.ErrUserNotFound) {
			Error(w, http.StatusNotFound, apperrors.NewNotFoundError("user", id))
			return
		}
		Error(w, http.StatusInternalServerError, apperrors.NewInternalError("failed to get user", err))
		return
	}

	Success(w, http.StatusOK, user)
}

// ListUsersRequest 列出用户请求。
type ListUsersRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ListUsers 列出用户。
// GET /api/v1/users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 从查询参数获取分页信息
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.service.ListUsers(ctx, limit, offset)
	if err != nil {
		Error(w, http.StatusInternalServerError, apperrors.NewInternalError("failed to list users", err))
		return
	}

	Success(w, http.StatusOK, users)
}

// UpdateUserRequest 更新用户请求。
type UpdateUserRequest struct {
	Name string `json:"name"`
}

// UpdateUser 更新用户。
// PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		Error(w, http.StatusBadRequest, apperrors.NewInvalidInputError("user id is required"))
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, apperrors.NewInvalidInputError("invalid request body"))
		return
	}

	if req.Name == "" {
		Error(w, http.StatusBadRequest, apperrors.NewInvalidInputError("name is required"))
		return
	}

	if err := h.service.UpdateUserName(ctx, id, req.Name); err != nil {
		if errors.Is(err, domainuser.ErrUserNotFound) {
			Error(w, http.StatusNotFound, apperrors.NewNotFoundError("user", id))
			return
		}
		Error(w, http.StatusInternalServerError, apperrors.NewInternalError("failed to update user", err))
		return
	}

	// 获取更新后的用户
	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		Error(w, http.StatusInternalServerError, apperrors.NewInternalError("failed to get updated user", err))
		return
	}

	Success(w, http.StatusOK, user)
}

// DeleteUser 删除用户。
// DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		Error(w, http.StatusBadRequest, apperrors.NewInvalidInputError("user id is required"))
		return
	}

	if err := h.service.DeleteUser(ctx, id); err != nil {
		if errors.Is(err, domainuser.ErrUserNotFound) {
			Error(w, http.StatusNotFound, apperrors.NewNotFoundError("user", id))
			return
		}
		Error(w, http.StatusInternalServerError, apperrors.NewInternalError("failed to delete user", err))
		return
	}

	Success(w, http.StatusNoContent, nil)
}

