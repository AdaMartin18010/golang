//go:build !arenas
// +build !arenas

package memory

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestArenaMemoryLeak 测试内存泄漏检测
func TestArenaMemoryLeak(t *testing.T) {
	var m1, m2 runtime.MemStats

	// 第一次GC，建立基线
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// 执行多次处理
	for i := 0; i < 100; i++ {
		records := make([]Record, 100)
		for j := range records {
			records[j] = Record{
				ID:        j,
				Name:      "test",
				Value:     float64(j),
				Timestamp: time.Now(),
			}
		}
		_ = processWithArena(records)
	}

	// 第二次GC并读取内存
	runtime.GC()
	time.Sleep(50 * time.Millisecond)
	runtime.ReadMemStats(&m2)

	// 验证内存增长是否在合理范围内
	increase := m2.Alloc - m1.Alloc
	t.Logf("Memory increase: %d bytes", increase)

	// 内存增长应该相对较小（允许一些增长用于元数据）
	if increase > 10*1024*1024 { // 10MB
		t.Logf("Warning: Significant memory increase detected: %d MB", increase/1024/1024)
	}
}

// TestArenaConcurrentSafety 测试并发安全性
func TestArenaConcurrentSafety(t *testing.T) {
	const numGoroutines = 50
	const numRecordsPerGoroutine = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			records := make([]Record, numRecordsPerGoroutine)
			for j := range records {
				records[j] = Record{
					ID:        id*numRecordsPerGoroutine + j,
					Name:      "concurrent-test",
					Value:     float64(j),
					Timestamp: time.Now(),
				}
			}

			results := processWithArena(records)
			if len(results) != len(records) {
				errors <- nil
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// 检查是否有错误
	errorCount := 0
	for range errors {
		errorCount++
	}

	if errorCount > 0 {
		t.Errorf("Found %d errors in concurrent processing", errorCount)
	}
}

// TestArenaLargeDataset 测试大数据集处理
func TestArenaLargeDataset(t *testing.T) {
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		t.Run(string(rune(size)), func(t *testing.T) {
			records := make([]Record, size)
			for i := range records {
				records[i] = Record{
					ID:        i,
					Name:      "large-dataset",
					Value:     float64(i),
					Timestamp: time.Now(),
				}
			}

			start := time.Now()
			results := processWithArena(records)
			duration := time.Since(start)

			if len(results) != size {
				t.Errorf("Expected %d results, got %d", size, len(results))
			}

			t.Logf("Processed %d records in %v", size, duration)
		})
	}
}

// TestArenaEdgeCases 测试边界条件
func TestArenaEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		records []Record
		wantLen int
	}{
		{
			name:    "Empty slice",
			records: []Record{},
			wantLen: 0,
		},
		{
			name: "Single record",
			records: []Record{
				{ID: 1, Name: "single", Value: 1.0, Timestamp: time.Now()},
			},
			wantLen: 1,
		},
		{
			name: "Nil fields",
			records: []Record{
				{ID: 1, Name: "", Value: 0, Timestamp: time.Time{}},
			},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := processWithArena(tt.records)
			if len(results) != tt.wantLen {
				t.Errorf("Expected %d results, got %d", tt.wantLen, len(results))
			}
		})
	}
}

// TestTraditionalVsArenaConsistency 测试结果一致性
func TestTraditionalVsArenaConsistency(t *testing.T) {
	records := make([]Record, 1000)
	for i := range records {
		records[i] = Record{
			ID:        i,
			Name:      "consistency-test",
			Value:     float64(i) * 2.5,
			Timestamp: time.Now(),
		}
	}

	arenaResults := processWithArena(records)
	tradResults := processTraditional(records)

	if len(arenaResults) != len(tradResults) {
		t.Fatalf("Length mismatch: arena=%d, traditional=%d", len(arenaResults), len(tradResults))
	}

	// 验证每个结果
	for i := range arenaResults {
		if arenaResults[i].RecordID != tradResults[i].RecordID {
			t.Errorf("RecordID mismatch at index %d", i)
		}
		if arenaResults[i].Output != tradResults[i].Output {
			t.Errorf("Output mismatch at index %d", i)
		}
		if arenaResults[i].Score != tradResults[i].Score {
			t.Errorf("Score mismatch at index %d", i)
		}
	}
}

