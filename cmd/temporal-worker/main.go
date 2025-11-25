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
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建 Temporal 客户端
	temporalClient, err := temporal.NewClient(cfg.Workflow.Temporal.Address)
	if err != nil {
		log.Fatalf("Failed to create temporal client: %v", err)
	}
	defer temporalClient.Close()

	// 初始化数据库（使用 Ent 客户端）
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
	defer entClient.Close()
	log.Println("Database connected successfully")

	userRepo := entrepo.NewUserRepository(entClient)
	userService := appuser.NewService(userRepo)

	// 创建 Worker
	w := temporal.NewWorkerFromClient(temporalClient, cfg.Workflow.Temporal.TaskQueue)

	// 注册工作流
	w.RegisterWorkflow(appworkflow.UserWorkflow)

	// 注册活动（使用依赖注入）
	// 注意：Temporal 活动需要是函数，不能直接注入依赖
	// 可以通过 context 传递依赖，或使用活动结构体
	w.RegisterActivity(appworkflow.ValidateUserActivity)

	// 创建带依赖的活动包装器
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

	// 启动 Worker
	go func() {
		if err := w.Run(); err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()

	log.Printf("Temporal Worker started on task queue: %s", cfg.Workflow.Temporal.TaskQueue)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	w.Stop()
	log.Println("Worker stopped")
}
