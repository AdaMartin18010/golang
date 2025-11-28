package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"ai-agent-architecture/core"
)

// 测试辅助函数
func NewSimpleLearningEngine() *core.LearningEngine {
	return core.NewLearningEngine(nil)
}

func NewSimpleDecisionEngine() *core.DecisionEngine {
	return core.NewDecisionEngine(nil)
}

// createMockAgentForTest 创建模拟Agent用于测试
func createMockAgentForTest(id string) core.Agent {
	return &mockTestAgent{
		id:     id,
		status: core.Status{ID: id, State: core.AgentStateRunning, Health: "healthy"},
	}
}

// mockTestAgent 实现Agent接口的简单模拟
type mockTestAgent struct {
	id     string
	status core.Status
}

func (m *mockTestAgent) ID() string {
	return m.id
}

func (m *mockTestAgent) Start(ctx context.Context) error {
	m.status.State = core.AgentStateRunning
	return nil
}

func (m *mockTestAgent) Stop() error {
	m.status.State = core.AgentStateStopped
	return nil
}

func (m *mockTestAgent) Process(input core.Input) (core.Output, error) {
	return core.Output{
		ID:        input.ID,
		Type:      "result",
		Data:      map[string]interface{}{"processed": true},
		Timestamp: time.Now(),
	}, nil
}

func (m *mockTestAgent) Learn(experience core.Experience) error {
	return nil
}

func (m *mockTestAgent) GetStatus() core.Status {
	return m.status
}

// SimpleMetricsCollector 简单的指标收集器实现
type SimpleMetricsCollector struct {
	metrics map[string]float64
	mu      sync.RWMutex
}

func NewSimpleMetricsCollector() core.MetricsCollector {
	return &SimpleMetricsCollector{
		metrics: make(map[string]float64),
	}
}

func (s *SimpleMetricsCollector) RecordProcess(duration time.Duration, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics["process_count"]++
	if success {
		s.metrics["success_count"]++
	}
}

func (s *SimpleMetricsCollector) RecordEvent(event string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics[event] = value
}

func (s *SimpleMetricsCollector) RecordMetric(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics[name] = value
}

func (s *SimpleMetricsCollector) GetMetrics() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]float64)
	for k, v := range s.metrics {
		result[k] = v
	}
	return result
}

func (s *SimpleMetricsCollector) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics = make(map[string]float64)
}

// TestBaseAgent 基础代理测试
func TestBaseAgent(t *testing.T) {
	// 创建基础代理
	config := core.AgentConfig{
		Name:         "测试代理",
		Type:         "test",
		MaxLoad:      0.8,
		Timeout:      5 * time.Second,
		Retries:      3,
		Capabilities: []string{"test_processing"},
		Parameters:   map[string]interface{}{"version": "1.0"},
	}

	agent := core.NewBaseAgent("test-agent-1", config)

	// 创建决策引擎并注册一个模拟代理
	decisionEngine := NewSimpleDecisionEngine()
	mockAgent := createMockAgentForTest("mock-executor")
	if err := decisionEngine.RegisterAgent(&mockAgent); err != nil {
		t.Fatalf("Failed to register mock agent: %v", err)
	}

	agent.SetLearningEngine(NewSimpleLearningEngine())
	agent.SetDecisionEngine(decisionEngine)
	agent.SetMetricsCollector(NewSimpleMetricsCollector())

	ctx := context.Background()

	// 测试启动代理
	err := agent.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start agent: %v", err)
	}
	defer agent.Stop()

	// 测试代理状态
	status := agent.GetStatus()
	if status.State != core.AgentStateRunning {
		t.Errorf("Expected agent state to be running, got %s", status.State)
	}

	// 测试处理任务
	input := core.Input{
		ID:   "test-task-1",
		Type: "test_task",
		Data: map[string]interface{}{"value": 42},
		Metadata: map[string]interface{}{
			"priority": 1,
			"timeout":  "2s",
		},
		Timestamp: time.Now(),
	}

	output, err := agent.Process(input)
	if err != nil {
		t.Fatalf("Failed to process task: %v", err)
	}

	if output.Error != nil {
		t.Errorf("Expected task to succeed, but got error: %v", output.Error)
	}

	// 输出应该有有效的ID（可能与输入ID不同，取决于决策引擎的实现）
	if output.ID == "" {
		t.Error("Expected non-empty output ID")
	}

	// 测试学习功能
	experience := core.Experience{
		Input:     input,
		Output:    output,
		Reward:    1.0,
		Timestamp: time.Now(),
	}

	err = agent.Learn(experience)
	if err != nil {
		t.Fatalf("Failed to learn from experience: %v", err)
	}
}

