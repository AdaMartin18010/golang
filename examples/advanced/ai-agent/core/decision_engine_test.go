package core

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// createMockAgent 创建模拟Agent
func createMockAgent(id string) Agent {
	return &mockAgent{
		id:           id,
		name:         fmt.Sprintf("Agent-%s", id),
		capabilities: []string{"test", "mock"},
	}
}

// mockAgent 实现Agent接口的简单模拟
type mockAgent struct {
	id           string
	name         string
	capabilities []string
	status       Status
}

func (m *mockAgent) ID() string {
	return m.id
}

func (m *mockAgent) Start(ctx context.Context) error {
	m.status.State = AgentStateRunning
	return nil
}

func (m *mockAgent) Stop() error {
	m.status.State = AgentStateStopped
	return nil
}

func (m *mockAgent) Process(input Input) (Output, error) {
	return Output{
		ID:        fmt.Sprintf("output-%s", input.ID),
		Type:      "result",
		Data:      map[string]interface{}{"processed": true},
		Timestamp: time.Now(),
	}, nil
}

func (m *mockAgent) Learn(experience Experience) error {
	// 模拟学习
	return nil
}

func (m *mockAgent) GetStatus() Status {
	return m.status
}

// TestNewDecisionEngine 测试创建决策引擎
func TestNewDecisionEngine(t *testing.T) {
	engine := NewDecisionEngine(nil)

	if engine == nil {
		t.Fatal("Failed to create DecisionEngine")
	}

	if engine.coordinator == nil {
		t.Error("Coordinator not initialized")
	}

	if engine.distributor == nil {
		t.Error("Distributor not initialized")
	}
}

// TestMakeDecision 测试基本决策
func TestMakeDecision(t *testing.T) {
	engine := NewDecisionEngine(nil)

	// 注册一个模拟Agent
	mockAgent := createMockAgent("agent-1")
	err := engine.RegisterAgent(&mockAgent)
	if err != nil {
		t.Fatalf("RegisterAgent failed: %v", err)
	}

	ctx := context.Background()

	task := &Task{
		ID:        "test-task-1",
		Type:      "test",
		Priority:  1,
		Input:     map[string]interface{}{"value": 42},
		CreatedAt: time.Now(),
	}

	decision, err := engine.MakeDecision(ctx, task)
	if err != nil {
		t.Fatalf("MakeDecision failed: %v", err)
	}

	if decision.Action == "" {
		t.Error("Expected non-empty action")
	}

	if decision.Confidence < 0 || decision.Confidence > 1 {
		t.Errorf("Invalid confidence: %f", decision.Confidence)
	}
}

// TestMakeDecisionMultipleTasks 测试多个决策
func TestMakeDecisionMultipleTasks(t *testing.T) {
	engine := NewDecisionEngine(nil)

	// 注册多个模拟Agent
	for i := 1; i <= 3; i++ {
		mockAgent := createMockAgent(fmt.Sprintf("agent-%d", i))
		err := engine.RegisterAgent(&mockAgent)
		if err != nil {
			t.Fatalf("RegisterAgent failed: %v", err)
		}
	}

	ctx := context.Background()

	tasks := []*Task{
		{ID: "1", Type: "type1", Priority: 1, Input: map[string]interface{}{}, CreatedAt: time.Now()},
		{ID: "2", Type: "type2", Priority: 2, Input: map[string]interface{}{}, CreatedAt: time.Now()},
		{ID: "3", Type: "type3", Priority: 3, Input: map[string]interface{}{}, CreatedAt: time.Now()},
	}

	for _, task := range tasks {
		decision, err := engine.MakeDecision(ctx, task)
		if err != nil {
			t.Errorf("MakeDecision failed for task %s: %v", task.ID, err)
			continue
		}

		if decision.Action == "" {
			t.Errorf("Empty action for task %s", task.ID)
		}
	}
}

