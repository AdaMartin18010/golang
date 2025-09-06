package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Println("=== Go语言AI-Agent架构演示 ===")
	fmt.Println()

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 演示基础代理功能
	demoBasicAgent(ctx)

	// 演示多代理协作
	demoMultiAgentCollaboration(ctx)

	// 演示智能协调器
	demoSmartCoordinator(ctx)

	// 演示专业代理类型
	demoSpecializedAgents(ctx)

	fmt.Println("\n=== 演示完成 ===")
}

// demoBasicAgent 演示基础代理功能
func demoBasicAgent(ctx context.Context) {
	fmt.Println("1. 基础代理功能演示")
	fmt.Println("-------------------")

	// 创建基础代理
	config := AgentConfig{
		Name:         "基础代理",
		Type:         "basic",
		MaxLoad:      0.8,
		Timeout:      5 * time.Second,
		Retries:      3,
		Capabilities: []string{"basic_processing"},
		Parameters:   map[string]interface{}{"version": "1.0"},
	}

	agent := NewBaseAgent("basic-agent-1", config)
	agent.learning = NewSimpleLearningEngine()
	agent.decision = NewSimpleDecisionEngine()
	agent.metrics = NewSimpleMetricsCollector()

	// 启动代理
	if err := agent.Start(ctx); err != nil {
		log.Printf("Failed to start agent: %v", err)
		return
	}
	defer agent.Stop()

	// 处理一些任务
	for i := 0; i < 5; i++ {
		input := Input{
			ID:       fmt.Sprintf("task-%d", i),
			Type:     "basic_task",
			Data:     map[string]interface{}{"value": i, "timestamp": time.Now()},
			Priority: i%3 + 1,
			Timeout:  2 * time.Second,
		}

		output, err := agent.Process(ctx, input)
		if err != nil {
			log.Printf("Task processing failed: %v", err)
			continue
		}

		fmt.Printf("任务 %s 处理完成: %v\n", input.ID, output.Success)
	}

	// 显示代理状态
	status := agent.GetStatus()
	fmt.Printf("代理状态: %s, 负载: %.2f, 处理数量: %d\n",
		status.State, status.Load, status.Processed)
}

// demoMultiAgentCollaboration 演示多代理协作
func demoMultiAgentCollaboration(ctx context.Context) {
	fmt.Println("\n2. 多代理协作演示")
	fmt.Println("-----------------")

	// 创建多个代理
	agents := make([]Agent, 3)
	for i := 0; i < 3; i++ {
		config := AgentConfig{
			Name:         fmt.Sprintf("协作代理-%d", i+1),
			Type:         "collaboration",
			MaxLoad:      0.7,
			Timeout:      3 * time.Second,
			Retries:      2,
			Capabilities: []string{"collaboration", "data_processing"},
			Parameters:   map[string]interface{}{"role": fmt.Sprintf("worker-%d", i+1)},
		}

		agent := NewBaseAgent(fmt.Sprintf("collab-agent-%d", i+1), config)
		agent.learning = NewSimpleLearningEngine()
		agent.decision = NewSimpleDecisionEngine()
		agent.metrics = NewSimpleMetricsCollector()

		if err := agent.Start(ctx); err != nil {
			log.Printf("Failed to start agent %d: %v", i+1, err)
			continue
		}
		defer agent.Stop()

		agents[i] = agent
	}

	// 创建协作代理
	collabConfig := AgentConfig{
		Name:         "协作协调器",
		Type:         "collaboration_coordinator",
		MaxLoad:      0.9,
		Timeout:      10 * time.Second,
		Retries:      1,
		Capabilities: []string{"coordination", "task_decomposition"},
		Parameters:   map[string]interface{}{"max_peers": 10},
	}

	collabAgent := NewCollaborationAgent("collab-coordinator", collabConfig)
	collabAgent.learning = NewSimpleLearningEngine()
	collabAgent.decision = NewSimpleDecisionEngine()
	collabAgent.metrics = NewSimpleMetricsCollector()

	// 添加协作伙伴
	for _, agent := range agents {
		if agent != nil {
			collabAgent.AddPeer(agent)
		}
	}

	if err := collabAgent.Start(ctx); err != nil {
		log.Printf("Failed to start collaboration agent: %v", err)
		return
	}
	defer collabAgent.Stop()

	// 处理协作任务
	input := Input{
		ID:       "collaboration-task-1",
		Type:     "collaboration",
		Data:     map[string]interface{}{"complexity": "high", "participants": 3},
		Priority: 5,
		Timeout:  5 * time.Second,
	}

	output, err := collabAgent.Process(ctx, input)
	if err != nil {
		log.Printf("Collaboration task failed: %v", err)
		return
	}

	fmt.Printf("协作任务完成: %v, 参与代理数: %v\n",
		output.Success, output.Data["participants"])
}

