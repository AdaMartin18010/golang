package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	resp, err := client.Get(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if !resp.IsSuccess() {
		t.Error("Expected success response")
	}
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":1}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	body := map[string]interface{}{
		"name": "test",
	}

	resp, err := client.Post(context.Background(), "/test", body, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestClient_SetHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer token123" {
			t.Errorf("Expected Authorization header, got %s", auth)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})
	client.SetHeader("Authorization", "Bearer token123")

	resp, err := client.Get(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestResponse_JSON(t *testing.T) {
	resp := &Response{
		StatusCode: http.StatusOK,
		Body:       []byte(`{"name":"test","age":30}`),
	}

	var result map[string]interface{}
	err := resp.JSON(&result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected 'test', got %v", result["name"])
	}
}

func TestResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, true},
		{201, true},
		{299, true},
		{400, false},
		{500, false},
	}

	for _, tt := range tests {
		resp := &Response{StatusCode: tt.statusCode}
		if resp.IsSuccess() != tt.expected {
			t.Errorf("Status %d: expected %v, got %v", tt.statusCode, tt.expected, resp.IsSuccess())
		}
	}
}

func TestResponse_IsError(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, false},
		{300, false},
		{400, true},
		{500, true},
	}

	for _, tt := range tests {
		resp := &Response{StatusCode: tt.statusCode}
		if resp.IsError() != tt.expected {
			t.Errorf("Status %d: expected %v, got %v", tt.statusCode, tt.expected, resp.IsError())
		}
	}
}
