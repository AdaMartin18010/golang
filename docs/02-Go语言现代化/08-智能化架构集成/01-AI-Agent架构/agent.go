package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// AgentState 代理状态枚举
type AgentState string

const (
	AgentStateIdle    AgentState = "idle"
	AgentStateRunning AgentState = "running"
	AgentStateBusy    AgentState = "busy"
	AgentStateError   AgentState = "error"
	AgentStateStopped AgentState = "stopped"
)

// Input 代理输入
type Input struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Data     map[string]interface{} `json:"data"`
	Priority int                    `json:"priority"`
	Timeout  time.Duration          `json:"timeout"`
}

// Output 代理输出
type Output struct {
	ID       string                 `json:"id"`
	Success  bool                   `json:"success"`
	Data     map[string]interface{} `json:"data"`
	Error    string                 `json:"error,omitempty"`
	Duration time.Duration          `json:"duration"`
}

// Experience 学习经验
type Experience struct {
	Input     Input     `json:"input"`
	Output    Output    `json:"output"`
	Feedback  float64   `json:"feedback"`
	Timestamp time.Time `json:"timestamp"`
}

// Status 代理状态
type Status struct {
	State     AgentState         `json:"state"`
	Metrics   map[string]float64 `json:"metrics"`
	LastSeen  time.Time          `json:"last_seen"`
	Load      float64            `json:"load"`
	Processed int64              `json:"processed"`
	Errors    int64              `json:"errors"`
}

// Agent 智能代理接口
type Agent interface {
	ID() string
	Start(ctx context.Context) error
	Stop() error
	Process(ctx context.Context, input Input) (Output, error)
	Learn(experience Experience) error
	GetStatus() Status
	GetCapabilities() []string
}