// TestBatchProcessorStress 测试BatchProcessor压力测试
func TestBatchProcessorStress(t *testing.T) {
	processor := NewBatchProcessor(100)

	// 生成大量数据
	records := make([]Record, 10000)
	for i := range records {
		records[i] = Record{
			ID:        i,
			Name:      "stress-test",
			Value:     float64(i),
			Timestamp: time.Now(),
		}
	}

	start := time.Now()
	results := processor.ProcessBatch(records)
	duration := time.Since(start)

	if len(results) != len(records) {
		t.Errorf("Expected %d results, got %d", len(records), len(results))
	}

	t.Logf("Batch processed %d records in %v", len(records), duration)
}

// TestRecordValidation 测试Record验证
func TestRecordValidation(t *testing.T) {
	tests := []struct {
		name    string
		record  Record
		wantErr bool
	}{
		{
			name: "Valid record",
			record: Record{
				ID:        1,
				Name:      "valid",
				Value:     100.0,
				Timestamp: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Negative ID",
			record: Record{
				ID:        -1,
				Name:      "negative",
				Value:     100.0,
				Timestamp: time.Now(),
			},
			wantErr: false, // 当前实现不验证，但测试覆盖
		},
		{
			name: "Zero timestamp",
			record: Record{
				ID:        1,
				Name:      "zero-time",
				Value:     100.0,
				Timestamp: time.Time{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records := []Record{tt.record}
			results := processWithArena(records)
			if len(results) != 1 {
				t.Errorf("Expected 1 result, got %d", len(results))
			}
		})
	}
}

// TestResultValidation 测试Result验证
func TestResultValidation(t *testing.T) {
	record := Record{
		ID:        123,
		Name:      "test",
		Value:     50.0,
		Timestamp: time.Now(),
	}

	records := []Record{record}
	results := processWithArena(records)

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	result := results[0]

	// 验证Result字段
	if result.RecordID != record.ID {
		t.Errorf("Expected RecordID %d, got %d", record.ID, result.RecordID)
	}

	expectedOutput := "Processed-" + record.Name
	if result.Output != expectedOutput {
		t.Errorf("Expected Output %s, got %s", expectedOutput, result.Output)
	}

	expectedScore := record.Value * 1.5
	if result.Score != expectedScore {
		t.Errorf("Expected Score %f, got %f", expectedScore, result.Score)
	}
}

// BenchmarkArenaLargeDataset 基准测试大数据集
func BenchmarkArenaLargeDataset(b *testing.B) {
	records := make([]Record, 10000)
	for i := range records {
		records[i] = Record{
			ID:        i,
			Name:      "benchmark",
			Value:     float64(i),
			Timestamp: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processWithArena(records)
	}
}

// BenchmarkTraditionalLargeDataset 基准测试传统方式大数据集
func BenchmarkTraditionalLargeDataset(b *testing.B) {
	records := make([]Record, 10000)
	for i := range records {
		records[i] = Record{
			ID:        i,
			Name:      "benchmark",
			Value:     float64(i),
			Timestamp: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processTraditional(records)
	}
}

// BenchmarkBatchProcessor 基准测试BatchProcessor
func BenchmarkBatchProcessor(b *testing.B) {
	processor := NewBatchProcessor(1000)
	records := make([]Record, 10000)
	for i := range records {
		records[i] = Record{
			ID:        i,
			Name:      "benchmark",
			Value:     float64(i),
			Timestamp: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processor.ProcessBatch(records)
	}
}

// TestMemoryUsageComparison 测试内存使用对比
func TestMemoryUsageComparison(t *testing.T) {
	records := make([]Record, 1000)
	for i := range records {
		records[i] = Record{
			ID:        i,
			Name:      "memory-test",
			Value:     float64(i),
			Timestamp: time.Now(),
		}
	}

	// 测试Arena方式
	runtime.GC()
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	for i := 0; i < 100; i++ {
		_ = processWithArena(records)
	}

	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	arenaAlloc := m2.TotalAlloc - m1.TotalAlloc

	// 测试传统方式
	runtime.GC()
	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)

	for i := 0; i < 100; i++ {
		_ = processTraditional(records)
	}

	runtime.GC()
	var m4 runtime.MemStats
	runtime.ReadMemStats(&m4)

	tradAlloc := m4.TotalAlloc - m3.TotalAlloc

	t.Logf("Arena allocation: %d MB", arenaAlloc/1024/1024)
	t.Logf("Traditional allocation: %d MB", tradAlloc/1024/1024)

	// Arena方式应该分配更少或相似的内存
	if arenaAlloc > tradAlloc*2 {
		t.Logf("Warning: Arena allocated significantly more memory")
	}
}
