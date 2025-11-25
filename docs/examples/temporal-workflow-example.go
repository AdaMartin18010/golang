package main

import (
	"context"
	"fmt"
	"log"
	"time"

	appworkflow "github.com/yourusername/golang/internal/application/workflow"
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

	// 启动用户创建工作流
	input := appworkflow.UserWorkflowInput{
		UserID: "user-123",
		Email:  "test@example.com",
		Name:   "Test User",
		Action: "create",
	}

	workflowRun, err := handler.StartUserWorkflow(context.Background(), input)
	if err != nil {
		log.Fatalf("Failed to start workflow: %v", err)
	}

	fmt.Printf("Workflow started: %s, Run ID: %s\n", workflowRun.GetID(), workflowRun.GetRunID())

	// 等待工作流完成
	var result appworkflow.UserWorkflowOutput
	err = workflowRun.Get(context.Background(), &result)
	if err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	fmt.Printf("Workflow completed: %+v\n", result)
}
