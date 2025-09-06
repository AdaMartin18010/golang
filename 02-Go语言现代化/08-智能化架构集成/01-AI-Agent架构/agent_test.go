package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestBaseAgent(t *testing.T) {
	// 创建基础代理
	config := AgentConfig{
		Name:         "测试代理",
		Type:         "test",
		MaxLoad:      0.8,
		Timeout:      5 * time.Second,
		Retries:      3,
		Capabilities: []string{"test_processing"},
		Parameters:   map[string]interface{}{"version": "1.0"},
	}

	agent := NewBaseAgent("test-agent-1", config)
	agent.learning = NewSimpleLearningEngine()
	agent.decision = NewSimpleDecisionEngine()
	agent.metrics = NewSimpleMetricsCollector()

	ctx := context.Background()

	// 测试启动代理
	err := agent.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start agent: %v", err)
	}
	defer agent.Stop()

	// 测试代理状态
	status := agent.GetStatus()
	if status.State != AgentStateRunning {
		t.Errorf("Expected agent state to be running, got %s", status.State)
	}

	// 测试处理任务
	input := Input{
		ID:       "test-task-1",
		Type:     "test_task",
		Data:     map[string]interface{}{"value": 42},
		Priority: 1,
		Timeout:  2 * time.Second,
	}

	output, err := agent.Process(ctx, input)
	if err != nil {
		t.Fatalf("Failed to process task: %v", err)
	}

	if !output.Success {
		t.Errorf("Expected task to succeed, but it failed")
	}

	if output.ID != input.ID {
		t.Errorf("Expected output ID %s, got %s", input.ID, output.ID)
	}

	// 测试学习功能
	experience := Experience{
		Input:     input,
		Output:    output,
		Feedback:  1.0,
		Timestamp: time.Now(),
	}

	err = agent.Learn(experience)
	if err != nil {
		t.Errorf("Failed to learn from experience: %v", err)
	}

	// 测试获取能力
	capabilities := agent.GetCapabilities()
	if len(capabilities) == 0 {
		t.Errorf("Expected capabilities, got empty slice")
	}
}

func TestDataProcessingAgent(t *testing.T) {
	config := AgentConfig{
		Name:         "数据处理测试代理",
		Type:         "data_processing",
		MaxLoad:      0.8,
		Timeout:      5 * time.Second,
		Retries:      2,
		Capabilities: []string{"data_processing", "etl", "validation"},
		Parameters:   map[string]interface{}{"batch_size": 100},
	}

	agent := NewDataProcessingAgent("data-agent-1", config)
	agent.learning = NewSimpleLearningEngine()
	agent.decision = NewSimpleDecisionEngine()
	agent.metrics = NewSimpleMetricsCollector()

	ctx := context.Background()

	err := agent.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start data processing agent: %v", err)
	}
	defer agent.Stop()

	// 测试数据处理
	input := Input{
		ID:       "data-task-1",
		Type:     "data_processing",
		Data:     map[string]interface{}{"records": 1000, "format": "csv"},
		Priority: 4,
		Timeout:  3 * time.Second,
	}

	output, err := agent.Process(ctx, input)
	if err != nil {
		t.Fatalf("Failed to process data task: %v", err)
	}

	if !output.Success {
		t.Errorf("Expected data processing to succeed, but it failed")
	}

	// 检查输出数据
	if output.Data["agent_type"] != "data_processing" {
		t.Errorf("Expected agent_type to be 'data_processing', got %v", output.Data["agent_type"])
	}
}

func TestDecisionAgent(t *testing.T) {
	config := AgentConfig{
		Name:         "决策测试代理",
		Type:         "decision_making",
		MaxLoad:      0.6,
		Timeout:      3 * time.Second,
		Retries:      1,
		Capabilities: []string{"decision_making", "rule_engine", "optimization"},
		Parameters:   map[string]interface{}{"confidence_threshold": 0.8},
	}

	agent := NewDecisionAgent("decision-agent-1", config)
	agent.learning = NewSimpleLearningEngine()
	agent.decision = NewSimpleDecisionEngine()
	agent.metrics = NewSimpleMetricsCollector()

	ctx := context.Background()

	err := agent.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start decision agent: %v", err)
	}
	defer agent.Stop()

	// 添加一些规则
	rule := Rule{
		Condition: "high_priority",
		Action:    "process_immediately",
		Weight:    0.9,
	}
	agent.rules.AddRule(rule)

	// 添加策略
	policy := Policy{
		ID:       "test_policy",
		Name:     "测试策略",
		Rules:    []Rule{rule},
		Priority: 1,
		Enabled:  true,
	}
	agent.policies.AddPolicy(policy)

	// 测试决策
	input := Input{
		ID:       "decision-task-1",
		Type:     "decision_making",
		Data:     map[string]interface{}{"scenario": "investment", "risk_level": "medium"},
		Priority: 6, // 高优先级
		Timeout:  2 * time.Second,
	}

	output, err := agent.Process(ctx, input)
	if err != nil {
		t.Fatalf("Failed to process decision task: %v", err)
	}

	if !output.Success {
		t.Errorf("Expected decision making to succeed, but it failed")
	}
}