// demoSmartCoordinator 演示智能协调器
func demoSmartCoordinator(ctx context.Context) {
	fmt.Println("\n3. 智能协调器演示")
	fmt.Println("-----------------")

	// 创建智能协调器
	coordinator := NewSmartCoordinator()
	if err := coordinator.Start(); err != nil {
		log.Printf("Failed to start coordinator: %v", err)
		return
	}
	defer coordinator.Stop()

	// 注册不同类型的代理
	agents := []Agent{
		createSpecializedAgent("data-agent", "data_processing"),
		createSpecializedAgent("decision-agent", "decision_making"),
		createSpecializedAgent("monitoring-agent", "monitoring"),
	}

	for _, agent := range agents {
		if err := coordinator.RegisterAgent(agent); err != nil {
			log.Printf("Failed to register agent: %v", err)
			continue
		}
		defer coordinator.UnregisterAgent(agent.ID())
	}

	// 提交不同类型的任务
	tasks := []Task{
		{
			ID:        "data-task-1",
			Type:      "data_processing",
			Priority:  3,
			Data:      map[string]interface{}{"size": 1000, "format": "json"},
			Timeout:   5 * time.Second,
			CreatedAt: time.Now(),
		},
		{
			ID:        "decision-task-1",
			Type:      "decision_making",
			Priority:  5,
			Data:      map[string]interface{}{"scenario": "risk_assessment"},
			Timeout:   3 * time.Second,
			CreatedAt: time.Now(),
		},
		{
			ID:        "monitoring-task-1",
			Type:      "monitoring",
			Priority:  2,
			Data:      map[string]interface{}{"metrics": []string{"cpu", "memory", "disk"}},
			Timeout:   2 * time.Second,
			CreatedAt: time.Now(),
		},
	}

	// 处理任务
	for _, task := range tasks {
		result, err := coordinator.ProcessTask(task)
		if err != nil {
			log.Printf("Task %s failed: %v", task.ID, err)
			continue
		}

		fmt.Printf("任务 %s 由代理 %s 处理完成，耗时: %v\n",
			task.ID, result.AgentID, result.Duration)
	}

	// 显示系统状态
	status := coordinator.GetSystemStatus()
	fmt.Printf("系统状态: 总代理数=%d, 活跃代理数=%d, 系统负载=%.2f\n",
		status.TotalAgents, status.ActiveAgents, status.SystemLoad)
}

