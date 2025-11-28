package main

import "context"

type contextKey string

const (
	userIDKey contextKey = "userID"
	requestIDKey contextKey = "requestID"
)

// WithValue 在context中传递值
func WithValue(userID, requestID string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, requestIDKey, requestID)
	return ctx
}

// GetUserID 从context获取用户ID
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetRequestID 从context获取请求ID
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}
