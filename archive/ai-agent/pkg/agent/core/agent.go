package core

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 类型定义已移至types.go

// BaseAgent 基础代理实现
type BaseAgent struct {
	id       string
	state    AgentState
	config   AgentConfig
	learning *LearningEngine
	decision *DecisionEngine
	metrics  MetricsCollector
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewBaseAgent 创建新的基础代理
func NewBaseAgent(id string, config AgentConfig) *BaseAgent {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseAgent{
		id:     id,
		state:  AgentStateIdle,
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

// ID 返回代理ID
func (a *BaseAgent) ID() string {
	return a.id
}

// Start 启动代理
func (a *BaseAgent) Start(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.state == AgentStateRunning {
		return fmt.Errorf("agent %s is already running", a.id)
	}

	// 初始化组件
	if err := a.initialize(); err != nil {
		a.state = AgentStateError
		return fmt.Errorf("failed to initialize agent: %w", err)
	}

	a.state = AgentStateRunning
	if a.metrics != nil {
		a.metrics.RecordEvent("agent_started", 1.0)
	}

	return nil
}

// Stop 停止代理
func (a *BaseAgent) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.state == AgentStateStopped {
		return nil
	}

	// 取消上下文
	if a.cancel != nil {
		a.cancel()
	}

	// 清理资源
	if err := a.cleanup(); err != nil {
		return fmt.Errorf("failed to cleanup agent: %w", err)
	}

	a.state = AgentStateStopped
	if a.metrics != nil {
		a.metrics.RecordEvent("agent_stopped", 1.0)
	}

	return nil
}

// Process 处理输入数据（简化版本）
func (a *BaseAgent) Process(input Input) (Output, error) {
	a.mu.RLock()
	if a.state != AgentStateRunning {
		a.mu.RUnlock()
		return Output{}, fmt.Errorf("agent %s is not running", a.id)
	}
	a.mu.RUnlock()

	startTime := time.Now()
	defer func() {
		if a.metrics != nil {
			duration := time.Since(startTime)
			a.metrics.RecordMetric("processing_time", duration.Seconds())
		}
	}()

	// 简化处理：直接转换Input为Task并调用决策引擎
	task := &Task{
		ID:         input.ID,
		Type:       input.Type,
		Priority:   1,
		Complexity: 0.5,
		CreatedAt:  input.Timestamp,
	}

	// 做出决策
	decision, err := a.decision.MakeDecision(a.ctx, task)
	if err != nil {
		if a.metrics != nil {
			a.metrics.RecordEvent("decision_error", 1.0)
		}
		return Output{}, fmt.Errorf("failed to make decision: %w", err)
	}

	// 执行决策
	output := a.execute(*decision)

	// 学习结果
	a.learnFromOutcome(*decision, output)

	// 记录成功处理
	if a.metrics != nil {
		a.metrics.RecordEvent("processing_success", 1.0)
	}

	return output, nil
}

// Learn 学习经验（简化版本）
func (a *BaseAgent) Learn(experience Experience) error {
	// 简化处理：记录学习事件
	if a.metrics != nil {
		a.metrics.RecordMetric("learn_count", 1.0)
	}

	// TODO: 如需使用学习引擎，实现完整的学习逻辑
	return nil
}

// GetStatus 获取代理状态
func (a *BaseAgent) GetStatus() Status {
	a.mu.RLock()
	defer a.mu.RUnlock()

	metrics := make(map[string]float64)
	if a.metrics != nil {
		m := a.metrics.GetMetrics()
		if m != nil {
			metrics = m
		}
	}

	return Status{
		ID:       a.id,
		State:    a.state,
		Metrics:  metrics,
		LastSeen: time.Now(),
		Load:     a.calculateLoad(),
	}
}

// 内部方法

// initialize 初始化代理组件（简化版本）
func (a *BaseAgent) initialize() error {
	// 简化：组件在外部设置，这里只做基本检查
	// 允许组件为nil，在使用时进行判断
	return nil
}

// cleanup 清理资源
func (a *BaseAgent) cleanup() error {
	// 简化清理逻辑：组件清理由外部管理
	// TODO: 如果接口添加了Cleanup方法，在这里调用
	return nil
}

// execute 执行决策
func (a *BaseAgent) execute(decision Decision) Output {
	// 这里实现具体的决策执行逻辑
	// 可以根据不同的决策类型执行不同的操作
	return Output{
		ID:   decision.ID,
		Type: "result",
		Data: map[string]interface{}{
			"action":     decision.Action,
			"confidence": decision.Confidence,
		},
		Timestamp: time.Now(),
	}
}

// learnFromOutcome 从结果中学习
func (a *BaseAgent) learnFromOutcome(decision Decision, output Output) {
	// 简化的学习逻辑
	reward := 1.0
	if output.Error != nil {
		reward = 0.0
	}

	if a.metrics != nil {
		a.metrics.RecordMetric("outcome_reward", reward)
		a.metrics.RecordMetric("outcome_count", 1.0)
	}
}

// calculateLoad 计算当前负载
func (a *BaseAgent) calculateLoad() float64 {
	// 这里可以实现负载计算逻辑
	// 例如：基于处理队列长度、CPU使用率等
	return 0.0
}

// SetLearningEngine 设置学习引擎
func (a *BaseAgent) SetLearningEngine(learning *LearningEngine) {
	a.learning = learning
}

// SetDecisionEngine 设置决策引擎
func (a *BaseAgent) SetDecisionEngine(decision *DecisionEngine) {
	a.decision = decision
}

// SetMetricsCollector 设置指标收集器
func (a *BaseAgent) SetMetricsCollector(metrics MetricsCollector) {
	a.metrics = metrics
}
