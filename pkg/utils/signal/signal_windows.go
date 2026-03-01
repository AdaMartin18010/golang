//go:build windows
// +build windows

package signal

import (
	"syscall"
)

// Windows signals - SIGUSR1 and SIGUSR2 are not available on Windows
// We use alternative signals for compatibility

// NotifyAll 注册所有信号通知（Windows版本）
func NotifyAllWindows() []syscall.Signal {
	return []syscall.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	}
}