// demoSpecializedAgents 演示专业代理类型
func demoSpecializedAgents(ctx context.Context) {
	fmt.Println("\n4. 专业代理类型演示")
	fmt.Println("-------------------")

	// 数据处理代理
	dataConfig := AgentConfig{
		Name:         "数据处理代理",
		Type:         "data_processing",
		MaxLoad:      0.8,
		Timeout:      5 * time.Second,
		Retries:      2,
		Capabilities: []string{"data_processing", "etl", "validation"},
		Parameters:   map[string]interface{}{"batch_size": 100},
	}

	dataAgent := NewDataProcessingAgent("data-agent-1", dataConfig)
	dataAgent.learning = NewSimpleLearningEngine()
	dataAgent.decision = NewSimpleDecisionEngine()
	dataAgent.metrics = NewSimpleMetricsCollector()

	if err := dataAgent.Start(ctx); err != nil {
		log.Printf("Failed to start data agent: %v", err)
		return
	}
	defer dataAgent.Stop()

	// 处理数据任务
	dataInput := Input{
		ID:       "data-processing-task",
		Type:     "data_processing",
		Data:     map[string]interface{}{"records": 1000, "format": "csv"},
		Priority: 4,
		Timeout:  3 * time.Second,
	}

	dataOutput, err := dataAgent.Process(ctx, dataInput)
	if err != nil {
		log.Printf("Data processing failed: %v", err)
	} else {
		fmt.Printf("数据处理完成: %v, 处理时间: %v\n",
			dataOutput.Success, dataOutput.Data["processing_time"])
	}

	// 决策代理
	decisionConfig := AgentConfig{
		Name:         "决策代理",
		Type:         "decision_making",
		MaxLoad:      0.6,
		Timeout:      3 * time.Second,
		Retries:      1,
		Capabilities: []string{"decision_making", "rule_engine", "optimization"},
		Parameters:   map[string]interface{}{"confidence_threshold": 0.8},
	}

	decisionAgent := NewDecisionAgent("decision-agent-1", decisionConfig)
	decisionAgent.learning = NewSimpleLearningEngine()
	decisionAgent.decision = NewSimpleDecisionEngine()
	decisionAgent.metrics = NewSimpleMetricsCollector()

	if err := decisionAgent.Start(ctx); err != nil {
		log.Printf("Failed to start decision agent: %v", err)
		return
	}
	defer decisionAgent.Stop()

	// 处理决策任务
	decisionInput := Input{
		ID:       "decision-task",
		Type:     "decision_making",
		Data:     map[string]interface{}{"scenario": "investment", "risk_level": "medium"},
		Priority: 6,
		Timeout:  2 * time.Second,
	}

	decisionOutput, err := decisionAgent.Process(ctx, decisionInput)
	if err != nil {
		log.Printf("Decision making failed: %v", err)
	} else {
		fmt.Printf("决策完成: %v, 匹配规则数: %v\n",
			decisionOutput.Success, decisionOutput.Data["matched_rules"])
	}

	// 监控代理
	monitoringConfig := AgentConfig{
		Name:         "监控代理",
		Type:         "monitoring",
		MaxLoad:      0.5,
		Timeout:      2 * time.Second,
		Retries:      3,
		Capabilities: []string{"monitoring", "anomaly_detection", "alerting"},
		Parameters:   map[string]interface{}{"check_interval": "5s"},
	}

	monitoringAgent := NewMonitoringAgent("monitoring-agent-1", monitoringConfig)
	monitoringAgent.learning = NewSimpleLearningEngine()
	monitoringAgent.decision = NewSimpleDecisionEngine()
	monitoringAgent.metrics = NewSimpleMetricsCollector()

	if err := monitoringAgent.Start(ctx); err != nil {
		log.Printf("Failed to start monitoring agent: %v", err)
		return
	}
	defer monitoringAgent.Stop()

	// 处理监控任务
	monitoringInput := Input{
		ID:       "monitoring-task",
		Type:     "monitoring",
		Data:     map[string]interface{}{"targets": []string{"server1", "server2"}},
		Priority: 3,
		Timeout:  1 * time.Second,
	}

	monitoringOutput, err := monitoringAgent.Process(ctx, monitoringInput)
	if err != nil {
		log.Printf("Monitoring failed: %v", err)
	} else {
		fmt.Printf("监控完成: %v, 异常数: %v, 告警数: %v\n",
			monitoringOutput.Success,
			monitoringOutput.Data["anomalies"],
			monitoringOutput.Data["alerts"])
	}
}

// createSpecializedAgent 创建专业代理的辅助函数
func createSpecializedAgent(id, agentType string) Agent {
	config := AgentConfig{
		Name:         fmt.Sprintf("%s代理", agentType),
		Type:         agentType,
		MaxLoad:      0.8,
		Timeout:      5 * time.Second,
		Retries:      2,
		Capabilities: []string{agentType},
		Parameters:   map[string]interface{}{"version": "1.0"},
	}

	var agent Agent
	switch agentType {
	case "data_processing":
		agent = NewDataProcessingAgent(id, config)
	case "decision_making":
		agent = NewDecisionAgent(id, config)
	case "monitoring":
		agent = NewMonitoringAgent(id, config)
	default:
		agent = NewBaseAgent(id, config)
	}

	// 设置通用组件
	if baseAgent, ok := agent.(*BaseAgent); ok {
		baseAgent.learning = NewSimpleLearningEngine()
		baseAgent.decision = NewSimpleDecisionEngine()
		baseAgent.metrics = NewSimpleMetricsCollector()
	}

	return agent
}
