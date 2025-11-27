package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/golang/pkg/errors"
)

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	Success(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != http.StatusOK {
		t.Errorf("Expected code %d, got %d", http.StatusOK, resp.Code)
	}

	if resp.Message != "success" {
		t.Errorf("Expected message 'success', got '%s'", resp.Message)
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.NewNotFoundError("user", "123")

	Error(w, http.StatusNotFound, err)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Error("Expected error info, got nil")
		return
	}

	if resp.Error.Code != "NOT_FOUND" {
		t.Errorf("Expected error code 'NOT_FOUND', got '%s'", resp.Error.Code)
	}
}

func TestPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	data := []string{"item1", "item2", "item3"}

	Paginated(w, http.StatusOK, data, 1, 10, 25)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp PaginatedResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Pagination == nil {
		t.Error("Expected pagination info, got nil")
		return
	}

	if resp.Pagination.Page != 1 {
		t.Errorf("Expected page 1, got %d", resp.Pagination.Page)
	}

	if resp.Pagination.Total != 25 {
		t.Errorf("Expected total 25, got %d", resp.Pagination.Total)
	}

	if resp.Pagination.TotalPages != 3 {
		t.Errorf("Expected total pages 3, got %d", resp.Pagination.TotalPages)
	}
}

func TestSuccessWithTraceID(t *testing.T) {
	w := httptest.NewRecorder()
	traceID := "trace-123"

	SuccessWithTraceID(w, http.StatusOK, nil, traceID)

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.TraceID != traceID {
		t.Errorf("Expected trace ID '%s', got '%s'", traceID, resp.TraceID)
	}
}

func TestErrorWithTraceID(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.NewNotFoundError("user", "123")
	traceID := "trace-456"

	ErrorWithTraceID(w, http.StatusNotFound, err, traceID)

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.TraceID != traceID {
		t.Errorf("Expected trace ID '%s', got '%s'", traceID, resp.TraceID)
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := Response{
		Code:      http.StatusOK,
		Message:   "test",
		Timestamp: time.Now(),
	}

	writeJSON(w, http.StatusOK, data)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Message != "test" {
		t.Errorf("Expected message 'test', got '%s'", resp.Message)
	}
}

func TestMeta_WithExtra(t *testing.T) {
	meta := NewMeta("req-123", "v1.0.0")
	meta.WithExtra("key", "value")

	if meta.Extra["key"] != "value" {
		t.Errorf("Expected extra['key'] = 'value', got '%v'", meta.Extra["key"])
	}
}
