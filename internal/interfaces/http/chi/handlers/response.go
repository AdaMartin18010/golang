package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/golang/pkg/errors"
)

// APIResponse 统一 API 响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError API 错误信息
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success 成功响应
func Success(w http.ResponseWriter, code int, data interface{}) {
	response := APIResponse{
		Code:    code,
		Message: "success",
		Data:    data,
	}
	writeJSON(w, code, response)
}

// Error 错误响应
func Error(w http.ResponseWriter, code int, err error) {
	response := APIResponse{
		Code:    code,
		Message: "error",
	}

	// 检查是否是 AppError
	if appErr, ok := err.(*errors.AppError); ok {
		response.Error = &APIError{
			Code:    string(appErr.Code),
			Message: appErr.Message,
		}
	} else {
		response.Error = &APIError{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}
	}

	writeJSON(w, code, response)
}

// writeJSON 写入 JSON 响应
func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
