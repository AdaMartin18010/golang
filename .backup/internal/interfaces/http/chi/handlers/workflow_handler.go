// Package handlers provides HTTP handlers for workflow-related operations.
//
// 工作流 HTTP 处理器负责：
// 1. 接收 HTTP 请求
// 2. 解析请求参数
// 3. 调用 Temporal 客户端启动和查询工作流
// 4. 格式化 HTTP 响应
//
// 设计原则：
// 1. 协议适配：将 HTTP 协议转换为工作流操作
// 2. 参数验证：验证 HTTP 请求参数
// 3. 错误处理：将工作流错误映射为 HTTP 状态码
// 4. 响应格式化：统一 API 响应格式
//
// 架构位置：
// - 位置：Interfaces Layer (internal/interfaces/http/chi/handlers/)
// - 职责：工作流 HTTP 协议适配、请求处理、响应格式化
// - 依赖：Temporal 客户端
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	temporalhandler "github.com/yourusername/golang/internal/interfaces/workflow/temporal"
)

// WorkflowHandler 工作流 HTTP 处理器
type WorkflowHandler struct {
	temporalClient *temporalhandler.Handler
}

// NewWorkflowHandler 创建工作流 HTTP 处理器
func NewWorkflowHandler(temporalClient *temporalhandler.Handler) *WorkflowHandler {
	return &WorkflowHandler{
		temporalClient: temporalClient,
	}
}

// StartUserWorkflowRequest 启动用户工作流请求
type StartUserWorkflowRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// StartUserWorkflowResponse 启动用户工作流响应
type StartUserWorkflowResponse struct {
	WorkflowID string `json:"workflow_id"`
	RunID      string `json:"run_id"`
}

// StartUserWorkflow 启动用户工作流
// POST /api/v1/workflows/user
func (h *WorkflowHandler) StartUserWorkflow(w http.ResponseWriter, r *http.Request) {
	var req StartUserWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, NewInvalidInputError("invalid request body"))
		return
	}

	if req.Email == "" {
		Error(w, http.StatusBadRequest, NewInvalidInputError("email is required"))
		return
	}
	if req.Name == "" {
		Error(w, http.StatusBadRequest, NewInvalidInputError("name is required"))
		return
	}

	// 调用 temporal handler 的方法
	h.temporalClient.StartUserWorkflow(w, r)
}

// GetWorkflowResultResponse 获取工作流结果响应
type GetWorkflowResultResponse struct {
	Result string `json:"result"`
	Status string `json:"status"`
}

// GetWorkflowResult 获取工作流结果
// GET /api/v1/workflows/user/{workflow_id}/result
func (h *WorkflowHandler) GetWorkflowResult(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "workflow_id")
	if workflowID == "" {
		Error(w, http.StatusBadRequest, NewInvalidInputError("workflow_id is required"))
		return
	}

	// 调用 temporal handler 的方法
	h.temporalClient.GetWorkflowResult(w, r)
}

// NewInvalidInputError 创建无效输入错误（简化实现）
// 实际应该从 apperrors 包导入
func NewInvalidInputError(message string) error {
	return &simpleError{message: message, status: http.StatusBadRequest}
}

type simpleError struct {
	message string
	status  int
}

func (e *simpleError) Error() string {
	return e.message
}
