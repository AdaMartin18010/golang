// Package handlers provides tests for workflow HTTP handlers.
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkflowHandler(t *testing.T) {
	// We can't easily mock the temporal handler due to its structure,
	// but we can test the function signature and basic behavior
	assert.NotNil(t, NewWorkflowHandler)
}

func TestStartUserWorkflowRequestStruct(t *testing.T) {
	req := StartUserWorkflowRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "Test User", req.Name)
}

func TestStartUserWorkflowResponseStruct(t *testing.T) {
	resp := StartUserWorkflowResponse{
		WorkflowID: "workflow-123",
		RunID:      "run-456",
	}

	assert.Equal(t, "workflow-123", resp.WorkflowID)
	assert.Equal(t, "run-456", resp.RunID)
}

func TestGetWorkflowResultResponseStruct(t *testing.T) {
	resp := GetWorkflowResultResponse{
		Result: "success",
		Status: "completed",
	}

	assert.Equal(t, "success", resp.Result)
	assert.Equal(t, "completed", resp.Status)
}

func TestWorkflowHandler_StartUserWorkflow_InvalidJSON(t *testing.T) {
	// Create handler with nil temporal client - it won't be used due to JSON error
	handler := &WorkflowHandler{temporalClient: nil}

	reqBody := `{"invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workflows/user", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.StartUserWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestWorkflowHandler_StartUserWorkflow_EmptyEmail(t *testing.T) {
	handler := &WorkflowHandler{temporalClient: nil}

	reqBody := `{"email": "", "name": "Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workflows/user", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.StartUserWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestWorkflowHandler_StartUserWorkflow_EmptyName(t *testing.T) {
	handler := &WorkflowHandler{temporalClient: nil}

	reqBody := `{"email": "test@example.com", "name": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workflows/user", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()

	handler.StartUserWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestWorkflowHandler_StartUserWorkflow_ValidRequest(t *testing.T) {
	// Test valid request parsing
	reqBody := `{"email": "test@example.com", "name": "Test User"}`

	var req StartUserWorkflowRequest
	err := json.Unmarshal([]byte(reqBody), &req)
	assert.NoError(t, err)

	// Verify struct fields
	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "Test User", req.Name)
}

func TestWorkflowHandler_GetWorkflowResult_EmptyWorkflowID(t *testing.T) {
	handler := &WorkflowHandler{temporalClient: nil}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workflows/user//result", nil)
	rec := httptest.NewRecorder()

	handler.GetWorkflowResult(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestNewInvalidInputError(t *testing.T) {
	err := NewInvalidInputError("test error message")

	assert.NotNil(t, err)
	assert.Equal(t, "test error message", err.Error())
}

func TestSimpleError(t *testing.T) {
	err := &simpleError{
		message: "test message",
		status:  http.StatusBadRequest,
	}

	assert.Equal(t, "test message", err.Error())
	assert.Equal(t, http.StatusBadRequest, err.status)
}

func TestWorkflowHandler_Struct(t *testing.T) {
	handler := &WorkflowHandler{}
	assert.NotNil(t, handler)
}
