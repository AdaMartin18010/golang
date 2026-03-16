// Package workflow 提供 Temporal 工作流和活动的应用层实现
//
// 设计原理：
// 1. 这是 Application Layer 的工作流实现
// 2. 定义工作流（Workflow）和活动（Activity）
// 3. 实现业务逻辑编排
//
// 架构位置：
// - 位置：Application Layer (internal/application/workflow/)
// - 职责：工作流定义、活动实现
// - 依赖：Domain Layer（通过应用服务）
package workflow

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"

	appuser "github.com/yourusername/golang/internal/application/user"
)

// contextKey 是 context 键类型
type contextKey string

const userServiceKey contextKey = "userService"

// WithUserService 将用户服务注入 context
func WithUserService(ctx context.Context, service *appuser.Service) context.Context {
	return context.WithValue(ctx, userServiceKey, service)
}

// getUserService 从 context 获取用户服务
func getUserService(ctx context.Context) *appuser.Service {
	service, ok := ctx.Value(userServiceKey).(*appuser.Service)
	if !ok {
		panic("user service not found in context")
	}
	return service
}

// UserWorkflowInput 用户工作流输入
type UserWorkflowInput struct {
	Email string
	Name  string
}

// UserWorkflowResult 用户工作流结果
type UserWorkflowResult struct {
	UserID    string
	Email     string
	Name      string
	CreatedAt time.Time
}

// UserWorkflow 用户创建工作流
// 这是一个 Temporal 工作流，用于编排用户创建流程
func UserWorkflow(ctx workflow.Context, input UserWorkflowInput) (*UserWorkflowResult, error) {
	// 设置活动选项
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// 步骤 1: 验证用户输入
	var validationResult bool
	err := workflow.ExecuteActivity(ctx, ValidateUserActivity, input.Email, input.Name).Get(ctx, &validationResult)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	if !validationResult {
		return nil, fmt.Errorf("invalid user input")
	}

	// 步骤 2: 创建用户
	var userID string
	err = workflow.ExecuteActivity(ctx, CreateUserActivity, input.Email, input.Name).Get(ctx, &userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 步骤 3: 发送通知
	err = workflow.ExecuteActivity(ctx, SendNotificationActivity, userID, input.Email).Get(ctx, nil)
	if err != nil {
		// 通知失败不影响主流程，记录日志即可
		workflow.GetLogger(ctx).Error("Failed to send notification", "error", err)
	}

	return &UserWorkflowResult{
		UserID:    userID,
		Email:     input.Email,
		Name:      input.Name,
		CreatedAt: workflow.Now(ctx),
	}, nil
}

// ValidateUserActivity 验证用户输入活动
func ValidateUserActivity(ctx context.Context, email, name string) (bool, error) {
	// 简单的验证逻辑
	if email == "" || name == "" {
		return false, nil
	}
	// 可以添加更复杂的验证逻辑
	return true, nil
}

// CreateUserActivity 创建用户活动
func CreateUserActivity(ctx context.Context, email, name string) (string, error) {
	service := getUserService(ctx)
	user, err := service.CreateUser(ctx, email, name)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

// UpdateUserActivity 更新用户活动
func UpdateUserActivity(ctx context.Context, userID, name string) error {
	service := getUserService(ctx)
	return service.UpdateUserName(ctx, userID, name)
}

// DeleteUserActivity 删除用户活动
func DeleteUserActivity(ctx context.Context, userID string) error {
	service := getUserService(ctx)
	return service.DeleteUser(ctx, userID)
}

// SendNotificationActivity 发送通知活动
func SendNotificationActivity(ctx context.Context, userID, email string) error {
	// 模拟发送通知
	// 实际实现可以调用邮件服务、短信服务等
	fmt.Printf("Sending notification to user %s at %s\n", userID, email)
	return nil
}
