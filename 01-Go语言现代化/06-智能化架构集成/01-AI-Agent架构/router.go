package main

import (
	"fmt"
	"math"
	"time"
)

// SimpleTaskRouter 简单任务路由器实现
type SimpleTaskRouter struct {
	routingRules map[string]string
}

// NewSimpleTaskRouter 创建简单任务路由器
func NewSimpleTaskRouter() *SimpleTaskRouter {
	return &SimpleTaskRouter{
		routingRules: map[string]string{
			"data_processing": "data_agent",
			"decision_making": "decision_agent",
			"collaboration":   "collaboration_agent",
			"monitoring":      "monitoring_agent",
		},
	}
}

// RouteTask 路由任务
func (r *SimpleTaskRouter) RouteTask(task Task, agents []Agent) (Agent, error) {
	// 首先尝试基于任务类型路由
	if agentType, exists := r.routingRules[task.Type]; exists {
		for _, agent := range agents {
			if agent.ID() == agentType {
				return agent, nil
			}
		}
	}

	// 如果没有找到特定类型的代理，使用最佳代理选择
	return r.GetBestAgent(task, agents)
}

// GetBestAgent 获取最佳代理
func (r *SimpleTaskRouter) GetBestAgent(task Task, agents []Agent) (Agent, error) {
	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents available")
	}

	// 基于多个因素评分选择最佳代理
	bestAgent := agents[0]
	bestScore := r.calculateAgentScore(task, agents[0])

	for _, agent := range agents[1:] {
		score := r.calculateAgentScore(task, agent)
		if score > bestScore {
			bestScore = score
			bestAgent = agent
		}
	}

	return bestAgent, nil
}

// calculateAgentScore 计算代理评分
func (r *SimpleTaskRouter) calculateAgentScore(task Task, agent Agent) float64 {
	status := agent.GetStatus()
	capabilities := agent.GetCapabilities()

	score := 0.0

	// 1. 健康状态评分 (40%)
	if status.State == AgentStateRunning {
		score += 0.4
	} else if status.State == AgentStateBusy {
		score += 0.2
	}

	// 2. 负载评分 (30%) - 负载越低评分越高
	loadScore := math.Max(0, 1.0-status.Load)
	score += loadScore * 0.3

	// 3. 能力匹配评分 (20%)
	capabilityScore := r.calculateCapabilityScore(task, capabilities)
	score += capabilityScore * 0.2

	// 4. 历史性能评分 (10%)
	performanceScore := r.calculatePerformanceScore(status)
	score += performanceScore * 0.1

	return score
}

// calculateCapabilityScore 计算能力匹配评分
func (r *SimpleTaskRouter) calculateCapabilityScore(task Task, capabilities []string) float64 {
	if len(capabilities) == 0 {
		return 0.5 // 默认评分
	}

	// 检查任务要求是否匹配代理能力
	matches := 0
	for _, requirement := range task.Requirements {
		for _, capability := range capabilities {
			if requirement == capability {
				matches++
				break
			}
		}
	}

	if len(task.Requirements) == 0 {
		return 0.5 // 没有特定要求时给默认评分
	}

	return float64(matches) / float64(len(task.Requirements))
}

// calculatePerformanceScore 计算性能评分
func (r *SimpleTaskRouter) calculatePerformanceScore(status Status) float64 {
	// 基于错误率和处理数量计算性能评分
	errorRate := 0.0
	if status.Processed > 0 {
		errorRate = float64(status.Errors) / float64(status.Processed)
	}

	// 错误率越低，性能评分越高
	performanceScore := math.Max(0, 1.0-errorRate)

	return performanceScore
}

// SimpleLoadBalancer 简单负载均衡器实现
type SimpleLoadBalancer struct {
	loadThreshold float64
}

// NewSimpleLoadBalancer 创建简单负载均衡器
func NewSimpleLoadBalancer() *SimpleLoadBalancer {
	return &SimpleLoadBalancer{
		loadThreshold: 0.8, // 80%负载阈值
	}
}

