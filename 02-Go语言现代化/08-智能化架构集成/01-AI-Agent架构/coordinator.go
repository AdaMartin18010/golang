package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task 任务定义
type Task struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Priority    int                    `json:"priority"`
	Data        map[string]interface{} `json:"data"`
	Timeout     time.Duration          `json:"timeout"`
	Requirements []string              `json:"requirements"`
	CreatedAt   time.Time              `json:"created_at"`
}

// TaskResult 任务结果
type TaskResult struct {
	TaskID    string    `json:"task_id"`
	AgentID   string    `json:"agent_id"`
	Success   bool      `json:"success"`
	Result    Output    `json:"result"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
}

// Coordinator 代理协调器接口
type Coordinator interface {
	RegisterAgent(agent Agent) error
	UnregisterAgent(agentID string) error
	RouteTask(task Task) (Agent, error)
	MonitorAgents() []Status
	OptimizeDistribution() error
	GetSystemStatus() SystemStatus
}

// SystemStatus 系统状态
type SystemStatus struct {
	TotalAgents    int                    `json:"total_agents"`
	ActiveAgents   int                    `json:"active_agents"`
	TotalTasks     int64                  `json:"total_tasks"`
	CompletedTasks int64                  `json:"completed_tasks"`
	FailedTasks    int64                  `json:"failed_tasks"`
	SystemLoad     float64                `json:"system_load"`
	Metrics        map[string]interface{} `json:"metrics"`
	LastUpdated    time.Time              `json:"last_updated"`
}

// SmartCoordinator 智能协调器实现
type SmartCoordinator struct {
	agents       map[string]Agent
	router       TaskRouter
	balancer     LoadBalancer
	monitor      SystemMonitor
	optimizer    SystemOptimizer
	mu           sync.RWMutex
	taskQueue    chan Task
	resultQueue  chan TaskResult
	ctx          context.Context
	cancel       context.CancelFunc
	metrics      *SystemMetrics
}

// TaskRouter 任务路由器接口
type TaskRouter interface {
	RouteTask(task Task, agents []Agent) (Agent, error)
	GetBestAgent(task Task, agents []Agent) (Agent, error)
}

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
	IsOverloaded(agent Agent) bool
	Redistribute(agent Agent, task Task) (Agent, error)
	GetLoadDistribution() map[string]float64
}

// SystemMonitor 系统监控器接口
type SystemMonitor interface {
	MonitorAgents(agents map[string]Agent) []Status
	GetSystemHealth() SystemHealth
	DetectAnomalies() []Anomaly
}

// SystemOptimizer 系统优化器接口
type SystemOptimizer interface {
	OptimizeAgentDistribution() error
	OptimizeTaskRouting() error
	OptimizeResourceAllocation() error
}

// SystemHealth 系统健康状态
type SystemHealth struct {
	OverallHealth string                 `json:"overall_health"`
	AgentHealth   map[string]string      `json:"agent_health"`
	ResourceUsage map[string]float64     `json:"resource_usage"`
	Alerts        []Alert                `json:"alerts"`
	LastChecked   time.Time              `json:"last_checked"`
}

// Anomaly 异常
type Anomaly struct {
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Description string                 `json:"description"`
	AgentID     string                 `json:"agent_id,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Alert 告警
type Alert struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Message     string                 `json:"message"`
	AgentID     string                 `json:"agent_id,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Acknowledged bool                  `json:"acknowledged"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	TotalTasks     int64             `json:"total_tasks"`
	CompletedTasks int64             `json:"completed_tasks"`
	FailedTasks    int64             `json:"failed_tasks"`
	AvgTaskTime    time.Duration     `json:"avg_task_time"`
	AgentLoads     map[string]float64 `json:"agent_loads"`
	mu             sync.RWMutex
}

// NewSmartCoordinator 创建智能协调器
func NewSmartCoordinator() *SmartCoordinator {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &SmartCoordinator{
		agents:      make(map[string]Agent),
		router:      NewSimpleTaskRouter(),
		balancer:    NewSimpleLoadBalancer(),
		monitor:     NewSimpleSystemMonitor(),
		optimizer:   NewSimpleSystemOptimizer(),
		taskQueue:   make(chan Task, 1000),
		resultQueue: make(chan TaskResult, 1000),
		ctx:         ctx,
		cancel:      cancel,
		metrics:     &SystemMetrics{
			AgentLoads: make(map[string]float64),
		},
	}
}

// RegisterAgent 注册代理
func (c *SmartCoordinator) RegisterAgent(agent Agent) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	agentID := agent.ID()
	if _, exists := c.agents[agentID]; exists {
		return fmt.Errorf("agent %s already registered", agentID)
	}
	
	c.agents[agentID] = agent
	
	// 启动代理
	if err := agent.Start(c.ctx); err != nil {
		return fmt.Errorf("failed to start agent %s: %w", agentID, err)
	}
	
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
	
	return nil
}

// RouteTask 路由任务
func (c *SmartCoordinator) RouteTask(task Task) (Agent, error) {
	c.mu.RLock()
	agents := make([]Agent, 0, len(c.agents))
	for _, agent := range c.agents {
		agents = append(agents, agent)
	}
	c.mu.RUnlock()
	
	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents available")
	}
	
	// 使用路由器选择最佳代理
	selectedAgent, err := c.router.RouteTask(task, agents)
	if err != nil {
		return nil, fmt.Errorf("failed to route task: %w", err)
	}
	
	// 检查负载均衡
	if c.balancer.IsOverloaded(selectedAgent) {
		selectedAgent, err = c.balancer.Redistribute(selectedAgent, task)
		if err != nil {
			return nil, fmt.Errorf("failed to redistribute task: %w", err)
		}
	}
	
	return selectedAgent, nil
}

