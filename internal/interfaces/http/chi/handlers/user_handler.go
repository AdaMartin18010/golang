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
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/pkg/errors"
)

// UserService 用户应用服务接口（用于依赖注入和测试）。
//
// 功能说明：
// - 定义用户应用服务的接口
// - 便于测试时使用 Mock
// - 支持依赖注入
type UserService interface {
	CreateUser(ctx context.Context, req appuser.CreateUserRequest) (*appuser.UserResponse, error)
	GetUser(ctx context.Context, id string) (*appuser.UserResponse, error)
	ListUsers(ctx context.Context, req appuser.ListUsersRequest) (*appuser.ListUsersResponse, error)
	UpdateUser(ctx context.Context, id string, req appuser.UpdateUserRequest) (*appuser.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

// UserHandler 用户 HTTP 处理器。
//
// 功能说明：
// - 处理用户相关的 HTTP 请求
// - 调用应用层服务完成业务逻辑
// - 格式化 HTTP 响应
//
// 字段说明：
// - service: 用户应用服务（来自 Application Layer）
//
// 职责：
// 1. 接收 HTTP 请求
// 2. 解析请求参数（JSON、URL 参数等）
// 3. 调用应用服务
// 4. 处理错误并返回适当的 HTTP 状态码
// 5. 格式化响应
type UserHandler struct {
	service UserService
}

// NewUserHandler 创建用户 HTTP 处理器。
//
// 功能说明：
// - 接收应用层服务实例
// - 创建并返回配置好的处理器
//
// 参数：
// - service: 用户应用服务实例（可以是 *appuser.Service 或实现 UserService 接口的类型）
//
// 返回：
// - *UserHandler: 配置好的处理器实例
//
// 使用示例：
//
//	userService := appuser.NewService(userRepo)
//	handler := handlers.NewUserHandler(userService)
func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// CreateUserRequest 创建用户请求。
//
// 字段说明：
// - Email: 用户邮箱（必填）
// - Name: 用户名称（必填）
type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=2"`
}

// CreateUser 创建用户。
//
// 功能说明：
// 1. 解析请求体（JSON）
// 2. 验证请求参数
// 3. 调用应用服务创建用户
// 4. 返回创建的用户信息
//
// HTTP 方法：POST
// 路径：/api/v1/users
// 请求体：JSON 格式的 CreateUserRequest
//
// 响应：
// - 201 Created: 用户创建成功
// - 400 Bad Request: 请求参数无效
// - 409 Conflict: 用户已存在
// - 500 Internal Server Error: 服务器内部错误
//
// 使用示例：
//
//	POST /api/v1/users
//	Content-Type: application/json
//
//	{
//	  "email": "test@example.com",
//	  "name": "Test User"
//	}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. 解析请求体
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("invalid request body"))
		return
	}

	// 2. 验证请求参数
	if req.Email == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("email is required"))
		return
	}
	if req.Name == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("name is required"))
		return
	}

	// 3. 调用应用服务
	appReq := appuser.CreateUserRequest{
		Email: req.Email,
		Name:  req.Name,
	}
	user, err := h.service.CreateUser(ctx, appReq)
	if err != nil {
		// 4. 处理错误
		if err == appuser.ErrUserAlreadyExists {
			Error(w, http.StatusConflict, errors.NewConflictError("user already exists"))
			return
		}
		if err == appuser.ErrInvalidInput {
			Error(w, http.StatusBadRequest, errors.NewInvalidInputError("invalid input"))
			return
		}
		Error(w, http.StatusInternalServerError, errors.NewInternalError("failed to create user", err))
		return
	}

	// 5. 返回成功响应
	Success(w, http.StatusCreated, user)
}

// GetUser 获取用户。
//
// 功能说明：
// 1. 从 URL 路径提取用户 ID
// 2. 调用应用服务获取用户
// 3. 返回用户信息
//
// HTTP 方法：GET
// 路径：/api/v1/users/{id}
//
// 响应：
// - 200 OK: 用户信息
// - 404 Not Found: 用户不存在
// - 500 Internal Server Error: 服务器内部错误
//
// 使用示例：
//
//	GET /api/v1/users/123
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. 从 URL 路径提取用户 ID
	id := chi.URLParam(r, "id")
	if id == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("user id is required"))
		return
	}

	// 2. 调用应用服务
	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		// 3. 处理错误
		if err == appuser.ErrUserNotFound {
			Error(w, http.StatusNotFound, errors.NewNotFoundError("user", id))
			return
		}
		Error(w, http.StatusInternalServerError, errors.NewInternalError("failed to get user", err))
		return
	}

	// 4. 返回成功响应
	Success(w, http.StatusOK, user)
}

