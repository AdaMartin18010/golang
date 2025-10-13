package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"domain"
	"usecase"
)

// UserHandler HTTP处理器
type UserHandler struct {
	userService *usecase.UserService
}

// NewUserHandler 创建新的用户处理器
func NewUserHandler(userService *usecase.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Response 通用响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CreateUser 创建用户处理器
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(req.Email, req.Name, req.Age)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, Response{
		Success: true,
		Data:    user,
	})
}

// GetUser 获取用户处理器
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, Response{
		Success: true,
		Data:    user,
	})
}

// GetAllUsers 获取所有用户处理器
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.userService.GetAllUsers()
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, Response{
		Success: true,
		Data:    users,
	})
}

// UpdateUser 更新用户处理器
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.UpdateUserProfile(userID, req.Name, req.Age)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, Response{
		Success: true,
		Data:    user,
	})
}

// DeleteUser 删除用户处理器
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, Response{
		Success: true,
		Data:    "User deleted successfully",
	})
}

// GetUsersByAgeRange 根据年龄范围获取用户处理器
func (h *UserHandler) GetUsersByAgeRange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	minAgeStr := r.URL.Query().Get("min_age")
	maxAgeStr := r.URL.Query().Get("max_age")

	minAge, err := strconv.Atoi(minAgeStr)
	if err != nil {
		http.Error(w, "Invalid min_age parameter", http.StatusBadRequest)
		return
	}

	maxAge, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		http.Error(w, "Invalid max_age parameter", http.StatusBadRequest)
		return
	}

	users, err := h.userService.GetUsersByAgeRange(minAge, maxAge)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.sendResponse(w, Response{
		Success: true,
		Data:    users,
	})
}

// handleError 处理错误响应
func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var errorMessage string

	switch err {
	case domain.ErrUserNotFound:
		statusCode = http.StatusNotFound
		errorMessage = "User not found"
	case domain.ErrUserAlreadyExists:
		statusCode = http.StatusConflict
		errorMessage = "User already exists"
	case domain.ErrInvalidInput:
		statusCode = http.StatusBadRequest
		errorMessage = "Invalid input data"
	default:
		statusCode = http.StatusInternalServerError
		errorMessage = "Internal server error"
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   errorMessage,
	})
}

// sendResponse 发送成功响应
func (h *UserHandler) sendResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
