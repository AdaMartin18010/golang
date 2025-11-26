package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Handler 信号处理函数类型
type Handler func(os.Signal)

// Notify 注册信号通知
func Notify(sig ...os.Signal) <-chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	return c
}

// NotifyAll 注册所有信号通知
func NotifyAll() <-chan os.Signal {
	return Notify(
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
}

// NotifyInterrupt 注册中断信号通知（SIGINT, SIGTERM）
func NotifyInterrupt() <-chan os.Signal {
	return Notify(syscall.SIGINT, syscall.SIGTERM)
}

// NotifyQuit 注册退出信号通知（SIGQUIT）
func NotifyQuit() <-chan os.Signal {
	return Notify(syscall.SIGQUIT)
}

// NotifyHangup 注册挂起信号通知（SIGHUP）
func NotifyHangup() <-chan os.Signal {
	return Notify(syscall.SIGHUP)
}

// NotifyUser1 注册用户信号1通知（SIGUSR1）
func NotifyUser1() <-chan os.Signal {
	return Notify(syscall.SIGUSR1)
}

// NotifyUser2 注册用户信号2通知（SIGUSR2）
func NotifyUser2() <-chan os.Signal {
	return Notify(syscall.SIGUSR2)
}

// Handle 处理信号
func Handle(handler Handler, sig ...os.Signal) {
	c := Notify(sig...)
	go func() {
		for s := range c {
			handler(s)
		}
	}()
}

// HandleInterrupt 处理中断信号
func HandleInterrupt(handler Handler) {
	Handle(handler, syscall.SIGINT, syscall.SIGTERM)
}

// HandleQuit 处理退出信号
func HandleQuit(handler Handler) {
	Handle(handler, syscall.SIGQUIT)
}

// HandleHangup 处理挂起信号
func HandleHangup(handler Handler) {
	Handle(handler, syscall.SIGHUP)
}

// Wait 等待信号
func Wait(sig ...os.Signal) os.Signal {
	c := Notify(sig...)
	return <-c
}

// WaitInterrupt 等待中断信号
func WaitInterrupt() os.Signal {
	return Wait(syscall.SIGINT, syscall.SIGTERM)
}

// WaitQuit 等待退出信号
func WaitQuit() os.Signal {
	return Wait(syscall.SIGQUIT)
}

// WaitHangup 等待挂起信号
func WaitHangup() os.Signal {
	return Wait(syscall.SIGHUP)
}

// WaitWithContext 使用context等待信号
func WaitWithContext(ctx context.Context, sig ...os.Signal) (os.Signal, error) {
	c := Notify(sig...)
	select {
	case s := <-c:
		return s, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// WaitInterruptWithContext 使用context等待中断信号
func WaitInterruptWithContext(ctx context.Context) (os.Signal, error) {
	return WaitWithContext(ctx, syscall.SIGINT, syscall.SIGTERM)
}

// GracefulShutdown 优雅关闭处理
type GracefulShutdown struct {
	sigChan  <-chan os.Signal
	handlers []func()
}

// NewGracefulShutdown 创建优雅关闭处理器
func NewGracefulShutdown() *GracefulShutdown {
	return &GracefulShutdown{
		sigChan:  NotifyInterrupt(),
		handlers: make([]func(), 0),
	}
}

// AddHandler 添加关闭处理函数
func (gs *GracefulShutdown) AddHandler(handler func()) {
	gs.handlers = append(gs.handlers, handler)
}

// Wait 等待关闭信号并执行处理函数
func (gs *GracefulShutdown) Wait() os.Signal {
	sig := <-gs.sigChan
	for _, handler := range gs.handlers {
		handler()
	}
	return sig
}

// WaitWithContext 使用context等待关闭信号
func (gs *GracefulShutdown) WaitWithContext(ctx context.Context) (os.Signal, error) {
	select {
	case sig := <-gs.sigChan:
		for _, handler := range gs.handlers {
			handler()
		}
		return sig, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Start 启动优雅关闭处理（异步）
func (gs *GracefulShutdown) Start() <-chan os.Signal {
	done := make(chan os.Signal, 1)
	go func() {
		sig := gs.Wait()
		done <- sig
	}()
	return done
}

// Ignore 忽略信号
func Ignore(sig ...os.Signal) {
	signal.Ignore(sig...)
}

// Reset 重置信号处理
func Reset(sig ...os.Signal) {
	signal.Reset(sig...)
}

// Stop 停止信号通知
func Stop(c chan<- os.Signal) {
	signal.Stop(c)
}

// Send 发送信号到进程
func Send(pid int, sig os.Signal) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(sig)
}

// SendInterrupt 发送中断信号到进程
func SendInterrupt(pid int) error {
	return Send(pid, syscall.SIGINT)
}

// SendTerminate 发送终止信号到进程
func SendTerminate(pid int) error {
	return Send(pid, syscall.SIGTERM)
}

// SendKill 发送杀死信号到进程
func SendKill(pid int) error {
	return Send(pid, syscall.SIGKILL)
}

// IsInterrupt 检查信号是否为中断信号
func IsInterrupt(sig os.Signal) bool {
	return sig == syscall.SIGINT || sig == syscall.SIGTERM
}

// IsQuit 检查信号是否为退出信号
func IsQuit(sig os.Signal) bool {
	return sig == syscall.SIGQUIT
}

// IsHangup 检查信号是否为挂起信号
func IsHangup(sig os.Signal) bool {
	return sig == syscall.SIGHUP
}

// IsUser1 检查信号是否为用户信号1
func IsUser1(sig os.Signal) bool {
	return sig == syscall.SIGUSR1
}

// IsUser2 检查信号是否为用户信号2
func IsUser2(sig os.Signal) bool {
	return sig == syscall.SIGUSR2
}

// SignalName 获取信号名称
func SignalName(sig os.Signal) string {
	switch sig {
	case syscall.SIGINT:
		return "SIGINT"
	case syscall.SIGTERM:
		return "SIGTERM"
	case syscall.SIGQUIT:
		return "SIGQUIT"
	case syscall.SIGHUP:
		return "SIGHUP"
	case syscall.SIGUSR1:
		return "SIGUSR1"
	case syscall.SIGUSR2:
		return "SIGUSR2"
	case syscall.SIGKILL:
		return "SIGKILL"
	default:
		return sig.String()
	}
}

// WithContext 创建带信号取消的context
func WithContext(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	c := Notify(sig...)
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

// WithInterruptContext 创建带中断信号取消的context
func WithInterruptContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithContext(ctx, syscall.SIGINT, syscall.SIGTERM)
}
