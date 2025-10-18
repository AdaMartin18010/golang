// Arena Allocator示例：批量内存管理
// 注意：arena包在Go 1.25是实验性特性
// 构建时需要：GOEXPERIMENT=arenas go build
package main

import (
	"arena"
	"fmt"
	"runtime"
	"time"
)

// Record 数据记录
type Record struct {
	ID        int
	Name      string
	Value     float64
	Timestamp time.Time
}

// Result 处理结果
type Result struct {
	RecordID int
	Output   string
	Score    float64
}

// processWithArena 使用Arena处理批量数据
func processWithArena(records []Record) []Result {
	// 创建Arena
	a := arena.NewArena()
	defer a.Free() // 批量释放所有分配

	// 在Arena中分配结果
	results := make([]Result, len(records))

	for i, record := range records {
		// 在Arena中分配
		result := arena.New[Result](a)
		result.RecordID = record.ID
		result.Output = fmt.Sprintf("Processed-%s", record.Name)
		result.Score = record.Value * 1.5

		results[i] = *result
	}

	return results
}

// processTraditional 传统方式处理（对比）
func processTraditional(records []Record) []Result {
	results := make([]Result, len(records))

	for i, record := range records {
		// 每次分配新对象
		result := &Result{
			RecordID: record.ID,
			Output:   fmt.Sprintf("Processed-%s", record.Name),
			Score:    record.Value * 1.5,
		}

		results[i] = *result
	}

	return results
}

// benchmark 性能测试
func benchmark(name string, fn func([]Record) []Result, records []Record, rounds int) {
	var totalDuration time.Duration
	var totalGC uint32

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startGC := m.NumGC

	for i := 0; i < rounds; i++ {
		start := time.Now()
		_ = fn(records)
		totalDuration += time.Since(start)

		if i%10 == 0 {
			runtime.GC() // 定期触发GC
		}
	}

	runtime.ReadMemStats(&m)
	totalGC = m.NumGC - startGC

	avgDuration := totalDuration / time.Duration(rounds)

	fmt.Printf("\n📊 %s:\n", name)
	fmt.Printf("  Average time: %v\n", avgDuration)
	fmt.Printf("  Total time: %v\n", totalDuration)
	fmt.Printf("  GC count: %d\n", totalGC)
	fmt.Printf("  Alloc: %v MB\n", m.Alloc/1024/1024)
}

// BatchProcessor 批处理器
type BatchProcessor struct {
	batchSize int
	arena     *arena.Arena
}

func NewBatchProcessor(batchSize int) *BatchProcessor {
	return &BatchProcessor{
		batchSize: batchSize,
	}
}

// ProcessBatch 处理一批数据
func (p *BatchProcessor) ProcessBatch(records []Record) []Result {
	// 为每个batch创建新的Arena
	p.arena = arena.NewArena()
	defer p.arena.Free()

	results := make([]Result, 0, len(records))

	// 分批处理
	for i := 0; i < len(records); i += p.batchSize {
		end := i + p.batchSize
		if end > len(records) {
			end = len(records)
		}

		batch := records[i:end]
		batchResults := p.processBatch(batch)
		results = append(results, batchResults...)
	}

	return results
}

func (p *BatchProcessor) processBatch(batch []Record) []Result {
	results := make([]Result, len(batch))

	for i, record := range batch {
		// 在Arena中分配
		result := arena.New[Result](p.arena)
		result.RecordID = record.ID
		result.Output = fmt.Sprintf("Batch-Processed-%s", record.Name)
		result.Score = record.Value * 2.0

		results[i] = *result
	}

	return results
}

func main() {
	fmt.Println("🔬 Arena Allocator Demo\n")

	// 准备测试数据
	const numRecords = 10000
	records := make([]Record, numRecords)
	for i := 0; i < numRecords; i++ {
		records[i] = Record{
			ID:        i,
			Name:      fmt.Sprintf("Record-%d", i),
			Value:     float64(i) * 1.1,
			Timestamp: time.Now(),
		}
	}

	fmt.Printf("✅ Prepared %d records\n", numRecords)

	// === 测试1: 单次处理对比 ===
	fmt.Println("\n=== Test 1: Single Processing ===")

	// Arena方式
	start := time.Now()
	arenaResults := processWithArena(records)
	arenaDuration := time.Since(start)
	fmt.Printf("Arena: Processed %d records in %v\n",
		len(arenaResults), arenaDuration)

	// 传统方式
	start = time.Now()
	tradResults := processTraditional(records)
	tradDuration := time.Since(start)
	fmt.Printf("Traditional: Processed %d records in %v\n",
		len(tradResults), tradDuration)

	improvement := float64(tradDuration-arenaDuration) / float64(tradDuration) * 100
	fmt.Printf("💡 Arena is %.1f%% faster\n", improvement)

	// === 测试2: 批量处理 ===
	fmt.Println("\n=== Test 2: Batch Processing ===")

	processor := NewBatchProcessor(1000)
	start = time.Now()
	batchResults := processor.ProcessBatch(records)
	batchDuration := time.Since(start)

	fmt.Printf("Batch (Arena): Processed %d records in %v\n",
		len(batchResults), batchDuration)

	// === 测试3: 性能基准 ===
	fmt.Println("\n=== Test 3: Performance Benchmark ===")

	const rounds = 100
	fmt.Printf("Running %d rounds...\n", rounds)

	benchmark("Arena Allocator", processWithArena, records, rounds)
	benchmark("Traditional Allocator", processTraditional, records, rounds)

	// === 测试4: 内存使用对比 ===
	fmt.Println("\n=== Test 4: Memory Usage ===")

	// Arena方式
	runtime.GC()
	time.Sleep(50 * time.Millisecond)

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	for i := 0; i < 100; i++ {
		_ = processWithArena(records)
	}

	runtime.GC()
	time.Sleep(50 * time.Millisecond)

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	fmt.Printf("Arena:\n")
	fmt.Printf("  Before: %v MB\n", m1.Alloc/1024/1024)
	fmt.Printf("  After: %v MB\n", m2.Alloc/1024/1024)
	fmt.Printf("  GC count: %d\n", m2.NumGC-m1.NumGC)

	// 传统方式
	runtime.GC()
	time.Sleep(50 * time.Millisecond)

	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)

	for i := 0; i < 100; i++ {
		_ = processTraditional(records)
	}

	runtime.GC()
	time.Sleep(50 * time.Millisecond)

	var m4 runtime.MemStats
	runtime.ReadMemStats(&m4)

	fmt.Printf("\nTraditional:\n")
	fmt.Printf("  Before: %v MB\n", m3.Alloc/1024/1024)
	fmt.Printf("  After: %v MB\n", m4.Alloc/1024/1024)
	fmt.Printf("  GC count: %d\n", m4.NumGC-m3.NumGC)

	// === 使用建议 ===
	fmt.Println("\n📝 Usage Guidelines:")
	fmt.Println("✅ Use Arena when:")
	fmt.Println("  - Processing batches of short-lived objects")
	fmt.Println("  - All objects have the same lifecycle")
	fmt.Println("  - Want to minimize GC pressure")
	fmt.Println("  - Example: request handling, batch jobs")

	fmt.Println("\n❌ Don't use Arena when:")
	fmt.Println("  - Objects have different lifetimes")
	fmt.Println("  - Need to return long-lived objects")
	fmt.Println("  - Working with small number of objects")

	fmt.Println("\n✅ Demo completed!")
}
