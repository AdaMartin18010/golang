package temporal

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
	appworkflow "github.com/yourusername/golang/internal/application/workflow"
)

// Handler Temporal 工作流处理器
type Handler struct {
	client client.Client
}

// NewHandler 创建处理器
func NewHandler(client client.Client) *Handler {
	return &Handler{client: client}
}

// StartUserWorkflow 启动用户工作流
func (h *Handler) StartUserWorkflow(ctx context.Context, input appworkflow.UserWorkflowInput) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("user-workflow-%s-%s", input.Action, input.UserID),
		TaskQueue: "user-task-queue",
	}

	workflowRun, err := h.client.ExecuteWorkflow(ctx, options, appworkflow.UserWorkflow, input)
	if err != nil {
		return nil, fmt.Errorf("failed to start workflow: %w", err)
	}

	return workflowRun, nil
}

// GetWorkflowResult 获取工作流结果
func (h *Handler) GetWorkflowResult(ctx context.Context, workflowID, runID string) (appworkflow.UserWorkflowOutput, error) {
	var result appworkflow.UserWorkflowOutput

	workflowRun := h.client.GetWorkflow(ctx, workflowID, runID)
	err := workflowRun.Get(ctx, &result)
	if err != nil {
		return result, fmt.Errorf("failed to get workflow result: %w", err)
	}

	return result, nil
}

// SignalWorkflow 发送信号给工作流
func (h *Handler) SignalWorkflow(ctx context.Context, workflowID, runID, signalName string, arg interface{}) error {
	return h.client.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
}

// QueryWorkflow 查询工作流
func (h *Handler) QueryWorkflow(ctx context.Context, workflowID, runID, queryType string, args ...interface{}) (interface{}, error) {
	return h.client.QueryWorkflow(ctx, workflowID, runID, queryType, args...)
}
