//go:build !windows
// +build !windows

package signal

import (
	"syscall"
)

// Unix signals
var (
	SIGUSR1 = syscall.SIGUSR1
	SIGUSR2 = syscall.SIGUSR2
)

// NotifyAll 注册所有信号通知（Unix版本）
func NotifyAllUnix() []syscall.Signal {
	return []syscall.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	}
}
