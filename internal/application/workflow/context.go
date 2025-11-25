package workflow

import (
	"context"

	appuser "github.com/yourusername/golang/internal/application/user"
)

type contextKey string

const userServiceKey contextKey = "user_service"

// WithUserService 将 UserService 注入到 context
func WithUserService(ctx context.Context, service appuser.Service) context.Context {
	return context.WithValue(ctx, userServiceKey, service)
}

// GetUserServiceFromContext 从 context 获取 UserService
func GetUserServiceFromContext(ctx context.Context) (appuser.Service, bool) {
	service, ok := ctx.Value(userServiceKey).(appuser.Service)
	return service, ok
}