// IsOverloaded 检查代理是否过载
func (l *SimpleLoadBalancer) IsOverloaded(agent Agent) bool {
	status := agent.GetStatus()
	return status.Load > l.loadThreshold
}

// Redistribute 重新分配任务
func (l *SimpleLoadBalancer) Redistribute(overloadedAgent Agent, task Task) (Agent, error) {
	// 这里需要访问协调器的代理列表，简化实现
	// 在实际实现中，应该通过协调器获取其他代理
	return nil, fmt.Errorf("redistribution not implemented in simple load balancer")
}

// GetLoadDistribution 获取负载分布
func (l *SimpleLoadBalancer) GetLoadDistribution() map[string]float64 {
	// 简化实现，返回空映射
	return make(map[string]float64)
}

// AdvancedLoadBalancer 高级负载均衡器实现
type AdvancedLoadBalancer struct {
	loadThreshold     float64
	agentCapabilities map[string][]string
	loadHistory       map[string][]float64
	maxHistorySize    int
}

// NewAdvancedLoadBalancer 创建高级负载均衡器
func NewAdvancedLoadBalancer() *AdvancedLoadBalancer {
	return &AdvancedLoadBalancer{
		loadThreshold:     0.8,
		agentCapabilities: make(map[string][]string),
		loadHistory:       make(map[string][]float64),
		maxHistorySize:    100,
	}
}

// IsOverloaded 检查代理是否过载
func (a *AdvancedLoadBalancer) IsOverloaded(agent Agent) bool {
	status := agent.GetStatus()

	// 检查当前负载
	if status.Load > a.loadThreshold {
		return true
	}

	// 检查负载趋势
	agentID := agent.ID()
	if history, exists := a.loadHistory[agentID]; exists && len(history) > 5 {
		// 如果最近5次的平均负载超过阈值，认为过载
		recentLoad := 0.0
		for i := len(history) - 5; i < len(history); i++ {
			recentLoad += history[i]
		}
		recentLoad /= 5.0

		if recentLoad > a.loadThreshold {
			return true
		}
	}

	return false
}

// Redistribute 重新分配任务
func (a *AdvancedLoadBalancer) Redistribute(overloadedAgent Agent, task Task) (Agent, error) {
	// 在实际实现中，这里应该从协调器获取所有可用代理
	// 然后选择最适合的代理进行重新分配
	return nil, fmt.Errorf("redistribution requires access to all agents")
}

// GetLoadDistribution 获取负载分布
func (a *AdvancedLoadBalancer) GetLoadDistribution() map[string]float64 {
	distribution := make(map[string]float64)

	for agentID, history := range a.loadHistory {
		if len(history) > 0 {
			// 计算平均负载
			totalLoad := 0.0
			for _, load := range history {
				totalLoad += load
			}
			distribution[agentID] = totalLoad / float64(len(history))
		}
	}

	return distribution
}

// UpdateLoadHistory 更新负载历史
func (a *AdvancedLoadBalancer) UpdateLoadHistory(agentID string, load float64) {
	if a.loadHistory[agentID] == nil {
		a.loadHistory[agentID] = make([]float64, 0)
	}

	history := a.loadHistory[agentID]
	history = append(history, load)

	// 保持历史记录在合理范围内
	if len(history) > a.maxHistorySize {
		history = history[len(history)-a.maxHistorySize:]
	}

	a.loadHistory[agentID] = history
}

// GetLoadTrend 获取负载趋势
func (a *AdvancedLoadBalancer) GetLoadTrend(agentID string) string {
	history, exists := a.loadHistory[agentID]
	if !exists || len(history) < 3 {
		return "stable"
	}

	// 计算最近3次的趋势
	recent := history[len(history)-3:]

	// 简单趋势分析
	if recent[2] > recent[1] && recent[1] > recent[0] {
		return "increasing"
	} else if recent[2] < recent[1] && recent[1] < recent[0] {
		return "decreasing"
	}

	return "stable"
}

// SimpleSystemMonitor 简单系统监控器实现
type SimpleSystemMonitor struct {
	healthThresholds map[string]float64
}

