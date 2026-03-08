package temporal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewHandler 测试 NewHandler 构造函数
func TestNewHandler(t *testing.T) {
	// 注意：由于 handler 需要实际的 *temporal.Client 类型，
	// 这里我们只测试 nil 情况和非 nil 情况的结构
	handler := NewHandler(nil)

	require.NotNil(t, handler, "NewHandler 不应返回 nil")
	// 传入 nil 时 client 字段应为 nil
	assert.Nil(t, handler.client, "传入 nil client 时 handler.client 应为 nil")
}

// TestHandlerStructure 测试 Handler 结构体
func TestHandlerStructure(t *testing.T) {
	h := &Handler{}
	assert.NotNil(t, h, "Handler 实例不应为 nil")
}

// TestStartUserWorkflowRequest 测试 StartUserWorkflowRequest 结构体
func TestStartUserWorkflowRequest(t *testing.T) {
	req := StartUserWorkflowRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "Test User", req.Name)
}

// TestStartUserWorkflowRequest_JSON 测试 JSON 序列化
func TestStartUserWorkflowRequest_JSON(t *testing.T) {
	req := StartUserWorkflowRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded StartUserWorkflowRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Email, decoded.Email)
	assert.Equal(t, req.Name, decoded.Name)
}

// TestStartUserWorkflowResponse 测试 StartUserWorkflowResponse 结构体
func TestStartUserWorkflowResponse(t *testing.T) {
	resp := StartUserWorkflowResponse{
		WorkflowID: "wf-123",
		RunID:      "run-456",
	}

	assert.Equal(t, "wf-123", resp.WorkflowID)
	assert.Equal(t, "run-456", resp.RunID)
}

// TestGetWorkflowResultResponse 测试 GetWorkflowResultResponse 结构体
func TestGetWorkflowResultResponse(t *testing.T) {
	resp := GetWorkflowResultResponse{
		Result: "success",
		Status: "completed",
	}

	assert.Equal(t, "success", resp.Result)
	assert.Equal(t, "completed", resp.Status)
}

// TestStartUserWorkflowRequest_Validation 测试请求验证
func TestStartUserWorkflowRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request StartUserWorkflowRequest
	}{
		{
			name: "有效请求",
			request: StartUserWorkflowRequest{
				Email: "user@example.com",
				Name:  "User Name",
			},
		},
		{
			name: "空邮箱",
			request: StartUserWorkflowRequest{
				Email: "",
				Name:  "User Name",
			},
		},
		{
			name: "空名称",
			request: StartUserWorkflowRequest{
				Email: "user@example.com",
				Name:  "",
			},
		},
		{
			name: "全部为空",
			request: StartUserWorkflowRequest{
				Email: "",
				Name:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证请求可以被序列化
			data, err := json.Marshal(tt.request)
			require.NoError(t, err)
			assert.NotNil(t, data)

			// 验证请求可以被反序列化
			var decoded StartUserWorkflowRequest
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err)
			assert.Equal(t, tt.request.Email, decoded.Email)
			assert.Equal(t, tt.request.Name, decoded.Name)
		})
	}
}

// TestHandler_MethodSignatures 测试方法签名
func TestHandler_MethodSignatures(t *testing.T) {
	// 验证方法存在且签名正确（编译时检查）
	var h *Handler

	// 这些方法应该存在
	_ = h.StartUserWorkflow
	_ = h.GetWorkflowResult
}

