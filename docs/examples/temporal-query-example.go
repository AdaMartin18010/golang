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

	// 查询工作流状态
	workflowID := "user-workflow-create-user-123"
	runID := "" // 如果为空，会查询最新的运行

	result, err := handler.QueryWorkflow(context.Background(), workflowID, runID, "status", nil)
	if err != nil {
		log.Fatalf("Failed to query workflow: %v", err)
	}

	fmt.Printf("Workflow status: %v\n", result)
}
