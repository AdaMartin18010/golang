package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GracefulShutdown 优雅关闭服务
func GracefulShutdown() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// 监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		cancel()
	}()
	
	// 运行服务
	<-ctx.Done()
	
	// 给予清理时间
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	
	<-shutdownCtx.Done()
}
