package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	appworkflow "github.com/yourusername/golang/internal/application/workflow"
	"github.com/yourusername/golang/pkg/errors"
	temporalhandler "github.com/yourusername/golang/internal/interfaces/workflow/temporal"
)

// WorkflowHandler 工作流处理器
type WorkflowHandler struct {
	workflowHandler *temporalhandler.Handler
}

// NewWorkflowHandler 创建工作流处理器
func NewWorkflowHandler(workflowHandler *temporalhandler.Handler) *WorkflowHandler {
	return &WorkflowHandler{workflowHandler: workflowHandler}
}

// StartUserWorkflow 启动用户工作流
func (h *WorkflowHandler) StartUserWorkflow(w http.ResponseWriter, r *http.Request) {
	var input appworkflow.UserWorkflowInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Invalid request body"))
		return
	}

	workflowRun, err := h.workflowHandler.StartUserWorkflow(r.Context(), input)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	Success(w, http.StatusAccepted, map[string]interface{}{
		"workflow_id": workflowRun.GetID(),
		"run_id":      workflowRun.GetRunID(),
	})
}

// GetWorkflowResult 获取工作流结果
func (h *WorkflowHandler) GetWorkflowResult(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "workflow_id")
	runID := r.URL.Query().Get("run_id")

	if workflowID == "" {
		Error(w, http.StatusBadRequest, errors.NewInvalidInputError("Workflow ID is required"))
		return
	}

	result, err := h.workflowHandler.GetWorkflowResult(r.Context(), workflowID, runID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	Success(w, http.StatusOK, result)
}