// ListUsers 列出用户。
//
// 功能说明：
// 1. 解析查询参数（分页、过滤等）
// 2. 调用应用服务获取用户列表
// 3. 返回用户列表
//
// HTTP 方法：GET
// 路径：/api/v1/users
// 查询参数：
// - limit: 每页数量（可选，默认 10，最大 100）
// - offset: 偏移量（可选，默认 0）
//
// 响应：
// - 200 OK: 用户列表
// - 500 Internal Server Error: 服务器内部错误
//
// 使用示例：
//
//	GET /api/v1/users?limit=10&offset=0
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. 解析查询参数
	limit := 10
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100
			}
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := parseInt(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// 2. 调用应用服务
	appReq := appuser.ListUsersRequest{
		Limit:  limit,
		Offset: offset,
	}
	resp, err := h.service.ListUsers(ctx, appReq)
	if err != nil {
		Error(w, http.StatusInternalServerError, errors.NewInternalError("failed to list users", err))
		return
	}

	// 3. 返回成功响应
	Success(w, http.StatusOK, resp)
}

// UpdateUserRequest 更新用户请求。
//
// 字段说明：
// - Name: 用户名称（可选）
// - Email: 用户邮箱（可选）
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// UpdateUser 更新用户。
//
// 功能说明：
// 1. 从 URL 路径提取用户 ID
// 2. 解析请求体（JSON）
// 3. 调用应用服务更新用户
// 4. 返回更新后的用户信息
//
// HTTP 方法：PUT
// 路径：/api/v1/users/{id}
// 请求体：JSON 格式的 UpdateUserRequest
//
// 响应：
// - 200 OK: 用户更新成功
// - 400 Bad Request: 请求参数无效
// - 404 Not Found: 用户不存在
// - 409 Conflict: 邮箱已被使用
// - 500 Internal Server Error: 服务器内部错误
//
// 使用示例：
//
//	PUT /api/v1/users/123
//	Content-Type: application/json
//
//	{
//	  "name": "Updated Name"
//	}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. 从 URL 路径提取用户 ID
	id := chi.URLParam(r, "id")
	if id == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("user id is required"))
		return
	}

	// 2. 解析请求体
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("invalid request body"))
		return
	}

	// 3. 验证至少有一个字段要更新
	if req.Name == nil && req.Email == nil {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("at least one field (name or email) is required"))
		return
	}

	// 4. 调用应用服务
	appReq := appuser.UpdateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	}
	user, err := h.service.UpdateUser(ctx, id, appReq)
	if err != nil {
		// 5. 处理错误
		if err == appuser.ErrUserNotFound {
			Error(w, http.StatusNotFound, errors.NewNotFoundError("user", id))
			return
		}
		if err == appuser.ErrUserAlreadyExists {
			Error(w, http.StatusConflict, errors.NewConflictError("email already exists"))
			return
		}
		if err == appuser.ErrInvalidInput {
			Error(w, http.StatusBadRequest, errors.NewInvalidInputError("invalid input"))
			return
		}
		Error(w, http.StatusInternalServerError, errors.NewInternalError("failed to update user", err))
		return
	}

	// 6. 返回成功响应
	Success(w, http.StatusOK, user)
}

// DeleteUser 删除用户。
//
// 功能说明：
// 1. 从 URL 路径提取用户 ID
// 2. 调用应用服务删除用户
// 3. 返回删除结果
//
// HTTP 方法：DELETE
// 路径：/api/v1/users/{id}
//
// 响应：
// - 204 No Content: 用户删除成功
// - 404 Not Found: 用户不存在
// - 500 Internal Server Error: 服务器内部错误
//
// 使用示例：
//
//	DELETE /api/v1/users/123
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. 从 URL 路径提取用户 ID
	id := chi.URLParam(r, "id")
	if id == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("user id is required"))
		return
	}

	// 2. 调用应用服务
	err := h.service.DeleteUser(ctx, id)
	if err != nil {
		// 3. 处理错误
		if err == appuser.ErrUserNotFound {
			Error(w, http.StatusNotFound, errors.NewNotFoundError("user", id))
			return
		}
		Error(w, http.StatusInternalServerError, errors.NewInternalError("failed to delete user", err))
		return
	}

	// 4. 返回成功响应（204 No Content）
	w.WriteHeader(http.StatusNoContent)
}

// parseInt 解析字符串为整数（辅助函数）。
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
