//go:build !linux
// +build !linux

package ebpf

// 非 Linux 系统的 stub 实现
// 这些函数在非 Linux 平台上不会实际加载 eBPF 程序

import (
	"errors"
)

// loadSyscallObjects stub 实现
func loadSyscallObjects(objs *syscallObjects, opts *ebpf.CollectionOptions) error {
	return errors.New("eBPF is only supported on Linux")
}

// syscallObjects stub Close 方法
func (o *syscallObjects) Close() error {
	return nil
}
