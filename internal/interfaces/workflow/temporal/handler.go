// Package temporal 提供 Temporal 工作流的 HTTP 接口处理器
//
// 设计原理：
// 1. 这是 Interfaces Layer 的工作流 HTTP 接口实现
// 2. 将 HTTP 请求转换为工作流操作
// 3. 调用 Temporal 客户端启动和查询工作流
//
// 架构位置：
// - 位置：Interfaces Layer (internal/interfaces/workflow/temporal/)
// - 职责：工作流 HTTP 协议适配
// - 依赖：Temporal 客户端
package temporal

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/golang/internal/infra/workflow/temporal"
)

// Handler Temporal 工作流 HTTP 处理器
type Handler struct {
	client *temporal.Client
}

// NewHandler 创建 Temporal 工作流处理器
func NewHandler(client *temporal.Client) *Handler {
	return &Handler{client: client}
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
func (h *Handler) StartUserWorkflow(w http.ResponseWriter, r *http.Request) {
	var req StartUserWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 启动工作流
	workflowID, runID, err := h.client.StartUserWorkflow(r.Context(), req.Email, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := StartUserWorkflowResponse{
		WorkflowID: workflowID,
		RunID:      runID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetWorkflowResultResponse 获取工作流结果响应
type GetWorkflowResultResponse struct {
	Result string `json:"result"`
	Status string `json:"status"`
}

// GetWorkflowResult 获取工作流结果
// GET /api/v1/workflows/user/{workflow_id}/result
func (h *Handler) GetWorkflowResult(w http.ResponseWriter, r *http.Request) {
	workflowID := r.URL.Path[len("/api/v1/workflows/user/"):]
	// 去掉 "/result" 后缀
	if len(workflowID) > 7 && workflowID[len(workflowID)-7:] == "/result" {
		workflowID = workflowID[:len(workflowID)-7]
	}

	result, err := h.client.GetWorkflowResult(r.Context(), workflowID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GetWorkflowResultResponse{
		Result: result,
		Status: "completed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

