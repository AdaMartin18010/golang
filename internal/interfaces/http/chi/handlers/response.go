// Package handlers provides HTTP response helpers for the interfaces layer.
//
// HTTP 响应处理器负责：
// 1. 统一 API 响应格式
// 2. 处理成功和错误响应
// 3. 将应用层错误映射为 HTTP 状态码
// 4. 提供 JSON 响应格式化
//
// 设计原则：
// 1. 统一格式：所有 API 响应使用相同的格式
// 2. 错误处理：统一的错误响应格式
// 3. 类型安全：使用结构体定义响应格式
// 4. 易于使用：提供便捷的辅助函数
//
// 响应格式：
// - 成功响应：包含 code、message、data 字段
// - 错误响应：包含 code、message、error 字段
//
// 使用示例：
//
//	// 成功响应
//	handlers.Success(w, http.StatusOK, user)
//
//	// 错误响应
//	handlers.Error(w, http.StatusBadRequest, err)
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/golang/pkg/errors"
)

// APIResponse 是统一的 API 响应格式。
//
// 功能说明：
// - 所有 API 响应都使用此格式
// - 包含状态码、消息、数据和错误信息
//
// 字段说明：
// - Code: HTTP 状态码（如 200、400、500）
// - Message: 响应消息（"success" 或 "error"）
// - Data: 响应数据（成功时包含，使用 omitempty 表示可选）
// - Error: 错误信息（错误时包含，使用 omitempty 表示可选）
//
// 响应示例：
//
//	// 成功响应
//	{
//	  "code": 200,
//	  "message": "success",
//	  "data": {
//	    "id": "123",
//	    "email": "user@example.com"
//	  }
//	}
//
//	// 错误响应
//	{
//	  "code": 400,
//	  "message": "error",
//	  "error": {
//	    "code": "INVALID_EMAIL",
//	    "message": "Invalid email format"
//	  }
//	}
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError 是 API 错误信息的结构。
//
// 功能说明：
// - 定义错误响应的格式
// - 包含错误代码和错误消息
//
// 字段说明：
// - Code: 错误代码（如 "INVALID_EMAIL"、"NOT_FOUND"）
// - Message: 错误消息（人类可读的错误描述）
//
// 错误代码约定：
// - 使用大写字母和下划线（如 "INVALID_EMAIL"）
// - 错误代码应该是唯一的、有意义的
// - 错误消息应该清晰、可操作
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success 发送成功响应。
//
// 功能说明：
// - 构造成功响应并写入 HTTP 响应
// - 设置适当的 HTTP 状态码
// - 将数据序列化为 JSON
//
// 参数：
// - w: HTTP 响应写入器
// - code: HTTP 状态码（如 http.StatusOK、http.StatusCreated）
// - data: 响应数据（可以是任意类型，会自动序列化为 JSON）
//
// 使用示例：
//
//	// 返回用户信息
//	handlers.Success(w, http.StatusOK, user)
//
//	// 返回用户列表
//	handlers.Success(w, http.StatusOK, users)
//
//	// 返回创建的资源
//	handlers.Success(w, http.StatusCreated, createdUser)
//
// 注意事项：
// - 状态码应使用标准 HTTP 状态码
// - 数据会被序列化为 JSON，确保数据可序列化
// - 响应会自动设置 Content-Type 为 application/json
func Success(w http.ResponseWriter, code int, data interface{}) {
	response := APIResponse{
		Code:    code,
		Message: "success",
		Data:    data,
	}
	writeJSON(w, code, response)
}

// Error 发送错误响应。
//
// 功能说明：
// - 构造错误响应并写入 HTTP 响应
// - 处理应用层错误（AppError）和普通错误
// - 设置适当的 HTTP 状态码
//
// 参数：
// - w: HTTP 响应写入器
// - code: HTTP 状态码（如 http.StatusBadRequest、http.StatusInternalServerError）
// - err: 错误对象
//   如果是 AppError，会提取错误代码和消息
//   如果是普通错误，会使用默认错误代码
//
// 使用示例：
//
//	// 处理应用层错误
//	if err != nil {
//	    handlers.Error(w, http.StatusBadRequest, err)
//	    return
//	}
//
//	// 处理验证错误
//	if err := validateRequest(req); err != nil {
//	    handlers.Error(w, http.StatusBadRequest, err)
//	    return
//	}
//
// 错误处理：
// - AppError: 提取错误代码和消息
// - 普通错误: 使用 "INTERNAL_ERROR" 作为错误代码
//
// 注意事项：
// - 状态码应使用标准 HTTP 状态码
// - 错误信息会被序列化为 JSON
// - 不应在生产环境暴露敏感错误信息
func Error(w http.ResponseWriter, code int, err error) {
	response := APIResponse{
		Code:    code,
		Message: "error",
	}

	// 检查是否是 AppError
	// AppError 是应用层定义的错误类型，包含错误代码和消息
	if appErr, ok := err.(*errors.AppError); ok {
		response.Error = &APIError{
			Code:    string(appErr.Code),
			Message: appErr.Message,
		}
	} else {
		// 普通错误，使用默认错误代码
		response.Error = &APIError{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}
	}

	writeJSON(w, code, response)
}

// writeJSON 将数据写入 HTTP 响应为 JSON 格式。
//
// 功能说明：
// - 设置 Content-Type 为 application/json
// - 设置 HTTP 状态码
// - 将数据序列化为 JSON 并写入响应
//
// 参数：
// - w: HTTP 响应写入器
// - code: HTTP 状态码
// - data: 要序列化的数据
//
// 注意事项：
// - 数据必须是可序列化为 JSON 的类型
// - 如果序列化失败，会返回错误（但不会处理）
// - 应在调用此函数前确保数据有效
func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