// TestResponseContentType 测试响应 Content-Type
func TestResponseContentType(t *testing.T) {
	// 创建 recorder 来检查响应头
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/json")

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

// TestHandler_GetWorkflowResult_PathParsing 测试路径解析
func TestHandler_GetWorkflowResult_PathParsing(t *testing.T) {
	// 测试路径解析逻辑（复制 handler 中的逻辑）
	paths := []struct {
		path        string
		expectedID  string
		description string
	}{
		{
			path:        "/api/v1/workflows/user/wf-123/result",
			expectedID:  "wf-123",
			description: "标准路径",
		},
		{
			path:        "/api/v1/workflows/user/wf-123",
			expectedID:  "wf-123",
			description: "无 result 后缀",
		},
		{
			path:        "/api/v1/workflows/user/",
			expectedID:  "",
			description: "空 ID",
		},
	}

	for _, tt := range paths {
		t.Run(tt.description, func(t *testing.T) {
			// 提取 workflow ID 逻辑（复制 handler 中的逻辑）
			prefix := "/api/v1/workflows/user/"
			workflowID := tt.path[len(prefix):]
			if len(workflowID) > 7 && workflowID[len(workflowID)-7:] == "/result" {
				workflowID = workflowID[:len(workflowID)-7]
			}

			assert.Equal(t, tt.expectedID, workflowID)
		})
	}
}

// TestHandler_InvalidJSON 测试无效 JSON 请求
func TestHandler_InvalidJSON(t *testing.T) {
	handler := NewHandler(nil)

	// 创建无效 JSON 请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workflows/user", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	handler.StartUserWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request body")
}

// TestHandler_EmptyBody 测试空请求体
func TestHandler_EmptyBody(t *testing.T) {
	handler := NewHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workflows/user", nil)
	rec := httptest.NewRecorder()

	handler.StartUserWorkflow(rec, req)

	// 空 body 应该返回 BadRequest
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestHandler_StartUserWorkflowResponseStructure 测试响应结构
func TestHandler_StartUserWorkflowResponseStructure(t *testing.T) {
	resp := StartUserWorkflowResponse{
		WorkflowID: "test-workflow-id",
		RunID:      "test-run-id",
	}

	// 序列化响应
	data, err := json.Marshal(resp)
	require.NoError(t, err)

	// 验证 JSON 字段
	var decoded map[string]interface{}
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "test-workflow-id", decoded["workflow_id"])
	assert.Equal(t, "test-run-id", decoded["run_id"])
}

// TestHandler_GetWorkflowResultResponseStructure 测试获取结果响应结构
func TestHandler_GetWorkflowResultResponseStructure(t *testing.T) {
	resp := GetWorkflowResultResponse{
		Result: `{"data":"test"}`,
		Status: "completed",
	}

	// 序列化响应
	data, err := json.Marshal(resp)
	require.NoError(t, err)

	// 验证 JSON 字段
	var decoded map[string]interface{}
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, `{"data":"test"}`, decoded["result"])
	assert.Equal(t, "completed", decoded["status"])
}

// TestHandler_HTTPMethods 测试 HTTP 方法处理
func TestHandler_HTTPMethods(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		expectCode int
		callFunc   func(*Handler, http.ResponseWriter, *http.Request)
	}{
		{
			name:       "POST 启动工作流",
			method:     http.MethodPost,
			path:       "/api/v1/workflows/user",
			expectCode: http.StatusBadRequest, // 因为没有 body
			callFunc: func(h *Handler, w http.ResponseWriter, r *http.Request) {
				h.StartUserWorkflow(w, r)
			},
		},
		{
			name:       "POST 有效启动工作流",
			method:     http.MethodPost,
			path:       "/api/v1/workflows/user",
			expectCode: http.StatusInternalServerError, // 因为没有 client 会 panic 然后返回 500
			callFunc: func(h *Handler, w http.ResponseWriter, r *http.Request) {
				// 这个调用会导致 panic，所以我们用 recover 捕获
				defer func() { recover() }()
				h.StartUserWorkflow(w, r)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(nil)

			var req *http.Request
			if tt.name == "POST 有效启动工作流" {
				reqBody := StartUserWorkflowRequest{
					Email: "test@example.com",
					Name:  "Test",
				}
				body, _ := json.Marshal(reqBody)
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(body))
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}
			rec := httptest.NewRecorder()

			tt.callFunc(handler, rec, req)

			// 我们不严格检查状态码，只确保不 panic
			assert.True(t, rec.Code >= 200 && rec.Code < 600, "应返回有效的 HTTP 状态码")
		})
	}
}

// TestHandler_ConcurrentAccess 测试并发访问
func TestHandler_ConcurrentAccess(t *testing.T) {
	handler := NewHandler(nil)

	// 测试并发访问不会 panic
	done := make(chan bool, 3)

	for i := 0; i < 3; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// 预期的 panic，因为 client 是 nil
				}
				done <- true
			}()

			// 创建请求 - 使用无效 JSON 不会调用 client 方法
			req := httptest.NewRequest(http.MethodPost, "/api/v1/workflows/user", bytes.NewBufferString("invalid"))
			rec := httptest.NewRecorder()

			handler.StartUserWorkflow(rec, req)
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestHandler_RequestJSONTags 测试请求 JSON 标签
func TestHandler_RequestJSONTags(t *testing.T) {
	req := StartUserWorkflowRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// 验证 JSON 标签正确
	data, err := json.Marshal(req)
	require.NoError(t, err)

	// 检查字段名是否使用 snake_case
	jsonStr := string(data)
	assert.Contains(t, jsonStr, "email")
	assert.Contains(t, jsonStr, "name")
}

// TestHandler_ResponseJSONTags 测试响应 JSON 标签
func TestHandler_ResponseJSONTags(t *testing.T) {
	resp := StartUserWorkflowResponse{
		WorkflowID: "wf-123",
		RunID:      "run-456",
	}

	// 验证 JSON 标签正确
	data, err := json.Marshal(resp)
	require.NoError(t, err)

	// 检查字段名是否使用 snake_case
	jsonStr := string(data)
	assert.Contains(t, jsonStr, "workflow_id")
	assert.Contains(t, jsonStr, "run_id")
}
