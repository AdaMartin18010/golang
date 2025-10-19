package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"ai-agent-architecture/core"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║       AI-Agent 架构演示程序                      ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	// 创建简单的代理配置
	config := core.AgentConfig{
		Name:         "demo-agent",
		Type:         "general",
		MaxLoad:      100,
		Timeout:      30 * time.Second,
		Retries:      3,
		Capabilities: []string{"text-processing", "data-analysis"},
		Parameters: map[string]interface{}{
			"model":   "demo-v1",
			"version": "1.0",
			"enabled": true,
		},
	}

	// 创建基础代理
	agent := core.NewBaseAgent("agent-001", config)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 启动代理
	fmt.Println("启动代理...")
	if err := agent.Start(ctx); err != nil {
		log.Fatalf("启动代理失败: %v", err)
	}
	fmt.Println("✓ 代理已启动")

	// 创建测试输入
	input := core.Input{
		ID:   "input-001",
		Type: "text",
		Data: map[string]interface{}{
			"message": "Hello, AI Agent!",
			"action":  "process",
		},
		Metadata: map[string]interface{}{
			"source":   "demo",
			"priority": 1,
		},
		Timestamp: time.Now(),
	}

	// 处理输入
	fmt.Println("\n处理输入...")
	output, err := agent.Process(input)
	if err != nil {
		log.Printf("处理失败: %v", err)
	} else {
		fmt.Printf("✓ 处理成功\n")
		fmt.Printf("  - 输出ID: %s\n", output.ID)
		fmt.Printf("  - 输出类型: %s\n", output.Type)
		if output.Data != nil {
			fmt.Printf("  - 输出数据: %v\n", output.Data)
		}
	}

	// 获取代理状态
	fmt.Println("\n获取代理状态...")
	status := agent.GetStatus()
	fmt.Printf("✓ 当前状态: %s\n", status.State)
	fmt.Printf("  - 负载: %.2f\n", status.Load)
	if len(status.Metrics) > 0 {
		fmt.Println("  - 指标:")
		for k, v := range status.Metrics {
			fmt.Printf("    • %s: %.2f\n", k, v)
		}
	}

	// 创建学习经验
	experience := core.Experience{
		Input:     input,
		Output:    output,
		Reward:    0.85,
		Timestamp: time.Now(),
	}

	// 学习经验
	fmt.Println("\n学习经验...")
	if err := agent.Learn(experience); err != nil {
		log.Printf("学习失败: %v", err)
	} else {
		fmt.Println("✓ 学习完成")
	}

	// 停止代理
	fmt.Println("\n停止代理...")
	if err := agent.Stop(); err != nil {
		log.Printf("停止失败: %v", err)
	} else {
		fmt.Println("✓ 代理已停止")
	}

	fmt.Println("\n演示完成！")
}
