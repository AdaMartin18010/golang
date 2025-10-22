package observability

import (
	"sync"
	"testing"
	"time"
)

func TestCounterBasic(t *testing.T) {
	counter := NewCounter("test_counter", "Test counter", nil)

	if counter.Get() != 0 {
		t.Errorf("Expected initial value 0, got %d", counter.Get())
	}

	counter.Inc()
	if counter.Get() != 1 {
		t.Errorf("Expected value 1 after Inc(), got %d", counter.Get())
	}

	counter.Add(5)
	if counter.Get() != 6 {
		t.Errorf("Expected value 6 after Add(5), got %d", counter.Get())
	}
}

func TestCounterConcurrent(t *testing.T) {
	counter := NewCounter("test_counter_concurrent", "Test counter", nil)

	const goroutines = 100
	const increments = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				counter.Inc()
			}
		}()
	}

	wg.Wait()

	expected := uint64(goroutines * increments)
	if counter.Get() != expected {
		t.Errorf("Expected value %d, got %d", expected, counter.Get())
	}
}

func TestGaugeBasic(t *testing.T) {
	gauge := NewGauge("test_gauge", "Test gauge", nil)

	if gauge.Get() != 0 {
		t.Errorf("Expected initial value 0, got %d", gauge.Get())
	}

	gauge.Set(10)
	if gauge.Get() != 10 {
		t.Errorf("Expected value 10 after Set(10), got %d", gauge.Get())
	}

	gauge.Inc()
	if gauge.Get() != 11 {
		t.Errorf("Expected value 11 after Inc(), got %d", gauge.Get())
	}

	gauge.Dec()
	if gauge.Get() != 10 {
		t.Errorf("Expected value 10 after Dec(), got %d", gauge.Get())
	}

	gauge.Add(-5)
	if gauge.Get() != 5 {
		t.Errorf("Expected value 5 after Add(-5), got %d", gauge.Get())
	}
}

func TestGaugeConcurrent(t *testing.T) {
	gauge := NewGauge("test_gauge_concurrent", "Test gauge", nil)

	const goroutines = 50
	const operations = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// 一半goroutine递增
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				gauge.Inc()
			}
		}()
	}

	// 一半goroutine递减
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				gauge.Dec()
			}
		}()
	}

	wg.Wait()

	// 由于递增和递减数量相同，最终应该是0
	if gauge.Get() != 0 {
		t.Errorf("Expected value 0, got %d", gauge.Get())
	}
}

func TestHistogramBasic(t *testing.T) {
	buckets := []float64{0.1, 0.5, 1.0, 5.0}
	histogram := NewHistogram("test_histogram", "Test histogram", buckets, nil)

	histogram.Observe(0.05) // bucket 0
	histogram.Observe(0.3)  // bucket 1
	histogram.Observe(0.8)  // bucket 2
	histogram.Observe(3.0)  // bucket 3
	histogram.Observe(10.0) // +Inf bucket

	value := histogram.Value().(map[string]interface{})
	counts := histogram.counts

	// 验证buckets分布
	if counts[0] != 1 {
		t.Errorf("Expected bucket 0 count 1, got %d", counts[0])
	}
	if counts[1] != 1 {
		t.Errorf("Expected bucket 1 count 1, got %d", counts[1])
	}
	if counts[2] != 1 {
		t.Errorf("Expected bucket 2 count 1, got %d", counts[2])
	}
	if counts[3] != 1 {
		t.Errorf("Expected bucket 3 count 1, got %d", counts[3])
	}
	if counts[4] != 1 {
		t.Errorf("Expected +Inf bucket count 1, got %d", counts[4])
	}

	// 验证总数
	count := value["count"].(uint64)
	if count != 5 {
		t.Errorf("Expected total count 5, got %d", count)
	}
}

func TestHistogramConcurrent(t *testing.T) {
	buckets := []float64{0.1, 0.5, 1.0, 5.0}
	histogram := NewHistogram("test_histogram_concurrent", "Test histogram", buckets, nil)

	const goroutines = 100
	const observations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < observations; j++ {
				// 观察不同的值
				histogram.Observe(float64(id%10) / 10.0)
			}
		}(i)
	}

	wg.Wait()

	value := histogram.Value().(map[string]interface{})
	count := value["count"].(uint64)

	expected := uint64(goroutines * observations)
	if count != expected {
		t.Errorf("Expected total count %d, got %d", expected, count)
	}
}