func TestMonitoringAgent(t *testing.T) {
	config := AgentConfig{
		Name:         "监控测试代理",
		Type:         "monitoring",
		MaxLoad:      0.5,
		Timeout:      2 * time.Second,
		Retries:      3,
		Capabilities: []string{"monitoring", "anomaly_detection", "alerting"},
		Parameters:   map[string]interface{}{"check_interval": "5s"},
	}

	agent := NewMonitoringAgent("monitoring-agent-1", config)
	agent.learning = NewSimpleLearningEngine()
	agent.decision = NewSimpleDecisionEngine()
	agent.metrics = NewSimpleMetricsCollector()

	ctx := context.Background()

	err := agent.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start monitoring agent: %v", err)
	}
	defer agent.Stop()

	// 测试监控
	input := Input{
		ID:       "monitoring-task-1",
		Type:     "monitoring",
		Data:     map[string]interface{}{"targets": []string{"server1", "server2"}},
		Priority: 3,
		Timeout:  1 * time.Second,
	}

	output, err := agent.Process(ctx, input)
	if err != nil {
		t.Fatalf("Failed to process monitoring task: %v", err)
	}

	if !output.Success {
		t.Errorf("Expected monitoring to succeed, but it failed")
	}

	// 检查输出数据
	if output.Data["agent_type"] != "monitoring" {
		t.Errorf("Expected agent_type to be 'monitoring', got %v", output.Data["agent_type"])
	}
}

func TestSmartCoordinator(t *testing.T) {
	coordinator := NewSmartCoordinator()
	err := coordinator.Start()
	if err != nil {
		t.Fatalf("Failed to start coordinator: %v", err)
	}
	defer coordinator.Stop()

	// 注册代理
	config := AgentConfig{
		Name:         "测试代理",
		Type:         "test",
		MaxLoad:      0.8,
		Timeout:      5 * time.Second,
		Retries:      2,
		Capabilities: []string{"test_processing"},
		Parameters:   map[string]interface{}{"version": "1.0"},
	}

	agent := NewBaseAgent("test-agent-coord", config)
	agent.learning = NewSimpleLearningEngine()
	agent.decision = NewSimpleDecisionEngine()
	agent.metrics = NewSimpleMetricsCollector()

	err = coordinator.RegisterAgent(agent)
	if err != nil {
		t.Fatalf("Failed to register agent: %v", err)
	}
	defer coordinator.UnregisterAgent(agent.ID())

	// 测试任务处理
	task := Task{
		ID:        "coordinator-test-task",
		Type:      "test_task",
		Priority:  3,
		Data:      map[string]interface{}{"test": true},
		Timeout:   2 * time.Second,
		CreatedAt: time.Now(),
	}

	result, err := coordinator.ProcessTask(task)
	if err != nil {
		t.Fatalf("Failed to process task: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected task to succeed, but it failed")
	}

	if result.AgentID != agent.ID() {
		t.Errorf("Expected agent ID %s, got %s", agent.ID(), result.AgentID)
	}

	// 测试系统状态
	status := coordinator.GetSystemStatus()
	if status.TotalAgents != 1 {
		t.Errorf("Expected 1 total agent, got %d", status.TotalAgents)
	}

	if status.ActiveAgents != 1 {
		t.Errorf("Expected 1 active agent, got %d", status.ActiveAgents)
	}
}

// 基准测试
func BenchmarkBaseAgentProcess(b *testing.B) {
	config := AgentConfig{
		Name:         "基准测试代理",
		Type:         "benchmark",
		MaxLoad:      0.8,
		Timeout:      5 * time.Second,
		Retries:      1,
		Capabilities: []string{"benchmark_processing"},
		Parameters:   map[string]interface{}{"version": "1.0"},
	}

	agent := NewBaseAgent("benchmark-agent", config)
	agent.learning = NewSimpleLearningEngine()
	agent.decision = NewSimpleDecisionEngine()
	agent.metrics = NewSimpleMetricsCollector()

	ctx := context.Background()
	agent.Start(ctx)
	defer agent.Stop()

	input := Input{
		ID:       "benchmark-task",
		Type:     "benchmark_task",
		Data:     map[string]interface{}{"value": 42},
		Priority: 1,
		Timeout:  1 * time.Second,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input.ID = fmt.Sprintf("benchmark-task-%d", i)
		_, err := agent.Process(ctx, input)
		if err != nil {
			b.Fatalf("Failed to process task: %v", err)
		}
	}
}

func BenchmarkDataProcessingAgent(b *testing.B) {
	config := AgentConfig{
		Name:         "基准数据处理代理",
		Type:         "data_processing",
		MaxLoad:      0.8,
		Timeout:      5 * time.Second,
		Retries:      1,
		Capabilities: []string{"data_processing"},
		Parameters:   map[string]interface{}{"batch_size": 100},
	}

	agent := NewDataProcessingAgent("benchmark-data-agent", config)
	agent.learning = NewSimpleLearningEngine()
	agent.decision = NewSimpleDecisionEngine()
	agent.metrics = NewSimpleMetricsCollector()

	ctx := context.Background()
	agent.Start(ctx)
	defer agent.Stop()

	input := Input{
		ID:       "benchmark-data-task",
		Type:     "data_processing",
		Data:     map[string]interface{}{"records": 100},
		Priority: 1,
		Timeout:  1 * time.Second,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input.ID = fmt.Sprintf("benchmark-data-task-%d", i)
		_, err := agent.Process(ctx, input)
		if err != nil {
			b.Fatalf("Failed to process data task: %v", err)
		}
	}
}
