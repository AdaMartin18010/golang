package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/golang/internal/infrastructure/workflow/temporal"
	temporalhandler "github.com/yourusername/golang/internal/interfaces/workflow/temporal"
)

func main() {
	// 创建 Temporal 客户端
	temporalClient, err := temporal.NewClient("localhost:7233")
	if err != nil {
		log.Fatalf("Failed to create temporal client: %v", err)
	}
	defer temporalClient.Close()

	// 创建 Handler
	handler := temporalhandler.NewHandler(temporalClient.Client())

	// 发送信号给工作流
	workflowID := "user-workflow-create-user-123"
	runID := "" // 如果为空，会发送给最新的运行

	err = handler.SignalWorkflow(context.Background(), workflowID, runID, "cancel", nil)
	if err != nil {
		log.Fatalf("Failed to signal workflow: %v", err)
	}

	fmt.Printf("Signal sent to workflow: %s\n", workflowID)
}
