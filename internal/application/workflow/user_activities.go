package workflow

import (
	"context"
	"fmt"

	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/pkg/validator"
)

// ValidateUserActivity 验证用户活动
func ValidateUserActivity(ctx context.Context, email, name string) (string, error) {
	if !validator.ValidateEmail(email) {
		return "", fmt.Errorf("invalid email: %s", email)
	}
	if !validator.ValidateName(name) {
		return "", fmt.Errorf("invalid name: %s", name)
	}
	return "validation passed", nil
}

// CreateUserActivity 创建用户活动
func CreateUserActivity(ctx context.Context, email, name string) (string, error) {
	userService, ok := GetUserServiceFromContext(ctx)
	if !ok {
		// 如果没有注入 UserService，返回临时实现
		return fmt.Sprintf("user-%s", email), nil
	}

	user, err := userService.CreateUser(ctx, appuser.CreateUserRequest{
		Email: email,
		Name:  name,
	})
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

// UpdateUserActivity 更新用户活动
func UpdateUserActivity(ctx context.Context, userID, email, name string) error {
	userService, ok := GetUserServiceFromContext(ctx)
	if !ok {
		// 如果没有注入 UserService，返回成功（临时实现）
		return nil
	}

	_, err := userService.UpdateUser(ctx, userID, appuser.UpdateUserRequest{
		Email: &email,
		Name:  &name,
	})
	return err
}

// DeleteUserActivity 删除用户活动
func DeleteUserActivity(ctx context.Context, userID string) error {
	userService, ok := GetUserServiceFromContext(ctx)
	if !ok {
		// 如果没有注入 UserService，返回成功（临时实现）
		return nil
	}

	return userService.DeleteUser(ctx, userID)
}

// SendNotificationActivity 发送通知活动
func SendNotificationActivity(ctx context.Context, userID, eventType string) error {
	// TODO: 实现通知发送逻辑
	// 可以发送邮件、短信、推送通知等
	fmt.Printf("Sending notification: userID=%s, eventType=%s\n", userID, eventType)
	return nil
}
