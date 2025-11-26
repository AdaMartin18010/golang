package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yourusername/golang/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	TraceID   string      `json:"trace_id,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Meta 元数据
type Meta struct {
	RequestID string                 `json:"request_id,omitempty"`
	Version   string                 `json:"version,omitempty"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

// PaginationMeta 分页元数据
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse 分页响应
type PaginatedResponse struct {
	Response
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

// Success 成功响应
func Success(w http.ResponseWriter, code int, data interface{}) {
	writeJSON(w, code, Response{
		Code:      code,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now(),
	})
}

// SuccessWithMeta 带元数据的成功响应
func SuccessWithMeta(w http.ResponseWriter, code int, data interface{}, meta *Meta) {
	writeJSON(w, code, Response{
		Code:      code,
		Message:   "success",
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
	})
}

// SuccessWithTraceID 带追踪ID的成功响应
func SuccessWithTraceID(w http.ResponseWriter, code int, data interface{}, traceID string) {
	writeJSON(w, code, Response{
		Code:      code,
		Message:   "success",
		Data:      data,
		TraceID:   traceID,
		Timestamp: time.Now(),
	})
}

// Paginated 分页响应
func Paginated(w http.ResponseWriter, code int, data interface{}, page, pageSize int, total int64) {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	writeJSON(w, code, PaginatedResponse{
		Response: Response{
			Code:      code,
			Message:   "success",
			Data:      data,
			Timestamp: time.Now(),
		},
		Pagination: &PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// Error 错误响应
func Error(w http.ResponseWriter, code int, err error) {
	response := Response{
		Code:      code,
		Message:   "error",
		Timestamp: time.Now(),
	}

	// 检查是否是 AppError
	if appErr, ok := err.(*errors.AppError); ok {
		response.Error = &ErrorInfo{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		}
		response.TraceID = appErr.TraceID
		// 使用 AppError 的 HTTP 状态码
		code = appErr.HTTPStatusCode()
	} else {
		response.Error = &ErrorInfo{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}
	}

	writeJSON(w, code, response)
}

// ErrorWithTraceID 带追踪ID的错误响应
func ErrorWithTraceID(w http.ResponseWriter, code int, err error, traceID string) {
	response := Response{
		Code:      code,
		Message:   "error",
		TraceID:   traceID,
		Timestamp: time.Now(),
	}

	// 检查是否是 AppError
	if appErr, ok := err.(*errors.AppError); ok {
		response.Error = &ErrorInfo{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		}
		if appErr.TraceID != "" {
			response.TraceID = appErr.TraceID
		}
		code = appErr.HTTPStatusCode()
	} else {
		response.Error = &ErrorInfo{
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

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// 如果编码失败，记录错误但不返回错误响应（避免循环）
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NewMeta 创建元数据
func NewMeta(requestID, version string) *Meta {
	return &Meta{
		RequestID: requestID,
		Version:   version,
		Extra:     make(map[string]interface{}),
	}
}

// WithExtra 添加额外信息
func (m *Meta) WithExtra(key string, value interface{}) *Meta {
	if m.Extra == nil {
		m.Extra = make(map[string]interface{})
	}
	m.Extra[key] = value
	return m
}
