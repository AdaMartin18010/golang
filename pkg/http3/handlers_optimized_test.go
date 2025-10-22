package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// =============================================================================
// 优化处理器的测试
// =============================================================================

// TestHandleRootOptimized 测试优化的根处理器
func TestHandleRootOptimized(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handleRootOptimized(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}

// TestHandleStatsOptimized 测试优化的统计处理器
func TestHandleStatsOptimized(t *testing.T) {
	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()

	handleStatsOptimized(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestHandleHealthOptimized 测试优化的健康检查处理器
func TestHandleHealthOptimized(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handleHealthOptimized(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestHandleDataOptimized 测试优化的数据处理器
func TestHandleDataOptimized(t *testing.T) {
	req := httptest.NewRequest("GET", "/data", nil)
	w := httptest.NewRecorder()

	handleDataOptimized(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestHandleDataOptimizedV2 测试进一步优化的数据处理器
func TestHandleDataOptimizedV2(t *testing.T) {
	req := httptest.NewRequest("GET", "/data", nil)
	w := httptest.NewRecorder()

	handleDataOptimizedV2(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestHandleHealthCached 测试缓存的健康检查
func TestHandleHealthCached(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handleHealthCached(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// =============================================================================
// 对象池测试
// =============================================================================

// TestResponsePool 测试Response对象池
func TestResponsePool(t *testing.T) {
	resp := GetResponse()
	if resp == nil {
		t.Fatal("GetResponse returned nil")
	}

	resp.Message = "test"
	PutResponse(resp)

	// 再次获取，可能是同一个对象
	resp2 := GetResponse()
	if resp2 == nil {
		t.Fatal("Second GetResponse returned nil")
	}

	// 对象应该被重置
	if resp2.Message != "" {
		t.Error("Response was not reset properly")
	}

	PutResponse(resp2)
}

// TestBufferPool 测试Buffer对象池
func TestBufferPool(t *testing.T) {
	buf := GetBuffer()
	if buf == nil {
		t.Fatal("GetBuffer returned nil")
	}

	buf.WriteString("test")
	if buf.Len() == 0 {
		t.Error("Buffer write failed")
	}

	PutBuffer(buf)

	// 再次获取
	buf2 := GetBuffer()
	if buf2 == nil {
		t.Fatal("Second GetBuffer returned nil")
	}

	// Buffer应该被重置
	if buf2.Len() != 0 {
		t.Error("Buffer was not reset properly")
	}

	PutBuffer(buf2)
}

// TestDataItemPool 测试数据项对象池
func TestDataItemPool(t *testing.T) {
	item := GetDataItem()
	if item == nil {
		t.Fatal("GetDataItem returned nil")
	}

	item["key"] = "value"
	PutDataItem(item)

	// 再次获取
	item2 := GetDataItem()
	if item2 == nil {
		t.Fatal("Second GetDataItem returned nil")
	}

	// Map应该被清空
	if len(item2) != 0 {
		t.Error("DataItem was not cleared properly")
	}

	PutDataItem(item2)
}

// TestDataSlicePool 测试数据切片对象池
func TestDataSlicePool(t *testing.T) {
	slice := GetDataSlice()
	if slice == nil {
		t.Fatal("GetDataSlice returned nil")
	}

	slice = append(slice, map[string]interface{}{"test": "value"})
	if len(slice) == 0 {
		t.Error("Slice append failed")
	}

	PutDataSlice(slice)

	// 再次获取
	slice2 := GetDataSlice()
	if slice2 == nil {
		t.Fatal("Second GetDataSlice returned nil")
	}

	// 切片应该被清空但保留容量
	if len(slice2) != 0 {
		t.Error("Slice was not cleared properly")
	}
	if cap(slice2) < 100 {
		t.Error("Slice capacity was not preserved")
	}

	PutDataSlice(slice2)
}

// =============================================================================
// 性能基准测试
// =============================================================================

// BenchmarkHandleRootOriginal 原始处理器基准（对比）
func BenchmarkHandleRootOriginal(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleRoot(w, req)
	}
}

// BenchmarkHandleRootOptimized 优化处理器基准
func BenchmarkHandleRootOptimized(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleRootOptimized(w, req)
	}
}

// BenchmarkHandleStatsOriginal 原始统计处理器基准
func BenchmarkHandleStatsOriginal(b *testing.B) {
	req := httptest.NewRequest("GET", "/stats", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleStats(w, req)
	}
}

// BenchmarkHandleStatsOptimized 优化统计处理器基准
func BenchmarkHandleStatsOptimized(b *testing.B) {
	req := httptest.NewRequest("GET", "/stats", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleStatsOptimized(w, req)
	}
}

// BenchmarkHandleDataOriginal 原始数据处理器基准
func BenchmarkHandleDataOriginal(b *testing.B) {
	req := httptest.NewRequest("GET", "/data", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		// 注意：原始版本有10ms的sleep，所以基准测试会很慢
		// 在基准测试中应该移除sleep
		handleData(w, req)
	}
}

// BenchmarkHandleDataOptimized 优化数据处理器基准
func BenchmarkHandleDataOptimized(b *testing.B) {
	req := httptest.NewRequest("GET", "/data", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleDataOptimized(w, req)
	}
}

// BenchmarkHandleDataOptimizedV2 进一步优化的数据处理器基准
func BenchmarkHandleDataOptimizedV2(b *testing.B) {
	req := httptest.NewRequest("GET", "/data", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleDataOptimizedV2(w, req)
	}
}

// BenchmarkHandleHealthOriginal 原始健康检查基准
func BenchmarkHandleHealthOriginal(b *testing.B) {
	req := httptest.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleHealth(w, req)
	}
}

// BenchmarkHandleHealthOptimized 优化健康检查基准
func BenchmarkHandleHealthOptimized(b *testing.B) {
	req := httptest.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleHealthOptimized(w, req)
	}
}

// BenchmarkHandleHealthCached 缓存健康检查基准
func BenchmarkHandleHealthCached(b *testing.B) {
	req := httptest.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handleHealthCached(w, req)
	}
}

// BenchmarkObjectPooling 对象池性能基准
func BenchmarkObjectPooling(b *testing.B) {
	b.Run("ResponsePool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resp := GetResponse()
			resp.Message = "test"
			PutResponse(resp)
		}
	})

	b.Run("BufferPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := GetBuffer()
			buf.WriteString("test data")
			PutBuffer(buf)
		}
	})

	b.Run("DataItemPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			item := GetDataItem()
			item["key"] = "value"
			PutDataItem(item)
		}
	})
}

// BenchmarkConcurrentPooling 并发对象池基准
func BenchmarkConcurrentPooling(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp := GetResponse()
			resp.Message = "concurrent test"
			PutResponse(resp)
		}
	})
}