// NewSimpleSystemMonitor 创建简单系统监控器
func NewSimpleSystemMonitor() *SimpleSystemMonitor {
	return &SimpleSystemMonitor{
		healthThresholds: map[string]float64{
			"load_threshold":          0.8,
			"error_rate_threshold":    0.1,
			"response_time_threshold": 1000, // 毫秒
		},
	}
}

// MonitorAgents 监控代理
func (m *SimpleSystemMonitor) MonitorAgents(agents map[string]Agent) []Status {
	statuses := make([]Status, 0, len(agents))

	for _, agent := range agents {
		status := agent.GetStatus()
		statuses = append(statuses, status)
	}

	return statuses
}

// GetSystemHealth 获取系统健康状态
func (m *SimpleSystemMonitor) GetSystemHealth() SystemHealth {
	// 简化实现，返回基本健康状态
	return SystemHealth{
		OverallHealth: "healthy",
		AgentHealth:   make(map[string]string),
		ResourceUsage: make(map[string]float64),
		Alerts:        make([]Alert, 0),
		LastChecked:   time.Now(),
	}
}

// DetectAnomalies 检测异常
func (m *SimpleSystemMonitor) DetectAnomalies() []Anomaly {
	// 简化实现，返回空异常列表
	return make([]Anomaly, 0)
}

// SimpleSystemOptimizer 简单系统优化器实现
type SimpleSystemOptimizer struct {
	optimizationRules []OptimizationRule
}

// OptimizationRule 优化规则
type OptimizationRule struct {
	Name       string                 `json:"name"`
	Condition  string                 `json:"condition"`
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters"`
	Enabled    bool                   `json:"enabled"`
}

// NewSimpleSystemOptimizer 创建简单系统优化器
func NewSimpleSystemOptimizer() *SimpleSystemOptimizer {
	return &SimpleSystemOptimizer{
		optimizationRules: []OptimizationRule{
			{
				Name:      "load_balancing",
				Condition: "high_load",
				Action:    "redistribute_tasks",
				Enabled:   true,
			},
			{
				Name:      "error_recovery",
				Condition: "high_error_rate",
				Action:    "restart_agent",
				Enabled:   true,
			},
		},
	}
}

// OptimizeAgentDistribution 优化代理分布
func (o *SimpleSystemOptimizer) OptimizeAgentDistribution() error {
	// 简化实现，记录优化操作
	fmt.Println("Optimizing agent distribution...")
	return nil
}

// OptimizeTaskRouting 优化任务路由
func (o *SimpleSystemOptimizer) OptimizeTaskRouting() error {
	// 简化实现，记录优化操作
	fmt.Println("Optimizing task routing...")
	return nil
}

// OptimizeResourceAllocation 优化资源分配
func (o *SimpleSystemOptimizer) OptimizeResourceAllocation() error {
	// 简化实现，记录优化操作
	fmt.Println("Optimizing resource allocation...")
	return nil
}

// AddOptimizationRule 添加优化规则
func (o *SimpleSystemOptimizer) AddOptimizationRule(rule OptimizationRule) {
	o.optimizationRules = append(o.optimizationRules, rule)
}

// GetOptimizationRules 获取优化规则
func (o *SimpleSystemOptimizer) GetOptimizationRules() []OptimizationRule {
	return o.optimizationRules
}

// EnableOptimizationRule 启用优化规则
func (o *SimpleSystemOptimizer) EnableOptimizationRule(ruleName string) error {
	for i, rule := range o.optimizationRules {
		if rule.Name == ruleName {
			o.optimizationRules[i].Enabled = true
			return nil
		}
	}
	return fmt.Errorf("optimization rule %s not found", ruleName)
}

// DisableOptimizationRule 禁用优化规则
func (o *SimpleSystemOptimizer) DisableOptimizationRule(ruleName string) error {
	for i, rule := range o.optimizationRules {
		if rule.Name == ruleName {
			o.optimizationRules[i].Enabled = false
			return nil
		}
	}
	return fmt.Errorf("optimization rule %s not found", ruleName)
}
