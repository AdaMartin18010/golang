package coordination

import (
	"context"
	"fmt"
	"sync"
	"time"

	"core"
)

// Task 任务定义
type Task struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Priority     int                    `json:"priority"`
	Input        core.Input             `json:"input"`
	Requirements map[string]interface{} `json:"requirements"`
	CreatedAt    time.Time              `json:"created_at"`
}

// TaskProfile 任务特征
type TaskProfile struct {
	Complexity     float64            `json:"complexity"`
	ResourceNeeds  map[string]float64 `json:"resource_needs"`
	Specialization string             `json:"specialization"`
	EstimatedTime  time.Duration      `json:"estimated_time"`
}

// Coordinator 代理协调器接口
type Coordinator interface {
	RegisterAgent(agent core.Agent) error
	UnregisterAgent(agentID string) error
	RouteTask(task Task) (core.Agent, error)
	MonitorAgents() []core.Status
	OptimizeDistribution() error
	GetSystemStatus() SystemStatus
}

// SmartCoordinator 智能协调器实现
type SmartCoordinator struct {
	agents    map[string]core.Agent
	router    TaskRouter
	balancer  LoadBalancer
	monitor   SystemMonitor
	optimizer SystemOptimizer
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewSmartCoordinator 创建新的智能协调器
func NewSmartCoordinator() *SmartCoordinator {
	ctx, cancel := context.WithCancel(context.Background())
	return &SmartCoordinator{
		agents:    make(map[string]core.Agent),
		router:    NewDefaultTaskRouter(),
		balancer:  NewDefaultLoadBalancer(),
		monitor:   NewDefaultSystemMonitor(),
		optimizer: NewDefaultSystemOptimizer(),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// RegisterAgent 注册代理
func (c *SmartCoordinator) RegisterAgent(agent core.Agent) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	agentID := agent.ID()
	if _, exists := c.agents[agentID]; exists {
		return fmt.Errorf("agent %s already registered", agentID)
	}

	// 启动代理
	if err := agent.Start(c.ctx); err != nil {
		return fmt.Errorf("failed to start agent %s: %w", agentID, err)
	}

	c.agents[agentID] = agent

	// 通知监控器
	c.monitor.AgentRegistered(agentID)

	// 通知负载均衡器
	c.balancer.AgentAdded(agentID)

	return nil
}

// UnregisterAgent 注销代理
func (c *SmartCoordinator) UnregisterAgent(agentID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	agent, exists := c.agents[agentID]
	if !exists {
		return fmt.Errorf("agent %s not found", agentID)
	}

	// 停止代理
	if err := agent.Stop(); err != nil {
		return fmt.Errorf("failed to stop agent %s: %w", agentID, err)
	}

	delete(c.agents, agentID)

	// 通知监控器
	c.monitor.AgentUnregistered(agentID)

	// 通知负载均衡器
	c.balancer.AgentRemoved(agentID)

	return nil
}

// RouteTask 路由任务到合适的代理
func (c *SmartCoordinator) RouteTask(task Task) (core.Agent, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.agents) == 0 {
		return nil, fmt.Errorf("no agents available")
	}

	// 分析任务特征
	taskProfile := c.analyzeTask(task)

	// 选择最适合的代理
	agent := c.selectBestAgent(taskProfile)

	// 负载均衡检查
	if c.balancer.IsOverloaded(agent) {
		agent = c.balancer.Redistribute(agent, task)
	}

	// 记录任务分配
	c.monitor.TaskAssigned(task.ID, agent.ID())

	return agent, nil
}

// MonitorAgents 监控所有代理状态
func (c *SmartCoordinator) MonitorAgents() []core.Status {
	c.mu.RLock()
	defer c.mu.RUnlock()

	statuses := make([]core.Status, 0, len(c.agents))
	for _, agent := range c.agents {
		status := agent.GetStatus()
		statuses = append(statuses, status)
	}

	return statuses
}

// OptimizeDistribution 优化任务分布
func (c *SmartCoordinator) OptimizeDistribution() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 获取系统状态
	systemStatus := c.monitor.GetSystemStatus()

	// 分析性能瓶颈
	bottlenecks := c.analyzer.AnalyzeBottlenecks(systemStatus)

	// 生成优化建议
	recommendations := c.optimizer.GenerateRecommendations(bottlenecks)

	// 应用优化
	return c.optimizer.ApplyOptimizations(recommendations)
}

// GetSystemStatus 获取系统状态
func (c *SmartCoordinator) GetSystemStatus() SystemStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	agentStatuses := make([]core.Status, 0, len(c.agents))
	for _, agent := range c.agents {
		status := agent.GetStatus()
		agentStatuses = append(agentStatuses, status)
	}

	return SystemStatus{
		TotalAgents:   len(c.agents),
		ActiveAgents:  c.countActiveAgents(agentStatuses),
		AgentStatuses: agentStatuses,
		SystemLoad:    c.calculateSystemLoad(agentStatuses),
		LastUpdated:   time.Now(),
	}
}

