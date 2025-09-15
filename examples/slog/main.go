package main

import (
	"errors"
	"log/slog"
	"os"
	"time"
)

func main() {
	// JSON Handler 带级别与时间戳
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("service starting", slog.String("version", "v0.1.0"))

	// 带上下文字段
	reqLog := logger.With(slog.String("component", "http"), slog.String("route", "/api/users/{id}"))
	reqLog.Info("request received", slog.Int("user_id", 123))

	// 错误记录
	err := errors.New("database timeout")
	logger.Error("operation failed", slog.Any("err", err), slog.Duration("retry_in", 2*time.Second))

	logger.Info("service stopped")
}
