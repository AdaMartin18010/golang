//go:build linux
// +build linux

package ebpf

// Linux 系统的 eBPF 加载实现
// 当 bpf2go 生成文件存在时，这些函数将被替换

import (
	"errors"

	"github.com/cilium/ebpf"
)

// loadNetworkObjects 加载网络追踪 eBPF 对象
// 注意：实际实现应由 bpf2go 生成
func loadNetworkObjects(objs *networkObjects, opts *ebpf.CollectionOptions) error {
	// 这里应该调用 bpf2go 生成的代码
	// 例如：return loadNetworkObjectsImpl(objs, opts)
	return errors.New("eBPF objects not generated. Run: go generate ./pkg/observability/ebpf")
}

// networkObjects Close 方法
func (o *networkObjects) Close() error {
	// 关闭所有 maps 和 programs
	if o.Maps.TCPEvents != nil {
		o.Maps.TCPEvents.Close()
	}
	if o.Maps.TCPConnections != nil {
		o.Maps.TCPConnections.Close()
	}
	if o.Maps.TCPStats != nil {
		o.Maps.TCPStats.Close()
	}
	if o.Programs.TraceTCPConnect != nil {
		o.Programs.TraceTCPConnect.Close()
	}
	if o.Programs.TraceTCPAccept != nil {
		o.Programs.TraceTCPAccept.Close()
	}
	if o.Programs.TraceTCPSendMsg != nil {
		o.Programs.TraceTCPSendMsg.Close()
	}
	if o.Programs.TraceTCPClose != nil {
		o.Programs.TraceTCPClose.Close()
	}
	if o.Programs.TraceTCPRecvMsg != nil {
		o.Programs.TraceTCPRecvMsg.Close()
	}
	if o.Programs.TraceTCPRecvMsgRet != nil {
		o.Programs.TraceTCPRecvMsgRet.Close()
	}
	return nil
}
