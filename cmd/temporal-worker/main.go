// Package main Temporal Worker 应用入口
//
// 设计原理：
// 1. 这是 Temporal 工作流 Worker 的入口点
// 2. 负责执行工作流（Workflow）和活动（Activity）
// 3. 支持长时间运行的任务和复杂业务流程
// 4. 实现优雅关闭，确保正在执行的任务完成后再退出
//
// 启动流程：
// 1. 加载配置
// 2. 创建 Temporal 客户端
// 3. 初始化数据库连接
// 4. 组装依赖（仓储、服务）
// 5. 创建 Worker
// 6. 注册工作流和活动
// 7. 启动 Worker（在 goroutine 中）
// 8. 等待中断信号
// 9. 优雅关闭
//
// Temporal 说明：
// - Temporal 是工作流编排引擎
// - 支持长时间运行的任务（可以运行数天、数周）
// - 支持重试、超时、错误处理
// - 支持版本控制和迁移
//
// 使用方式：
//   go run cmd/temporal-worker/main.go
//   # 或
//   go build -o bin/temporal-worker cmd/temporal-worker/main.go
//   ./bin/temporal-worker
//
// 配置：
// - 配置文件：configs/config.yaml
// - Temporal.Address: Temporal 服务器地址（默认 localhost:7233）
// - Temporal.TaskQueue: 任务队列名称（默认 default）
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/golang/internal/config"
	appuser "github.com/yourusername/golang/internal/application/user"
	entdb "github.com/yourusername/golang/internal/infrastructure/database/ent"
	entrepo "github.com/yourusername/golang/internal/infrastructure/database/ent/repository"
	"github.com/yourusername/golang/internal/infrastructure/workflow/temporal"
	appworkflow "github.com/yourusername/golang/internal/application/workflow"
)

func main() {
	// 步骤 1: 加载配置
	//
	// 配置加载说明：
	// - 从 configs/config.yaml 加载配置
	// - 支持环境变量覆盖（APP_* 前缀）
	// - 需要配置 Temporal 服务器地址和任务队列
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 步骤 2: 创建 Temporal 客户端
	//
	// Temporal 客户端说明：
	// - 连接到 Temporal 服务器
	// - 用于启动工作流、查询工作流状态等
	// - 如果连接失败，程序无法继续运行
	temporalClient, err := temporal.NewClient(cfg.Workflow.Temporal.Address)
	if err != nil {
		log.Fatalf("Failed to create temporal client: %v", err)
	}
	// 确保在程序退出时关闭 Temporal 客户端
	defer temporalClient.Close()

	// 步骤 3: 初始化数据库（使用 Ent 客户端）
	//
	// 数据库初始化说明：
	// - 使用 Ent ORM 连接 PostgreSQL
	// - 活动（Activity）可能需要访问数据库
	// - 需要数据库连接来执行业务逻辑
	ctx := context.Background()
	entClient, err := entdb.NewClientFromConfig(
		ctx,
		cfg.Database.Host,
		fmt.Sprintf("%d", cfg.Database.Port),
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// 确保在程序退出时关闭数据库连接
	defer entClient.Close()
	log.Println("Database connected successfully")

	// 步骤 4: 组装依赖（手动依赖注入）
	//
	// 依赖注入说明：
	// - 创建仓储和应用服务
	// - 活动需要访问应用服务来执行业务逻辑
	// - 推荐迁移到 Wire 依赖注入
	userRepo := entrepo.NewUserRepository(entClient)
	userService := appuser.NewService(userRepo)

	// 步骤 5: 创建 Worker
	//
	// Worker 说明：
	// - Worker 负责执行工作流和活动
	// - 从任务队列（TaskQueue）中获取任务
	// - 可以配置并发数、重试策略等
	w := temporal.NewWorkerFromClient(temporalClient, cfg.Workflow.Temporal.TaskQueue)

	// 步骤 6: 注册工作流
	//
	// 工作流注册说明：
	// - 工作流是长时间运行的业务流程
	// - 工作流可以调用多个活动
	// - 工作流支持版本控制和迁移
	w.RegisterWorkflow(appworkflow.UserWorkflow)

	// 步骤 7: 注册活动（使用依赖注入）
	//
	// 活动注册说明：
	// - 活动是工作流中的单个步骤
	// - 活动需要是函数，不能直接注入依赖
	// - 可以通过 context 传递依赖，或使用活动包装器
	//
	// 依赖注入方式：
	// 1. 通过 context 传递依赖（推荐）
	// 2. 使用活动包装器（当前实现）
	// 3. 使用全局变量（不推荐）
	//
	// 注意：Temporal 活动需要是函数，不能直接注入依赖
	// 可以通过 context 传递依赖，或使用活动结构体
	w.RegisterActivity(appworkflow.ValidateUserActivity)

	// 创建带依赖的活动包装器
	//
	// 包装器说明：
	// - 将依赖注入到 context 中
	// - 调用实际的活动函数
	// - 这样可以在活动中访问应用服务
	createUserActivity := func(ctx context.Context, email, name string) (string, error) {
		ctx = appworkflow.WithUserService(ctx, userService)
		return appworkflow.CreateUserActivity(ctx, email, name)
	}
	updateUserActivity := func(ctx context.Context, userID, email, name string) error {
		ctx = appworkflow.WithUserService(ctx, userService)
		return appworkflow.UpdateUserActivity(ctx, userID, email, name)
	}
	deleteUserActivity := func(ctx context.Context, userID string) error {
		ctx = appworkflow.WithUserService(ctx, userService)
		return appworkflow.DeleteUserActivity(ctx, userID)
	}

	w.RegisterActivity(createUserActivity)
	w.RegisterActivity(updateUserActivity)
	w.RegisterActivity(deleteUserActivity)
	w.RegisterActivity(appworkflow.SendNotificationActivity)

	// 步骤 8: 启动 Worker
	//
	// Worker 启动说明：
	// - 在独立的 goroutine 中启动 Worker
	// - Worker 会持续从任务队列中获取任务并执行
	// - 如果 Worker 失败，记录错误并退出程序
	go func() {
		if err := w.Run(); err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()

	log.Printf("Temporal Worker started on task queue: %s", cfg.Workflow.Temporal.TaskQueue)

	// 步骤 9: 等待中断信号
	//
	// 信号监听说明：
	// - 监听 SIGINT（Ctrl+C）和 SIGTERM（kill 命令）
	// - 收到信号后，开始优雅关闭流程
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")

	// 步骤 10: 优雅关闭
	//
	// 优雅关闭说明：
	// - 停止 Worker，不再接收新任务
	// - 等待正在执行的任务完成
	// - 关闭数据库连接和 Temporal 客户端
	w.Stop()
	log.Println("Worker stopped")
}
