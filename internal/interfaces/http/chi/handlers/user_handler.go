package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/pkg/errors"
)

// UserHandler 用户处理器
type UserHandler struct {
	service appuser.Service
}

// NewUserHandler 创建用户处理器
func NewUserHandler(service appuser.Service) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req appuser.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Invalid request body"))
		return
	}

	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.ErrCodeConflict {
			Error(w, http.StatusConflict, err)
		} else {
			Error(w, http.StatusInternalServerError, err)
		}
		return
	}

	Success(w, http.StatusCreated, user)
}

// GetUser 获取用户
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("User ID is required"))
		return
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.ErrCodeNotFound {
			Error(w, http.StatusNotFound, err)
		} else {
			Error(w, http.StatusInternalServerError, err)
		}
		return
	}

	Success(w, http.StatusOK, user)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("User ID is required"))
		return
	}

	var req appuser.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Invalid request body"))
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			switch appErr.Code {
			case errors.ErrCodeNotFound:
				Error(w, http.StatusNotFound, err)
			case errors.ErrCodeConflict:
				Error(w, http.StatusConflict, err)
			default:
				Error(w, http.StatusInternalServerError, err)
			}
		} else {
			Error(w, http.StatusInternalServerError, err)
		}
		return
	}

	Success(w, http.StatusOK, user)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("User ID is required"))
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.ErrCodeNotFound {
			Error(w, http.StatusNotFound, err)
		} else {
			Error(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListUsers 列出用户
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// 解析分页参数
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.service.ListUsers(r.Context(), limit, offset)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	Success(w, http.StatusOK, users)
}
