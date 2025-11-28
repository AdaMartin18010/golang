package core

import (
	"context"
	"time"
)

// =============================================================================
// 基础配置类型
// =============================================================================

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

// =============================================================================
// 基础数据类型
// =============================================================================

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

// =============================================================================
// 状态类型
// =============================================================================

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

// Status 代理状态信息
type Status struct {
	ID         string             `json:"id"`
	State      AgentState         `json:"state"`
	Metrics    map[string]float64 `json:"metrics"`
	LastSeen   time.Time          `json:"last_seen"`
	Load       float64            `json:"load"`
	Health     string             `json:"health"`
	Experience int                `json:"experience"`
	Error      string             `json:"error,omitempty"`
}

// =============================================================================
// 基础接口类型 (Agent接口，其他接口在各自文件中定义)
// =============================================================================

// Agent 智能代理接口
type Agent interface {
	ID() string
	Start(ctx context.Context) error
	Stop() error
	Process(input Input) (Output, error)
	Learn(experience Experience) error
	GetStatus() Status
}

// Experience 学习经验 (跨模块使用的基础定义)
type Experience struct {
	Input     Input     `json:"input"`
	Output    Output    `json:"output"`
	Reward    float64   `json:"reward"`
	Timestamp time.Time `json:"timestamp"`
}

// =============================================================================
// 指标收集器接口
// =============================================================================

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
	RecordProcess(duration time.Duration, success bool)
	RecordEvent(event string, value float64)
	RecordMetric(name string, value float64)
	GetMetrics() map[string]float64
	Reset()
}

// =============================================================================
// 工具类型
// =============================================================================

// Context 上下文别名
type Context = context.Context