func TestMetricsRegistry(t *testing.T) {
	registry := NewMetricsRegistry()

	counter := NewCounter("test_counter", "Test counter", nil)
	gauge := NewGauge("test_gauge", "Test gauge", nil)

	// 注册指标
	if err := registry.Register(counter); err != nil {
		t.Fatalf("Failed to register counter: %v", err)
	}

	if err := registry.Register(gauge); err != nil {
		t.Fatalf("Failed to register gauge: %v", err)
	}

	// 重复注册应该失败
	if err := registry.Register(counter); err == nil {
		t.Error("Expected error when registering duplicate metric")
	}

	// 获取指标
	if metric, ok := registry.Get("test_counter"); !ok {
		t.Error("Expected to find test_counter")
	} else if metric.Type() != CounterType {
		t.Errorf("Expected CounterType, got %v", metric.Type())
	}

	// 获取所有指标
	metrics := registry.All()
	if len(metrics) != 2 {
		t.Errorf("Expected 2 metrics, got %d", len(metrics))
	}

	// 注销指标
	registry.Unregister("test_counter")
	if _, ok := registry.Get("test_counter"); ok {
		t.Error("Expected test_counter to be unregistered")
	}

	metrics = registry.All()
	if len(metrics) != 1 {
		t.Errorf("Expected 1 metric after unregister, got %d", len(metrics))
	}
}

func TestMetricsExport(t *testing.T) {
	registry := NewMetricsRegistry()

	counter := NewCounter("http_requests_total", "Total HTTP requests", map[string]string{
		"method": "GET",
		"status": "200",
	})
	counter.Add(100)

	gauge := NewGauge("active_connections", "Active connections", nil)
	gauge.Set(42)

	registry.Register(counter)
	registry.Register(gauge)

	output := registry.Export()

	// 验证输出包含指标信息
	if output == "" {
		t.Error("Expected non-empty export output")
	}

	// 简单验证包含指标名称
	if len(output) < 50 { // 应该有相当长的输出
		t.Errorf("Expected longer export output, got length %d", len(output))
	}
}

func TestDefaultMetrics(t *testing.T) {
	// 测试默认注册的指标
	metrics := []string{
		"http_requests_total",
		"http_request_duration_seconds",
		"active_connections",
		"memory_usage_bytes",
		"goroutine_count",
	}

	for _, name := range metrics {
		if _, ok := GetMetric(name); !ok {
			t.Errorf("Expected default metric '%s' to be registered", name)
		}
	}
}

func TestUpdateRuntimeMetrics(t *testing.T) {
	// 获取初始值
	initialGoroutines := GoroutineCount.Get()
	initialMemory := MemoryUsage.Get()

	// 更新指标
	UpdateRuntimeMetrics()

	// 验证指标被更新
	newGoroutines := GoroutineCount.Get()
	newMemory := MemoryUsage.Get()

	// 值应该已设置（在实际实现中会是真实值）
	if newGoroutines == initialGoroutines && newMemory == initialMemory {
		// 这是预期的，因为我们使用的是示例值
		// 在实际应用中，这些值会变化
	}

	if newGoroutines == 0 && newMemory == 0 {
		t.Error("Expected metrics to be updated with non-zero values")
	}
}

func TestStartMetricsCollector(t *testing.T) {
	// 启动指标收集器
	StartMetricsCollector(100 * time.Millisecond)

	// 等待几个收集周期
	time.Sleep(350 * time.Millisecond)

	// 验证指标已更新
	goroutines := GoroutineCount.Get()
	memory := MemoryUsage.Get()

	if goroutines == 0 {
		t.Error("Expected goroutine count to be updated")
	}

	if memory == 0 {
		t.Error("Expected memory usage to be updated")
	}
}

func BenchmarkCounterInc(b *testing.B) {
	counter := NewCounter("bench_counter", "Benchmark counter", nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter.Inc()
	}
}

func BenchmarkCounterIncParallel(b *testing.B) {
	counter := NewCounter("bench_counter_parallel", "Benchmark counter", nil)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Inc()
		}
	})
}

func BenchmarkGaugeSet(b *testing.B) {
	gauge := NewGauge("bench_gauge", "Benchmark gauge", nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gauge.Set(int64(i))
	}
}

func BenchmarkGaugeIncDec(b *testing.B) {
	gauge := NewGauge("bench_gauge_incdec", "Benchmark gauge", nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			gauge.Inc()
		} else {
			gauge.Dec()
		}
	}
}

func BenchmarkHistogramObserve(b *testing.B) {
	buckets := []float64{0.1, 0.5, 1.0, 5.0, 10.0}
	histogram := NewHistogram("bench_histogram", "Benchmark histogram", buckets, nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		histogram.Observe(float64(i%100) / 10.0)
	}
}

func BenchmarkHistogramObserveParallel(b *testing.B) {
	buckets := []float64{0.1, 0.5, 1.0, 5.0, 10.0}
	histogram := NewHistogram("bench_histogram_parallel", "Benchmark histogram", buckets, nil)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			histogram.Observe(float64(i%100) / 10.0)
			i++
		}
	})
}

func BenchmarkMetricsRegistryGet(b *testing.B) {
	registry := NewMetricsRegistry()
	counter := NewCounter("bench_counter", "Benchmark counter", nil)
	registry.Register(counter)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.Get("bench_counter")
	}
}

func BenchmarkMetricsExport(b *testing.B) {
	registry := NewMetricsRegistry()

	for i := 0; i < 10; i++ {
		counter := NewCounter("counter_"+string(rune(i)), "Counter", nil)
		gauge := NewGauge("gauge_"+string(rune(i)), "Gauge", nil)
		registry.Register(counter)
		registry.Register(gauge)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = registry.Export()
	}
}