// TestAgentConcurrency 测试代理并发处理
func TestAgentConcurrency(t *testing.T) {
	config := core.AgentConfig{
		Name:         "并发测试代理",
		Type:         "concurrent",
		MaxLoad:      0.9,
		Timeout:      10 * time.Second,
		Retries:      2,
		Capabilities: []string{"concurrent_processing"},
		Parameters:   map[string]interface{}{"workers": 5},
	}

	agent := core.NewBaseAgent("concurrent-agent", config)

	// 创建决策引擎并注册多个模拟代理
	decisionEngine := NewSimpleDecisionEngine()
	for i := 0; i < 3; i++ {
		mockAgent := createMockAgentForTest(string(rune('A' + i)))
		if err := decisionEngine.RegisterAgent(&mockAgent); err != nil {
			t.Fatalf("Failed to register mock agent: %v", err)
		}
	}

	agent.SetLearningEngine(NewSimpleLearningEngine())
	agent.SetDecisionEngine(decisionEngine)
	agent.SetMetricsCollector(NewSimpleMetricsCollector())

	ctx := context.Background()
	err := agent.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start agent: %v", err)
	}
	defer agent.Stop()

	// 并发处理多个任务
	const numTasks = 10
	errors := make(chan error, numTasks)
	var wg sync.WaitGroup

	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			input := core.Input{
				ID:   string(rune(taskID)),
				Type: "test_task",
				Data: map[string]interface{}{"task_id": taskID},
				Metadata: map[string]interface{}{
					"priority": 1,
				},
				Timestamp: time.Now(),
			}

			_, err := agent.Process(input)
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// 检查错误
	for err := range errors {
		t.Errorf("Task processing error: %v", err)
	}
}

// BenchmarkAgentProcess 性能测试
func BenchmarkAgentProcess(b *testing.B) {
	config := core.AgentConfig{
		Name:         "基准测试代理",
		Type:         "benchmark",
		MaxLoad:      1.0,
		Timeout:      30 * time.Second,
		Retries:      1,
		Capabilities: []string{"benchmark_processing"},
		Parameters:   map[string]interface{}{"mode": "fast"},
	}

	agent := core.NewBaseAgent("benchmark-agent", config)

	// 创建决策引擎并注册模拟代理
	decisionEngine := NewSimpleDecisionEngine()
	mockAgent := createMockAgentForTest("benchmark-executor")
	_ = decisionEngine.RegisterAgent(&mockAgent)

	agent.SetLearningEngine(NewSimpleLearningEngine())
	agent.SetDecisionEngine(decisionEngine)
	agent.SetMetricsCollector(NewSimpleMetricsCollector())

	ctx := context.Background()
	agent.Start(ctx)
	defer agent.Stop()

	input := core.Input{
		ID:   "benchmark-task",
		Type: "benchmark_task",
		Data: map[string]interface{}{"value": 42},
		Metadata: map[string]interface{}{
			"priority": 1,
			"timeout":  "1s",
		},
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := agent.Process(input)
		if err != nil {
			b.Fatalf("Failed to process task: %v", err)
		}
	}
}