// 内部方法

// analyzeTask 分析任务特征
func (c *SmartCoordinator) analyzeTask(task Task) TaskProfile {
	// 这里实现任务特征分析逻辑
	// 可以根据任务类型、复杂度、资源需求等进行分析
	return TaskProfile{
		Complexity:     c.calculateComplexity(task),
		ResourceNeeds:  c.estimateResourceNeeds(task),
		Specialization: c.determineSpecialization(task),
		EstimatedTime:  c.estimateProcessingTime(task),
	}
}

// selectBestAgent 选择最适合的代理
func (c *SmartCoordinator) selectBestAgent(profile TaskProfile) core.Agent {
	var bestAgent core.Agent
	var bestScore float64

	for _, agent := range c.agents {
		score := c.calculateAgentScore(agent, profile)
		if score > bestScore {
			bestScore = score
			bestAgent = agent
		}
	}

	return bestAgent
}

// calculateAgentScore 计算代理得分
func (c *SmartCoordinator) calculateAgentScore(agent core.Agent, profile TaskProfile) float64 {
	status := agent.GetStatus()

	// 基础得分：基于代理状态
	score := 1.0

	// 负载因子：负载越低得分越高
	loadFactor := 1.0 - status.Load
	score *= loadFactor

	// 专业匹配因子：根据任务专业需求匹配
	specializationFactor := c.calculateSpecializationMatch(agent, profile.Specialization)
	score *= specializationFactor

	// 性能因子：基于历史性能指标
	performanceFactor := c.calculatePerformanceFactor(agent)
	score *= performanceFactor

	return score
}

// calculateComplexity 计算任务复杂度
func (c *SmartCoordinator) calculateComplexity(task Task) float64 {
	// 这里实现任务复杂度计算逻辑
	// 可以根据任务类型、数据量、处理要求等计算
	return 0.5 // 示例值
}

// estimateResourceNeeds 估算资源需求
func (c *SmartCoordinator) estimateResourceNeeds(task Task) map[string]float64 {
	// 这里实现资源需求估算逻辑
	return map[string]float64{
		"cpu":     0.1,
		"memory":  0.2,
		"network": 0.05,
	}
}

// determineSpecialization 确定任务专业需求
func (c *SmartCoordinator) determineSpecialization(task Task) string {
	// 这里实现专业需求确定逻辑
	// 可以根据任务类型确定所需的专业能力
	switch task.Type {
	case "data_processing":
		return "data_processing"
	case "decision_making":
		return "decision_making"
	case "collaboration":
		return "collaboration"
	default:
		return "general"
	}
}

// estimateProcessingTime 估算处理时间
func (c *SmartCoordinator) estimateProcessingTime(task Task) time.Duration {
	// 这里实现处理时间估算逻辑
	// 可以根据任务复杂度和历史数据估算
	return 100 * time.Millisecond
}

// calculateSpecializationMatch 计算专业匹配度
func (c *SmartCoordinator) calculateSpecializationMatch(agent core.Agent, specialization string) float64 {
	// 这里实现专业匹配度计算逻辑
	// 可以根据代理的能力和任务需求计算匹配度
	return 0.8 // 示例值
}

// calculatePerformanceFactor 计算性能因子
func (c *SmartCoordinator) calculatePerformanceFactor(agent core.Agent) float64 {
	// 这里实现性能因子计算逻辑
	// 可以根据代理的历史性能指标计算
	return 0.9 // 示例值
}

// countActiveAgents 统计活跃代理数量
func (c *SmartCoordinator) countActiveAgents(statuses []core.Status) int {
	count := 0
	for _, status := range statuses {
		if status.State == core.AgentStateRunning {
			count++
		}
	}
	return count
}

// calculateSystemLoad 计算系统负载
func (c *SmartCoordinator) calculateSystemLoad(statuses []core.Status) float64 {
	if len(statuses) == 0 {
		return 0.0
	}

	totalLoad := 0.0
	for _, status := range statuses {
		totalLoad += status.Load
	}

	return totalLoad / float64(len(statuses))
}

// SetTaskRouter 设置任务路由器
func (c *SmartCoordinator) SetTaskRouter(router TaskRouter) {
	c.router = router
}

// SetLoadBalancer 设置负载均衡器
func (c *SmartCoordinator) SetLoadBalancer(balancer LoadBalancer) {
	c.balancer = balancer
}

// SetSystemMonitor 设置系统监控器
func (c *SmartCoordinator) SetSystemMonitor(monitor SystemMonitor) {
	c.monitor = monitor
}

// SetSystemOptimizer 设置系统优化器
func (c *SmartCoordinator) SetSystemOptimizer(optimizer SystemOptimizer) {
	c.optimizer = optimizer
}

// Stop 停止协调器
func (c *SmartCoordinator) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 停止所有代理
	for _, agent := range c.agents {
		agent.Stop()
	}

	// 取消上下文
	if c.cancel != nil {
		c.cancel()
	}

	return nil
}