// MonitorAgents 监控代理
func (c *SmartCoordinator) MonitorAgents() []Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.monitor.MonitorAgents(c.agents)
}

// OptimizeDistribution 优化分布
func (c *SmartCoordinator) OptimizeDistribution() error {
	return c.optimizer.OptimizeAgentDistribution()
}

// GetSystemStatus 获取系统状态
func (c *SmartCoordinator) GetSystemStatus() SystemStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	activeAgents := 0
	for _, agent := range c.agents {
		status := agent.GetStatus()
		if status.State == AgentStateRunning || status.State == AgentStateBusy {
			activeAgents++
		}
	}
	
	c.metrics.mu.RLock()
	totalTasks := c.metrics.TotalTasks
	completedTasks := c.metrics.CompletedTasks
	failedTasks := c.metrics.FailedTasks
	avgTaskTime := c.metrics.AvgTaskTime
	c.metrics.mu.RUnlock()
	
	// 计算系统负载
	systemLoad := 0.0
	if len(c.agents) > 0 {
		totalLoad := 0.0
		for _, agent := range c.agents {
			status := agent.GetStatus()
			totalLoad += status.Load
		}
		systemLoad = totalLoad / float64(len(c.agents))
	}
	
	return SystemStatus{
		TotalAgents:    len(c.agents),
		ActiveAgents:   activeAgents,
		TotalTasks:     totalTasks,
		CompletedTasks: completedTasks,
		FailedTasks:    failedTasks,
		SystemLoad:     systemLoad,
		Metrics: map[string]interface{}{
			"avg_task_time": avgTaskTime.Milliseconds(),
			"agent_loads":   c.metrics.AgentLoads,
		},
		LastUpdated: time.Now(),
	}
}

// ProcessTask 处理任务
func (c *SmartCoordinator) ProcessTask(task Task) (*TaskResult, error) {
	// 路由任务到合适的代理
	agent, err := c.RouteTask(task)
	if err != nil {
		return nil, fmt.Errorf("failed to route task: %w", err)
	}
	
	// 转换任务为输入
	input := Input{
		ID:       task.ID,
		Type:     task.Type,
		Data:     task.Data,
		Priority: task.Priority,
		Timeout:  task.Timeout,
	}
	
	// 处理任务
	start := time.Now()
	output, err := agent.Process(c.ctx, input)
	duration := time.Since(start)
	
	// 创建任务结果
	result := &TaskResult{
		TaskID:    task.ID,
		AgentID:   agent.ID(),
		Success:   err == nil && output.Success,
		Result:    output,
		Duration:  duration,
		Timestamp: time.Now(),
	}
	
	// 更新指标
	c.updateMetrics(result)
	
	return result, err
}

// updateMetrics 更新指标
func (c *SmartCoordinator) updateMetrics(result *TaskResult) {
	c.metrics.mu.Lock()
	defer c.metrics.mu.Unlock()
	
	c.metrics.TotalTasks++
	if result.Success {
		c.metrics.CompletedTasks++
	} else {
		c.metrics.FailedTasks++
	}
	
	// 更新平均任务时间
	if c.metrics.TotalTasks > 0 {
		totalTime := c.metrics.AvgTaskTime * time.Duration(c.metrics.TotalTasks-1)
		c.metrics.AvgTaskTime = (totalTime + result.Duration) / time.Duration(c.metrics.TotalTasks)
	} else {
		c.metrics.AvgTaskTime = result.Duration
	}
	
	// 更新代理负载
	c.metrics.AgentLoads[result.AgentID] = float64(result.Duration.Milliseconds()) / 1000.0
}

// Start 启动协调器
func (c *SmartCoordinator) Start() error {
	// 启动任务处理协程
	go c.taskProcessor()
	
	// 启动结果处理协程
	go c.resultProcessor()
	
	// 启动监控协程
	go c.monitorLoop()
	
	return nil
}

// Stop 停止协调器
func (c *SmartCoordinator) Stop() error {
	c.cancel()
	
	// 停止所有代理
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for _, agent := range c.agents {
		agent.Stop()
	}
	
	return nil
}

// taskProcessor 任务处理器
func (c *SmartCoordinator) taskProcessor() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case task := <-c.taskQueue:
			result, err := c.ProcessTask(task)
			if err != nil {
				// 处理错误
				fmt.Printf("Task processing error: %v\n", err)
			}
			if result != nil {
				c.resultQueue <- *result
			}
		}
	}
}

// resultProcessor 结果处理器
func (c *SmartCoordinator) resultProcessor() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case result := <-c.resultQueue:
			// 处理结果
			fmt.Printf("Task %s completed by agent %s in %v\n", 
				result.TaskID, result.AgentID, result.Duration)
		}
	}
}

// monitorLoop 监控循环
func (c *SmartCoordinator) monitorLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			// 执行系统优化
			c.OptimizeDistribution()
			
			// 检查系统健康状态
			health := c.monitor.GetSystemHealth()
			if health.OverallHealth != "healthy" {
				fmt.Printf("System health warning: %s\n", health.OverallHealth)
			}
		}
	}
}

// SubmitTask 提交任务
func (c *SmartCoordinator) SubmitTask(task Task) error {
	select {
	case c.taskQueue <- task:
		return nil
	case <-c.ctx.Done():
		return fmt.Errorf("coordinator is shutting down")
	default:
		return fmt.Errorf("task queue is full")
	}
}
