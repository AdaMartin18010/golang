//go:build !linux
// +build !linux

package ebpf

// 非 Linux 系统的 stub 实现
// 这些函数在非 Linux 平台上不会实际加载 eBPF 程序

import (
	"errors"

	"github.com/cilium/ebpf"
)

// loadNetworkObjects stub 实现
func loadNetworkObjects(objs *networkObjects, opts *ebpf.CollectionOptions) error {
	return errors.New("eBPF is only supported on Linux")
}

// networkObjects stub Close 方法
func (o *networkObjects) Close() error {
	return nil
}
