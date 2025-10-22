package observability

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLoggerBasic(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(DebugLevel, &buf)

	logger.Info("Test message")

	output := buf.String()
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected output to contain 'Test message', got: %s", output)
	}

	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected output to contain 'INFO', got: %s", output)
	}
}

func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  LogLevel
		logFunc   func(*Logger)
		shouldLog bool
	}{
		{
			name:      "Debug at Debug level",
			logLevel:  DebugLevel,
			logFunc:   func(l *Logger) { l.Debug("debug") },
			shouldLog: true,
		},
		{
			name:      "Debug at Info level",
			logLevel:  InfoLevel,
			logFunc:   func(l *Logger) { l.Debug("debug") },
			shouldLog: false,
		},
		{
			name:      "Info at Info level",
			logLevel:  InfoLevel,
			logFunc:   func(l *Logger) { l.Info("info") },
			shouldLog: true,
		},
		{
			name:      "Warn at Error level",
			logLevel:  ErrorLevel,
			logFunc:   func(l *Logger) { l.Warn("warn") },
			shouldLog: false,
		},
		{
			name:      "Error at Warn level",
			logLevel:  WarnLevel,
			logFunc:   func(l *Logger) { l.Error("error") },
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(tt.logLevel, &buf)

			tt.logFunc(logger)

			output := buf.String()
			if tt.shouldLog && output == "" {
				t.Error("Expected log output, got none")
			}
			if !tt.shouldLog && output != "" {
				t.Errorf("Expected no log output, got: %s", output)
			}
		})
	}
}

func TestLoggerWithField(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	logger.WithField("user_id", "123").Info("User action")

	output := buf.String()
	if !strings.Contains(output, "user_id") {
		t.Errorf("Expected output to contain 'user_id', got: %s", output)
	}
	if !strings.Contains(output, "123") {
		t.Errorf("Expected output to contain '123', got: %s", output)
	}
}

func TestLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	fields := map[string]interface{}{
		"user_id":    "123",
		"action":     "login",
		"ip_address": "192.168.1.1",
	}

	logger.WithFields(fields).Info("User login")

	output := buf.String()
	for key, value := range fields {
		if !strings.Contains(output, key) {
			t.Errorf("Expected output to contain '%s', got: %s", key, output)
		}
		// 检查值（转换为字符串）
		valueStr := ""
		switch v := value.(type) {
		case string:
			valueStr = v
		}
		if valueStr != "" && !strings.Contains(output, valueStr) {
			t.Errorf("Expected output to contain '%s', got: %s", valueStr, output)
		}
	}
}

func TestLoggerWithContext(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	// 创建带有span的context
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)

	ctx := context.Background()
	span, ctx := tracer.StartSpan(ctx, "test-operation")

	logger.WithContext(ctx).Info("Operation started")

	output := buf.String()
	if !strings.Contains(output, "trace_id") {
		t.Errorf("Expected output to contain 'trace_id', got: %s", output)
	}
	if !strings.Contains(output, "span_id") {
		t.Errorf("Expected output to contain 'span_id', got: %s", output)
	}

	span.Finish()
}

func TestLoggerFormatted(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	logger.Infof("User %s performed action %s", "Alice", "logout")

	output := buf.String()
	if !strings.Contains(output, "User Alice performed action logout") {
		t.Errorf("Expected formatted message, got: %s", output)
	}
}

func TestMetricsHook(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	hook := NewMetricsHook()
	logger.AddHook(hook)

	// 记录几条日志
	logger.Info("Info message 1")
	logger.Info("Info message 2")
	logger.Warn("Warn message 1")
	logger.Error("Error message 1")

	// 验证指标
	time.Sleep(10 * time.Millisecond) // 给hook时间处理

	// 检查Info级别的计数器
	infoCounter := hook.counters[InfoLevel]
	if infoCounter.Get() < 2 {
		t.Errorf("Expected at least 2 info logs, got %d", infoCounter.Get())
	}

	warnCounter := hook.counters[WarnLevel]
	if warnCounter.Get() < 1 {
		t.Errorf("Expected at least 1 warn log, got %d", warnCounter.Get())
	}

	errorCounter := hook.counters[ErrorLevel]
	if errorCounter.Get() < 1 {
		t.Errorf("Expected at least 1 error log, got %d", errorCounter.Get())
	}
}

func TestLoggerConcurrent(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	const goroutines = 100
	const logsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < logsPerGoroutine; j++ {
				logger.WithField("goroutine", id).Infof("Log %d", j)
			}
		}(i)
	}

	wg.Wait()

	output := buf.String()
	if output == "" {
		t.Error("Expected log output, got none")
	}

	// 简单验证输出不为空
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < goroutines*logsPerGoroutine {
		t.Errorf("Expected at least %d log lines, got %d", goroutines*logsPerGoroutine, len(lines))
	}
}

func TestGlobalLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	// 设置为全局日志记录器
	SetDefaultLogger(logger)

	// 使用全局函数
	Info("Global info message")
	Warn("Global warn message")

	output := buf.String()
	if !strings.Contains(output, "Global info message") {
		t.Errorf("Expected output to contain 'Global info message', got: %s", output)
	}
	if !strings.Contains(output, "Global warn message") {
		t.Errorf("Expected output to contain 'Global warn message', got: %s", output)
	}
}

func TestWithFieldImmutability(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	// 创建带字段的logger
	logger1 := logger.WithField("field1", "value1")
	logger2 := logger1.WithField("field2", "value2")

	// logger1应该只有field1
	logger1.Info("Message 1")
	output1 := buf.String()
	buf.Reset()

	// logger2应该有field1和field2
	logger2.Info("Message 2")
	output2 := buf.String()

	// 验证field1在两个输出中
	if !strings.Contains(output1, "field1") {
		t.Error("Expected output1 to contain 'field1'")
	}
	if !strings.Contains(output2, "field1") {
		t.Error("Expected output2 to contain 'field1'")
	}

	// 验证field2只在output2中
	if strings.Contains(output1, "field2") {
		t.Error("Expected output1 to NOT contain 'field2'")
	}
	if !strings.Contains(output2, "field2") {
		t.Error("Expected output2 to contain 'field2'")
	}
}

func BenchmarkLoggerInfo(b *testing.B) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message")
	}
}

func BenchmarkLoggerInfoParallel(b *testing.B) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Benchmark message")
		}
	})
}

func BenchmarkLoggerWithField(b *testing.B) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithField("key", "value").Info("Benchmark message")
	}
}

func BenchmarkLoggerWithFields(b *testing.B) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(fields).Info("Benchmark message")
	}
}

func BenchmarkLoggerWithContext(b *testing.B) {
	var buf bytes.Buffer
	logger := NewLogger(InfoLevel, &buf)

	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)

	ctx := context.Background()
	span, ctx := tracer.StartSpan(ctx, "test-operation")
	defer span.Finish()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithContext(ctx).Info("Benchmark message")
	}
}
