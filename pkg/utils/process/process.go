package process

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

// ProcessInfo 进程信息
type ProcessInfo struct {
	PID        int
	PPID       int
	Executable string
	Args       []string
	Env        []string
	Dir        string
}

// GetPID 获取当前进程PID
func GetPID() int {
	return os.Getpid()
}

// GetPPID 获取父进程PID
func GetPPID() int {
	return os.Getppid()
}

// GetExecutable 获取可执行文件路径
func GetExecutable() (string, error) {
	return os.Executable()
}

// GetArgs 获取命令行参数
func GetArgs() []string {
	return os.Args
}

// GetEnv 获取环境变量
func GetEnv(key string) string {
	return os.Getenv(key)
}

// SetEnv 设置环境变量
func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// GetEnvAll 获取所有环境变量
func GetEnvAll() []string {
	return os.Environ()
}

// GetWorkingDir 获取工作目录
func GetWorkingDir() (string, error) {
	return os.Getwd()
}

// ChangeDir 改变工作目录
func ChangeDir(dir string) error {
	return os.Chdir(dir)
}

// GetProcessInfo 获取进程信息
func GetProcessInfo() (*ProcessInfo, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &ProcessInfo{
		PID:        os.Getpid(),
		PPID:       os.Getppid(),
		Executable: executable,
		Args:       os.Args,
		Env:        os.Environ(),
		Dir:        dir,
	}, nil
}

// RunCommand 运行命令
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// RunCommandWithDir 在指定目录运行命令
func RunCommandWithDir(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// RunCommandWithEnv 使用指定环境变量运行命令
func RunCommandWithEnv(env []string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// RunCommandWithTimeout 带超时运行命令
func RunCommandWithTimeout(timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// StartCommand 启动命令（不等待完成）
func StartCommand(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// WaitCommand 等待命令完成
func WaitCommand(cmd *exec.Cmd) error {
	return cmd.Wait()
}

// KillProcess 杀死进程
func KillProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Kill()
}

// SignalProcess 向进程发送信号
func SignalProcess(pid int, sig os.Signal) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(sig)
}

// IsProcessRunning 检查进程是否运行
func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// WaitForProcess 等待进程结束
func WaitForProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	_, err = process.Wait()
	return err
}

// Exit 退出程序
func Exit(code int) {
	os.Exit(code)
}

// ExitSuccess 成功退出
func ExitSuccess() {
	os.Exit(0)
}

// ExitError 错误退出
func ExitError() {
	os.Exit(1)
}

// HandleSignals 处理信号
func HandleSignals(handler func(os.Signal), sigs ...os.Signal) {
	c := make(chan os.Signal, 1)
	if len(sigs) == 0 {
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	} else {
		signal.Notify(c, sigs...)
	}
	go func() {
		for sig := range c {
			handler(sig)
		}
	}()
}

// WaitForInterrupt 等待中断信号
func WaitForInterrupt() os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return <-c
}

// Daemonize 守护进程化（简单实现）
func Daemonize() error {
	// 注意：这是一个简化的实现，实际守护进程化需要更复杂的处理
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Dir = "/"
	return cmd.Start()
}

// IsRoot 检查是否以root权限运行
func IsRoot() bool {
	return os.Geteuid() == 0
}

// GetUserID 获取用户ID
func GetUserID() int {
	return os.Getuid()
}

// GetEffectiveUserID 获取有效用户ID
func GetEffectiveUserID() int {
	return os.Geteuid()
}

// GetGroupID 获取组ID
func GetGroupID() int {
	return os.Getgid()
}

// GetEffectiveGroupID 获取有效组ID
func GetEffectiveGroupID() int {
	return os.Getegid()
}

// GetHostname 获取主机名
func GetHostname() (string, error) {
	return os.Hostname()
}

// GetTempDir 获取临时目录
func GetTempDir() string {
	return os.TempDir()
}

// CreateTempFile 创建临时文件
func CreateTempFile(pattern string) (*os.File, error) {
	return os.CreateTemp("", pattern)
}

// CreateTempDir 创建临时目录
func CreateTempDir(pattern string) (string, error) {
	return os.MkdirTemp("", pattern)
}
