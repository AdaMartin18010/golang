//go:build linux
// +build linux

package ebpf

// Linux 系统的 eBPF 加载实现
// 当 bpf2go 生成文件存在时，这些函数将被替换

import (
	"errors"

	"github.com/cilium/ebpf"
)

// loadSyscallObjects 加载系统调用追踪 eBPF 对象
// 注意：实际实现应由 bpf2go 生成
func loadSyscallObjects(objs *syscallObjects, opts *ebpf.CollectionOptions) error {
	// 这里应该调用 bpf2go 生成的代码
	// 例如：return loadSyscallObjectsImpl(objs, opts)
	return errors.New("eBPF objects not generated. Run: go generate ./pkg/observability/ebpf")
}

// syscallObjects Close 方法
func (o *syscallObjects) Close() error {
	// 关闭所有 maps 和 programs
	if o.Maps.SyscallEvents != nil {
		o.Maps.SyscallEvents.Close()
	}
	if o.Maps.SyscallStats != nil {
		o.Maps.SyscallStats.Close()
	}
	if o.Maps.SyscallStartTime != nil {
		o.Maps.SyscallStartTime.Close()
	}
	if o.Programs.TraceSyscallEnter != nil {
		o.Programs.TraceSyscallEnter.Close()
	}
	if o.Programs.TraceSyscallExit != nil {
		o.Programs.TraceSyscallExit.Close()
	}
	return nil
}
