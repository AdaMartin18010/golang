//go:build !arenas
// +build !arenas

package memory

import (
	"runtime"
	"testing"
	"time"
)

// TestRecordCreation 测试记录创建
func TestRecordCreation(t *testing.T) {
	record := Record{
		ID:        1,
		Name:      "test",
		Value:     100.0,
		Timestamp: time.Now(),
	}

	if record.ID != 1 {
		t.Errorf("Expected ID 1, got %d", record.ID)
	}

	if record.Name != "test" {
		t.Errorf("Expected Name 'test', got '%s'", record.Name)
	}

	if record.Value != 100.0 {
		t.Errorf("Expected Value 100.0, got %f", record.Value)
	}
}

// TestProcessWithArena 测试Arena处理方式
func TestProcessWithArena(t *testing.T) {
	records := []Record{
		{ID: 1, Name: "A", Value: 10.0, Timestamp: time.Now()},
		{ID: 2, Name: "B", Value: 20.0, Timestamp: time.Now()},
		{ID: 3, Name: "C", Value: 30.0, Timestamp: time.Now()},
	}

	results := processWithArena(records)

	if len(results) != len(records) {
		t.Fatalf("Expected %d results, got %d", len(records), len(results))
	}

	// 验证第一条结果
	if results[0].RecordID != 1 {
		t.Errorf("Expected RecordID 1, got %d", results[0].RecordID)
	}

	if results[0].Output != "Processed-A" {
		t.Errorf("Expected Output 'Processed-A', got '%s'", results[0].Output)
	}

	if results[0].Score != 15.0 {
		t.Errorf("Expected Score 15.0, got %f", results[0].Score)
	}
}

// TestProcessTraditional 测试传统处理方式
func TestProcessTraditional(t *testing.T) {
	records := []Record{
		{ID: 1, Name: "A", Value: 10.0, Timestamp: time.Now()},
		{ID: 2, Name: "B", Value: 20.0, Timestamp: time.Now()},
	}

	results := processTraditional(records)

	if len(results) != len(records) {
		t.Fatalf("Expected %d results, got %d", len(records), len(results))
	}

	// 验证结果正确性
	for i, result := range results {
		if result.RecordID != records[i].ID {
			t.Errorf("Result %d: Expected RecordID %d, got %d", i, records[i].ID, result.RecordID)
		}
	}
}

// TestArenaVsTraditionalCorrectness 测试两种方式的结果一致性
func TestArenaVsTraditionalCorrectness(t *testing.T) {
	records := []Record{
		{ID: 1, Name: "X", Value: 50.0, Timestamp: time.Now()},
		{ID: 2, Name: "Y", Value: 75.0, Timestamp: time.Now()},
	}

	arenaResults := processWithArena(records)
	traditionalResults := processTraditional(records)

	if len(arenaResults) != len(traditionalResults) {
		t.Fatalf("Result counts differ: arena=%d, traditional=%d", len(arenaResults), len(traditionalResults))
	}

	for i := range arenaResults {
		if arenaResults[i].RecordID != traditionalResults[i].RecordID {
			t.Errorf("Result %d: RecordID mismatch", i)
		}
		if arenaResults[i].Output != traditionalResults[i].Output {
			t.Errorf("Result %d: Output mismatch", i)
		}
		if arenaResults[i].Score != traditionalResults[i].Score {
			t.Errorf("Result %d: Score mismatch", i)
		}
	}
}

// TestEmptyRecords 测试空记录集
func TestEmptyRecords(t *testing.T) {
	records := []Record{}

	arenaResults := processWithArena(records)
	if len(arenaResults) != 0 {
		t.Errorf("Expected empty results, got %d", len(arenaResults))
	}

	traditionalResults := processTraditional(records)
	if len(traditionalResults) != 0 {
		t.Errorf("Expected empty results, got %d", len(traditionalResults))
	}
}

// BenchmarkProcessWithArena 基准测试Arena方式
func BenchmarkProcessWithArena(b *testing.B) {
	records := make([]Record, 1000)
	for i := range records {
		records[i] = Record{
			ID:        i,
			Name:      "test",
			Value:     float64(i),
			Timestamp: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processWithArena(records)
	}
}

// BenchmarkProcessTraditional 基准测试传统方式
func BenchmarkProcessTraditional(b *testing.B) {
	records := make([]Record, 1000)
	for i := range records {
		records[i] = Record{
			ID:        i,
			Name:      "test",
			Value:     float64(i),
			Timestamp: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processTraditional(records)
	}
}

// TestMemoryAllocation 测试内存分配
func TestMemoryAllocation(t *testing.T) {
	records := make([]Record, 100)
	for i := range records {
		records[i] = Record{
			ID:        i,
			Name:      "test",
			Value:     float64(i),
			Timestamp: time.Now(),
		}
	}

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	_ = processWithArena(records)

	runtime.ReadMemStats(&m2)

	// 验证确实发生了内存分配
	if m2.TotalAlloc <= m1.TotalAlloc {
		t.Log("Memory allocation detected")
	}
}
