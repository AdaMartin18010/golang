package core

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// AgentState 代理状态枚举
type AgentState int

const (
	AgentStateIdle AgentState = iota
	AgentStateRunning
	AgentStateStopped
	AgentStateError
)

// String 返回状态字符串表示
func (s AgentState) String() string {
	switch s {
	case AgentStateIdle:
		return "idle"
	case AgentStateRunning:
		return "running"
	case AgentStateStopped:
		return "stopped"
	case AgentStateError:
		return "error"
	default:
		return "unknown"
	}
}

// Input 代理输入数据
type Input struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      interface{}            `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

// Output 代理输出数据
type Output struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      interface{}            `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	Error     error                  `json:"error,omitempty"`
}

// Status 代理状态信息
type Status struct {
	ID       string             `json:"id"`
	State    AgentState         `json:"state"`
	Metrics  map[string]float64 `json:"metrics"`
	LastSeen time.Time          `json:"last_seen"`
	Load     float64            `json:"load"`
	Error    string             `json:"error,omitempty"`
}

// Experience 学习经验
type Experience struct {
	Input     Input     `json:"input"`
	Output    Output    `json:"output"`
	Reward    float64   `json:"reward"`
	Timestamp time.Time `json:"timestamp"`
}

// Agent 智能代理接口
type Agent interface {
	ID() string
	Start(ctx context.Context) error
	Stop() error
	Process(input Input) (Output, error)
	Learn(experience Experience) error
	GetStatus() Status
}

// BaseAgent 基础代理实现
type BaseAgent struct {
	id       string
	state    AgentState
	config   AgentConfig
	learning LearningEngine
	decision DecisionEngine
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
	a.metrics.RecordEvent("agent_started", map[string]interface{}{
		"agent_id":  a.id,
		"timestamp": time.Now(),
	})

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
	a.metrics.RecordEvent("agent_stopped", map[string]interface{}{
		"agent_id":  a.id,
		"timestamp": time.Now(),
	})

	return nil
}

// Process 处理输入数据
func (a *BaseAgent) Process(input Input) (Output, error) {
	a.mu.RLock()
	if a.state != AgentStateRunning {
		a.mu.RUnlock()
		return Output{}, fmt.Errorf("agent %s is not running", a.id)
	}
	a.mu.RUnlock()

	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		a.metrics.RecordMetric("processing_time", duration.Seconds())
	}()

	// 1. 感知输入
	context := a.perceive(input)

	// 2. 学习历史经验
	a.learnFromHistory(context)

	// 3. 做出决策
	decision, err := a.decision.MakeDecision(context)
	if err != nil {
		a.metrics.RecordEvent("decision_error", map[string]interface{}{
			"agent_id": a.id,
			"error":    err.Error(),
		})
		return Output{}, fmt.Errorf("failed to make decision: %w", err)
	}

	// 4. 执行决策
	output := a.execute(decision)

	// 5. 学习结果
	a.learnFromOutcome(decision, output)

	// 记录成功处理
	a.metrics.RecordEvent("processing_success", map[string]interface{}{
		"agent_id": a.id,
		"input_id": input.ID,
	})

	return output, nil
}

// Learn 学习经验
func (a *BaseAgent) Learn(experience Experience) error {
	if a.learning == nil {
		return fmt.Errorf("learning engine not initialized")
	}

	return a.learning.UpdateModel(experience)
}

// GetStatus 获取代理状态
func (a *BaseAgent) GetStatus() Status {
	a.mu.RLock()
	defer a.mu.RUnlock()

	metrics := a.metrics.GetMetrics()
	if metrics == nil {
		metrics = make(map[string]float64)
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

// initialize 初始化代理组件
func (a *BaseAgent) initialize() error {
	// 初始化学习引擎
	if a.learning == nil {
		a.learning = NewDefaultLearningEngine()
	}

	// 初始化决策引擎
	if a.decision == nil {
		a.decision = NewDefaultDecisionEngine()
	}

	// 初始化指标收集器
	if a.metrics == nil {
		a.metrics = NewDefaultMetricsCollector()
	}

	return nil
}

// cleanup 清理资源
func (a *BaseAgent) cleanup() error {
	// 清理学习引擎
	if a.learning != nil {
		if err := a.learning.Cleanup(); err != nil {
			return err
		}
	}

	// 清理决策引擎
	if a.decision != nil {
		if err := a.decision.Cleanup(); err != nil {
			return err
		}
	}

	// 清理指标收集器
	if a.metrics != nil {
		if err := a.metrics.Cleanup(); err != nil {
			return err
		}
	}

	return nil
}

// perceive 感知输入
func (a *BaseAgent) perceive(input Input) Context {
	return Context{
		Input:     input,
		AgentID:   a.id,
		Timestamp: time.Now(),
	}
}

// learnFromHistory 从历史经验学习
func (a *BaseAgent) learnFromHistory(context Context) {
	// 这里可以实现从历史数据中学习的逻辑
	// 例如：分析历史决策的成功率，调整策略等
}

// execute 执行决策
func (a *BaseAgent) execute(decision Decision) Output {
	// 这里实现具体的决策执行逻辑
	// 可以根据不同的决策类型执行不同的操作
	return Output{
		ID:        decision.ID,
		Type:      decision.Type,
		Data:      decision.Data,
		Timestamp: time.Now(),
	}
}

// learnFromOutcome 从结果中学习
func (a *BaseAgent) learnFromOutcome(decision Decision, output Output) {
	// 这里可以实现从执行结果中学习的逻辑
	// 例如：评估决策效果，更新模型等
}

// calculateLoad 计算当前负载
func (a *BaseAgent) calculateLoad() float64 {
	// 这里可以实现负载计算逻辑
	// 例如：基于处理队列长度、CPU使用率等
	return 0.0
}

// SetLearningEngine 设置学习引擎
func (a *BaseAgent) SetLearningEngine(learning LearningEngine) {
	a.learning = learning
}

// SetDecisionEngine 设置决策引擎
func (a *BaseAgent) SetDecisionEngine(decision DecisionEngine) {
	a.decision = decision
}

// SetMetricsCollector 设置指标收集器
func (a *BaseAgent) SetMetricsCollector(metrics MetricsCollector) {
	a.metrics = metrics
}
