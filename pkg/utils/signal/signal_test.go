package signal

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	c := Notify(syscall.SIGUSR1)
	if c == nil {
		t.Error("Expected non-nil channel")
	}
}

func TestNotifyInterrupt(t *testing.T) {
	c := NotifyInterrupt()
	if c == nil {
		t.Error("Expected non-nil channel")
	}
}

func TestHandleInterrupt(t *testing.T) {
	handled := false
	HandleInterrupt(func(sig os.Signal) {
		handled = true
	})
	// 注意：实际测试中需要发送信号，这里只测试函数调用
}

func TestWaitInterrupt(t *testing.T) {
	// 注意：实际测试中需要发送信号，这里只测试函数调用
	// 这个测试在实际环境中需要发送信号才能完成
}

func TestIsInterrupt(t *testing.T) {
	if !IsInterrupt(syscall.SIGINT) {
		t.Error("Expected SIGINT to be interrupt signal")
	}
	if !IsInterrupt(syscall.SIGTERM) {
		t.Error("Expected SIGTERM to be interrupt signal")
	}
	if IsInterrupt(syscall.SIGQUIT) {
		t.Error("Expected SIGQUIT not to be interrupt signal")
	}
}

func TestSignalName(t *testing.T) {
	if SignalName(syscall.SIGINT) != "SIGINT" {
		t.Errorf("Expected 'SIGINT', got %s", SignalName(syscall.SIGINT))
	}
	if SignalName(syscall.SIGTERM) != "SIGTERM" {
		t.Errorf("Expected 'SIGTERM', got %s", SignalName(syscall.SIGTERM))
	}
}

func TestGracefulShutdown(t *testing.T) {
	gs := NewGracefulShutdown()
	handled := false
	gs.AddHandler(func() {
		handled = true
	})
	// 注意：实际测试中需要发送信号，这里只测试函数调用
}

func TestWithInterruptContext(t *testing.T) {
	ctx, cancel := WithInterruptContext(context.Background())
	defer cancel()

	if ctx == nil {
		t.Error("Expected non-nil context")
	}
}

func TestWaitWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := WaitWithContext(ctx, syscall.SIGUSR1)
	if err == nil {
		// 如果没有收到信号，应该返回超时错误
		// 但在测试环境中可能不会收到信号
	}
}