// BaseAgent 基础代理实现
type BaseAgent struct {
	id          string
	state       AgentState
	config      AgentConfig
	learning    LearningEngine
	decision    DecisionEngine
	metrics     MetricsCollector
	experiences []Experience
	LastSeen    time.Time
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// AgentConfig 代理配置
type AgentConfig struct {
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	MaxLoad      float64                `json:"max_load"`
	Timeout      time.Duration          `json:"timeout"`
	Retries      int                    `json:"retries"`
	Capabilities []string               `json:"capabilities"`
	Parameters   map[string]interface{} `json:"parameters"`
}

// NewBaseAgent 创建基础代理
func NewBaseAgent(id string, config AgentConfig) *BaseAgent {
	ctx, cancel := context.WithCancel(context.Background())

	return &BaseAgent{
		id:          id,
		state:       AgentStateIdle,
		config:      config,
		experiences: make([]Experience, 0),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// ID 获取代理ID
func (a *BaseAgent) ID() string {
	return a.id
}

// Start 启动代理
func (a *BaseAgent) Start(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.state != AgentStateIdle {
		return fmt.Errorf("agent %s is not in idle state", a.id)
	}

	a.state = AgentStateRunning
	a.LastSeen = time.Now()

	// 启动监控协程
	go a.monitor()

	return nil
}

// Stop 停止代理
func (a *BaseAgent) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.state == AgentStateStopped {
		return nil
	}

	a.cancel()
	a.state = AgentStateStopped

	return nil
}

// Process 处理输入
func (a *BaseAgent) Process(ctx context.Context, input Input) (Output, error) {
	start := time.Now()

	a.mu.Lock()
	if a.state != AgentStateRunning {
		a.mu.Unlock()
		return Output{}, fmt.Errorf("agent %s is not running", a.id)
	}
	a.state = AgentStateBusy
	a.mu.Unlock()

	defer func() {
		a.mu.Lock()
		a.state = AgentStateRunning
		a.LastSeen = time.Now()
		a.mu.Unlock()
	}()

	// 检查负载
	if a.getCurrentLoad() > a.config.MaxLoad {
		return Output{}, fmt.Errorf("agent %s is overloaded", a.id)
	}

	// 处理输入
	output, err := a.processInput(ctx, input)

	duration := time.Since(start)

	// 记录经验
	experience := Experience{
		Input:     input,
		Output:    output,
		Feedback:  a.calculateFeedback(output, err),
		Timestamp: time.Now(),
	}

	a.Learn(experience)

	// 更新指标
	a.updateMetrics(duration, err == nil)

	return output, nil
}

// Learn 学习经验
func (a *BaseAgent) Learn(experience Experience) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.experiences = append(a.experiences, experience)

	// 保持经验数量在合理范围内
	if len(a.experiences) > 1000 {
		a.experiences = a.experiences[len(a.experiences)-1000:]
	}

	// 如果有学习引擎，进行学习
	if a.learning != nil {
		return a.learning.UpdateModel(experience)
	}

	return nil
}

// GetStatus 获取状态
func (a *BaseAgent) GetStatus() Status {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return Status{
		State:     a.state,
		Metrics:   a.getMetrics(),
		LastSeen:  a.LastSeen,
		Load:      a.getCurrentLoad(),
		Processed: a.getProcessedCount(),
		Errors:    a.getErrorCount(),
	}
}

// GetCapabilities 获取能力
func (a *BaseAgent) GetCapabilities() []string {
	return a.config.Capabilities
}

// processInput 处理输入的核心逻辑
func (a *BaseAgent) processInput(ctx context.Context, input Input) (Output, error) {
	// 基础实现：简单的回显
	output := Output{
		ID:      input.ID,
		Success: true,
		Data: map[string]interface{}{
			"processed_by": a.id,
			"input_type":   input.Type,
			"timestamp":    time.Now(),
		},
	}

	// 如果有决策引擎，使用决策引擎
	if a.decision != nil {
		decision, err := a.decision.MakeDecision(ctx, input)
		if err != nil {
			output.Success = false
			output.Error = err.Error()
			return output, err
		}
		output.Data["decision"] = decision
	}

	return output, nil
}

// calculateFeedback 计算反馈分数
func (a *BaseAgent) calculateFeedback(output Output, err error) float64 {
	if err != nil {
		return 0.0
	}
	if output.Success {
		return 1.0
	}
	return 0.5
}

// getCurrentLoad 获取当前负载
func (a *BaseAgent) getCurrentLoad() float64 {
	// 简单的负载计算：基于处理时间和错误率
	metrics := a.getMetrics()
	avgProcessTime := metrics["avg_process_time"]
	errorRate := metrics["error_rate"]

	// 负载 = 处理时间权重 + 错误率权重
	load := (avgProcessTime/1000.0)*0.7 + errorRate*0.3
	if load > 1.0 {
		load = 1.0
	}
	return load
}

// getMetrics 获取指标
func (a *BaseAgent) getMetrics() map[string]float64 {
	if a.metrics != nil {
		return a.metrics.GetMetrics()
	}

	// 默认指标
	return map[string]float64{
		"avg_process_time": 100.0,
		"error_rate":       0.05,
		"throughput":       10.0,
	}
}

// getProcessedCount 获取处理数量
func (a *BaseAgent) getProcessedCount() int64 {
	return int64(len(a.experiences))
}

// getErrorCount 获取错误数量
func (a *BaseAgent) getErrorCount() int64 {
	count := int64(0)
	for _, exp := range a.experiences {
		if !exp.Output.Success {
			count++
		}
	}
	return count
}

// updateMetrics 更新指标
func (a *BaseAgent) updateMetrics(duration time.Duration, success bool) {
	if a.metrics != nil {
		a.metrics.RecordProcess(duration, success)
	}
}

// monitor 监控协程
func (a *BaseAgent) monitor() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.mu.Lock()
			a.LastSeen = time.Now()
			a.mu.Unlock()
		}
	}
}

// BaseAgent结构体已在上面定义，这里不需要重复定义
