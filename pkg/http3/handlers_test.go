package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestHandleRoot 测试根路径处理器
func TestHandleRoot(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handleRoot(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// 验证状态码
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// 验证Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// 解析响应
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// 验证响应内容
	if response.Message != "Welcome to HTTP/3 Server!" {
		t.Errorf("Unexpected message: %s", response.Message)
	}

	if response.Protocol == "" {
		t.Error("Protocol should not be empty")
	}

	if response.Server != "Go 1.23+" {
		t.Errorf("Expected server 'Go 1.23+', got '%s'", response.Server)
	}
}

// TestHandleRootHTTP2 测试HTTP/2协议检测
func TestHandleRootHTTP2(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.ProtoMajor = 2
	req.ProtoMinor = 0
	w := httptest.NewRecorder()

	handleRoot(w, req)

	var response Response
	json.NewDecoder(w.Result().Body).Decode(&response)

	if response.Protocol != "HTTP/2" {
		t.Errorf("Expected HTTP/2, got %s", response.Protocol)
	}
}

// TestHandleStats 测试统计信息处理器
func TestHandleStats(t *testing.T) {
	// 重置统计
	stats = Stats{StartTime: time.Now()}

	// 发送几个请求以生成统计数据
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handleRoot(w, req)
	}

	// 测试stats端点
	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()

	handleStats(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var statsData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&statsData); err != nil {
		t.Fatalf("Failed to decode stats: %v", err)
	}

	// 验证统计数据
	if requests, ok := statsData["requests"].(float64); !ok || requests < 5 {
		t.Errorf("Expected at least 5 requests, got %v", statsData["requests"])
	}

	if _, ok := statsData["uptime"]; !ok {
		t.Error("Stats should include uptime")
	}

	if _, ok := statsData["avg_duration"]; !ok {
		t.Error("Stats should include avg_duration")
	}

	if _, ok := statsData["req_per_sec"]; !ok {
		t.Error("Stats should include req_per_sec")
	}
}

// TestHandleHealth 测试健康检查处理器
func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handleHealth(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var healthData map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&healthData); err != nil {
		t.Fatalf("Failed to decode health response: %v", err)
	}

	if status, ok := healthData["status"]; !ok || status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", status)
	}

	if timeStr, ok := healthData["time"]; !ok || timeStr == "" {
		t.Error("Health response should include timestamp")
	}
}

// TestHandleData 测试数据处理器
func TestHandleData(t *testing.T) {
	req := httptest.NewRequest("GET", "/data", nil)
	w := httptest.NewRecorder()

	handleData(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode data response: %v", err)
	}

	// 验证数据数量
	if len(data) != 100 {
		t.Errorf("Expected 100 items, got %d", len(data))
	}

	// 验证第一个项目
	if len(data) > 0 {
		item := data[0]
		if id, ok := item["id"].(float64); !ok || id != 0 {
			t.Errorf("Expected first item id=0, got %v", item["id"])
		}
		if name, ok := item["name"].(string); !ok || name != "Item-0" {
			t.Errorf("Expected first item name='Item-0', got %v", item["name"])
		}
	}
}

// TestMultipleRequests 测试多个并发请求
func TestMultipleRequests(t *testing.T) {
	stats = Stats{StartTime: time.Now()}

	const numRequests = 20
	results := make(chan int, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			handleRoot(w, req)
			results <- w.Result().StatusCode
		}()
	}

	// 收集结果
	for i := 0; i < numRequests; i++ {
		statusCode := <-results
		if statusCode != http.StatusOK {
			t.Errorf("Request %d failed with status %d", i, statusCode)
		}
	}

	// 验证统计计数
	if stats.Requests != int64(numRequests) {
		t.Errorf("Expected %d requests, got %d", numRequests, stats.Requests)
	}
}

// TestResponseStructure 测试Response结构
func TestResponseStructure(t *testing.T) {
	resp := Response{
		Message:   "Test",
		Timestamp: time.Now(),
		Protocol:  "HTTP/2",
		Server:    "Test Server",
	}

	// 验证JSON编码
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal Response: %v", err)
	}

	// 验证JSON解码
	var decoded Response
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Response: %v", err)
	}

	if decoded.Message != resp.Message {
		t.Errorf("Message mismatch after JSON round-trip")
	}

	if decoded.Protocol != resp.Protocol {
		t.Errorf("Protocol mismatch after JSON round-trip")
	}
}

// TestStatsStructure 测试Stats结构
func TestStatsStructure(t *testing.T) {
	testStats := Stats{
		Requests:    100,
		StartTime:   time.Now(),
		AvgDuration: 5 * time.Millisecond,
	}

	if testStats.Requests != 100 {
		t.Errorf("Expected 100 requests, got %d", testStats.Requests)
	}

	if testStats.AvgDuration != 5*time.Millisecond {
		t.Errorf("Expected 5ms avg duration, got %v", testStats.AvgDuration)
	}
}

// BenchmarkHandleRoot 基准测试根处理器
func BenchmarkHandleRoot(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleRoot(w, req)
	}
}

// BenchmarkHandleStats 基准测试统计处理器
func BenchmarkHandleStats(b *testing.B) {
	req := httptest.NewRequest("GET", "/stats", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleStats(w, req)
	}
}

// BenchmarkHandleHealth 基准测试健康检查
func BenchmarkHandleHealth(b *testing.B) {
	req := httptest.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleHealth(w, req)
	}
}

// BenchmarkHandleData 基准测试数据处理器
func BenchmarkHandleData(b *testing.B) {
	req := httptest.NewRequest("GET", "/data", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleData(w, req)
	}
}
