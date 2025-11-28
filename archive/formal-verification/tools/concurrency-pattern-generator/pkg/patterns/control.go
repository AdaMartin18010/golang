// Package patterns - 控制流模式
package patterns

import "fmt"

// GenerateContextCancellation 生成Context取消模式
func GenerateContextCancellation(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport (\n\t\"context\"\n\t\"time\"\n)\n\n// WithCancel 创建可取消的context\nfunc WithCancel() {\n\tctx, cancel := context.WithCancel(context.Background())\n\tdefer cancel()\n\t\n\tgo func() {\n\t\tfor {\n\t\t\tselect {\n\t\t\tcase <-ctx.Done():\n\t\t\t\treturn\n\t\t\tdefault:\n\t\t\t\t// Do work\n\t\t\t}\n\t\t}\n\t}()\n\t\n\ttime.Sleep(time.Second)\n\tcancel() // 取消所有goroutine\n}\n", pkg)
}

// GenerateContextTimeout 生成Context超时模式
func GenerateContextTimeout(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport (\n\t\"context\"\n\t\"time\"\n)\n\n// WithTimeout 创建带超时的context\nfunc WithTimeout() error {\n\tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\n\tdefer cancel()\n\t\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tcase <-time.After(3 * time.Second):\n\t\treturn nil\n\t}\n}\n\n// WithDeadline 创建带截止时间的context\nfunc WithDeadline() error {\n\tdeadline := time.Now().Add(5 * time.Second)\n\tctx, cancel := context.WithDeadline(context.Background(), deadline)\n\tdefer cancel()\n\t\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\treturn nil\n\t}\n}\n", pkg)
}

// GenerateContextValue 生成Context传值模式
func GenerateContextValue(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"context\"\n\ntype contextKey string\n\nconst (\n\tuserIDKey contextKey = \"userID\"\n\trequestIDKey contextKey = \"requestID\"\n)\n\n// WithValue 在context中传递值\nfunc WithValue(userID, requestID string) context.Context {\n\tctx := context.Background()\n\tctx = context.WithValue(ctx, userIDKey, userID)\n\tctx = context.WithValue(ctx, requestIDKey, requestID)\n\treturn ctx\n}\n\n// GetUserID 从context获取用户ID\nfunc GetUserID(ctx context.Context) string {\n\tif userID, ok := ctx.Value(userIDKey).(string); ok {\n\t\treturn userID\n\t}\n\treturn \"\"\n}\n\n// GetRequestID 从context获取请求ID\nfunc GetRequestID(ctx context.Context) string {\n\tif reqID, ok := ctx.Value(requestIDKey).(string); ok {\n\t\treturn reqID\n\t}\n\treturn \"\"\n}\n", pkg)
}

// GenerateGracefulShutdown 生成优雅关闭模式
func GenerateGracefulShutdown(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport (\n\t\"context\"\n\t\"os\"\n\t\"os/signal\"\n\t\"syscall\"\n\t\"time\"\n)\n\n// GracefulShutdown 优雅关闭服务\nfunc GracefulShutdown() {\n\tctx, cancel := context.WithCancel(context.Background())\n\tdefer cancel()\n\t\n\t// 监听系统信号\n\tsigChan := make(chan os.Signal, 1)\n\tsignal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)\n\t\n\tgo func() {\n\t\t<-sigChan\n\t\tcancel()\n\t}()\n\t\n\t// 运行服务\n\t<-ctx.Done()\n\t\n\t// 给予清理时间\n\tshutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)\n\tdefer shutdownCancel()\n\t\n\t<-shutdownCtx.Done()\n}\n", pkg)
}

// GenerateRateLimiting 生成限流模式
func GenerateRateLimiting(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport (\n\t\"context\"\n\t\"time\"\n)\n\n// RateLimiter 简单的限流器\ntype RateLimiter struct {\n\tticker *time.Ticker\n}\n\n// NewRateLimiter 创建限流器\nfunc NewRateLimiter(rate time.Duration) *RateLimiter {\n\treturn &RateLimiter{\n\t\tticker: time.NewTicker(rate),\n\t}\n}\n\n// Wait 等待令牌\nfunc (r *RateLimiter) Wait(ctx context.Context) error {\n\tselect {\n\tcase <-r.ticker.C:\n\t\treturn nil\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\t}\n}\n\n// Close 关闭限流器\nfunc (r *RateLimiter) Close() {\n\tr.ticker.Stop()\n}\n\n// TokenBucket 令牌桶限流器\ntype TokenBucket struct {\n\ttokens chan struct{}\n}\n\n// NewTokenBucket 创建令牌桶\nfunc NewTokenBucket(capacity int, rate time.Duration) *TokenBucket {\n\ttb := &TokenBucket{\n\t\ttokens: make(chan struct{}, capacity),\n\t}\n\t\n\t// 填充初始令牌\n\tfor i := 0; i < capacity; i++ {\n\t\ttb.tokens <- struct{}{}\n\t}\n\t\n\t// 定期补充令牌\n\tgo func() {\n\t\tticker := time.NewTicker(rate)\n\t\tdefer ticker.Stop()\n\t\tfor range ticker.C {\n\t\t\tselect {\n\t\t\tcase tb.tokens <- struct{}{}:\n\t\t\tdefault:\n\t\t\t}\n\t\t}\n\t}()\n\t\n\treturn tb\n}\n\n// Acquire 获取令牌\nfunc (tb *TokenBucket) Acquire(ctx context.Context) error {\n\tselect {\n\tcase <-tb.tokens:\n\t\treturn nil\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\t}\n}\n", pkg)
}
