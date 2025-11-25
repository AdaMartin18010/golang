package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// UserWorkflowInput 用户工作流输入
type UserWorkflowInput struct {
	UserID  string
	Email   string
	Name    string
	Action  string // "create", "update", "delete"
}

// UserWorkflowOutput 用户工作流输出
type UserWorkflowOutput struct {
	UserID    string
	Success   bool
	Message   string
	Timestamp time.Time
}

// UserWorkflow 用户工作流
func UserWorkflow(ctx workflow.Context, input UserWorkflowInput) (UserWorkflowOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result UserWorkflowOutput
	var err error

	switch input.Action {
	case "create":
		result, err = createUserWorkflow(ctx, input)
	case "update":
		result, err = updateUserWorkflow(ctx, input)
	case "delete":
		result, err = deleteUserWorkflow(ctx, input)
	default:
		return UserWorkflowOutput{
			UserID:    input.UserID,
			Success:   false,
			Message:   "unknown action",
			Timestamp: workflow.Now(ctx),
		}, fmt.Errorf("unknown action: %s", input.Action)
	}

	return result, err
}

// createUserWorkflow 创建用户工作流
func createUserWorkflow(ctx workflow.Context, input UserWorkflowInput) (UserWorkflowOutput, error) {
	// 1. 验证用户信息
	var validationResult string
	err := workflow.ExecuteActivity(ctx, ValidateUserActivity, input.Email, input.Name).Get(ctx, &validationResult)
	if err != nil {
		return UserWorkflowOutput{
			UserID:    input.UserID,
			Success:   false,
			Message:   fmt.Sprintf("validation failed: %v", err),
			Timestamp: workflow.Now(ctx),
		}, err
	}

	// 2. 创建用户
	var userID string
	err = workflow.ExecuteActivity(ctx, CreateUserActivity, input.Email, input.Name).Get(ctx, &userID)
	if err != nil {
		return UserWorkflowOutput{
			UserID:    input.UserID,
			Success:   false,
			Message:   fmt.Sprintf("create failed: %v", err),
			Timestamp: workflow.Now(ctx),
		}, err
	}

	// 3. 发送通知
	_ = workflow.ExecuteActivity(ctx, SendNotificationActivity, userID, "user_created").Get(ctx, nil)

	return UserWorkflowOutput{
		UserID:    userID,
		Success:   true,
		Message:   "user created successfully",
		Timestamp: workflow.Now(ctx),
	}, nil
}

// updateUserWorkflow 更新用户工作流
func updateUserWorkflow(ctx workflow.Context, input UserWorkflowInput) (UserWorkflowOutput, error) {
	// 1. 更新用户
	err := workflow.ExecuteActivity(ctx, UpdateUserActivity, input.UserID, input.Email, input.Name).Get(ctx, nil)
	if err != nil {
		return UserWorkflowOutput{
			UserID:    input.UserID,
			Success:   false,
			Message:   fmt.Sprintf("update failed: %v", err),
			Timestamp: workflow.Now(ctx),
		}, err
	}

	// 2. 发送通知
	_ = workflow.ExecuteActivity(ctx, SendNotificationActivity, input.UserID, "user_updated").Get(ctx, nil)

	return UserWorkflowOutput{
		UserID:    input.UserID,
		Success:   true,
		Message:   "user updated successfully",
		Timestamp: workflow.Now(ctx),
	}, nil
}

// deleteUserWorkflow 删除用户工作流
func deleteUserWorkflow(ctx workflow.Context, input UserWorkflowInput) (UserWorkflowOutput, error) {
	// 1. 删除用户
	err := workflow.ExecuteActivity(ctx, DeleteUserActivity, input.UserID).Get(ctx, nil)
	if err != nil {
		return UserWorkflowOutput{
			UserID:    input.UserID,
			Success:   false,
			Message:   fmt.Sprintf("delete failed: %v", err),
			Timestamp: workflow.Now(ctx),
		}, err
	}

	// 2. 发送通知
	_ = workflow.ExecuteActivity(ctx, SendNotificationActivity, input.UserID, "user_deleted").Get(ctx, nil)

	return UserWorkflowOutput{
		UserID:    input.UserID,
		Success:   true,
		Message:   "user deleted successfully",
		Timestamp: workflow.Now(ctx),
	}, nil
}