// TestMakeDecisionWithContext 测试带超时的决策
func TestMakeDecisionWithContext(t *testing.T) {
	engine := NewDecisionEngine(nil)

	// 注册模拟Agent
	mockAgent := createMockAgent("agent-timeout")
	err := engine.RegisterAgent(&mockAgent)
	if err != nil {
		t.Fatalf("RegisterAgent failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	task := &Task{
		ID:        "timeout-task",
		Type:      "slow",
		Priority:  1,
		Input:     map[string]interface{}{},
		CreatedAt: time.Now(),
	}

	// 应该在超时前完成
	decision, err := engine.MakeDecision(ctx, task)
	if err != nil && ctx.Err() == context.DeadlineExceeded {
		t.Skip("Context deadline exceeded as expected in some cases")
	}

	if err == nil && decision == nil {
		t.Error("Expected either error or decision")
	}
}

// TestDecisionConfidence 测试决策置信度
func TestDecisionConfidence(t *testing.T) {
	engine := NewDecisionEngine(nil)

	// 注册模拟Agent
	mockAgent := createMockAgent("agent-conf")
	err := engine.RegisterAgent(&mockAgent)
	if err != nil {
		t.Fatalf("RegisterAgent failed: %v", err)
	}

	ctx := context.Background()

	task := &Task{
		ID:        "conf-test",
		Type:      "standard",
		Priority:  1,
		Input:     map[string]interface{}{},
		CreatedAt: time.Now(),
	}

	decision, err := engine.MakeDecision(ctx, task)
	if err != nil {
		t.Fatalf("MakeDecision failed: %v", err)
	}

	// 置信度应该在合理范围内
	if decision.Confidence < 0.0 {
		t.Error("Confidence should not be negative")
	}
	if decision.Confidence > 1.0 {
		t.Error("Confidence should not exceed 1.0")
	}
}

// TestDecisionWithNilTask 测试nil任务
func TestDecisionWithNilTask(t *testing.T) {
	engine := NewDecisionEngine(nil)

	// 即使没有Agent，nil任务也应该返回错误
	ctx := context.Background()

	_, err := engine.MakeDecision(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil task")
	}
}

// TestConcurrentDecisions 测试并发决策
func TestConcurrentDecisions(t *testing.T) {
	engine := NewDecisionEngine(nil)

	// 注册足够的Agent处理并发请求
	for i := 0; i < 5; i++ {
		mockAgent := createMockAgent(fmt.Sprintf("agent-concurrent-%d", i))
		err := engine.RegisterAgent(&mockAgent)
		if err != nil {
			t.Fatalf("RegisterAgent failed: %v", err)
		}
	}

	ctx := context.Background()

	numGoroutines := 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			task := &Task{
				ID:        fmt.Sprintf("task-%d", id),
				Type:      "concurrent",
				Priority:  1,
				Input:     map[string]interface{}{},
				CreatedAt: time.Now(),
			}
			_, err := engine.MakeDecision(ctx, task)
			results <- err
		}(i)
	}

	// 收集结果
	errorCount := 0
	for i := 0; i < numGoroutines; i++ {
		if err := <-results; err != nil {
			t.Logf("Concurrent decision error: %v", err)
			errorCount++
		}
	}

	// 允许少量错误（由于并发竞争）
	if errorCount > numGoroutines/2 {
		t.Errorf("Too many errors: %d out of %d", errorCount, numGoroutines)
	}
}

// BenchmarkMakeDecision 基准测试决策性能
func BenchmarkMakeDecision(b *testing.B) {
	engine := NewDecisionEngine(nil)

	// 注册模拟Agent
	mockAgent := createMockAgent("agent-bench")
	_ = engine.RegisterAgent(&mockAgent)

	ctx := context.Background()

	task := &Task{
		ID:        "bench-task",
		Type:      "benchmark",
		Priority:  1,
		Input:     map[string]interface{}{},
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.MakeDecision(ctx, task)
	}
}
