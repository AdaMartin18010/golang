//go:build linux
// +build linux

package ebpf

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type tcp_event network ./programs/network.bpf.c -- -I/usr/include -I/usr/include/x86_64-linux-gnu

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/rlimit"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// NetworkTracer 网络追踪器
// 使用 eBPF 追踪 TCP 连接和网络流量
type NetworkTracer struct {
	tracer trace.Tracer
	meter  metric.Meter

	// eBPF 对象
	objs       *networkObjects
	links      []link.Link
	perfReader *perf.Reader
	ctx        context.Context
	cancel     context.CancelFunc
	enabled    bool

	// 指标
	tcpConnectionCounter metric.Int64Counter
	tcpBytesCounter      metric.Int64Counter
	tcpLatencyHistogram  metric.Float64Histogram
	activeConnections    metric.Int64UpDownCounter
}

// networkObjects 将由 bpf2go 生成
type networkObjects struct {
	Programs networkPrograms
	Maps     networkMaps
}

type networkPrograms struct {
	TraceTCPConnect    *ebpf.Program
	TraceTCPAccept     *ebpf.Program
	TraceTCPSendMsg    *ebpf.Program
	TraceTCPClose      *ebpf.Program
	TraceTCPRecvMsg    *ebpf.Program
	TraceTCPRecvMsgRet *ebpf.Program
}

type networkMaps struct {
	TCPEvents      *ebpf.Map
	TCPConnections *ebpf.Map
	TCPStats       *ebpf.Map
}

// TCPEvent TCP 事件
type TCPEvent struct {
	Timestamp uint64
	PID       uint32
	TID       uint32
	EventType uint32 // 0=connect, 1=accept, 2=close
	SrcAddr   [4]byte
	DstAddr   [4]byte
	SrcPort   uint16
	DstPort   uint16
	BytesSent uint64
	BytesRecv uint64
	Duration  uint64
}

// TCPConnInfo TCP 连接信息（与 eBPF C 结构对应）
type TCPConnInfo struct {
	StartTime uint64
	BytesSent uint64
	BytesRecv uint64
	SrcAddr   [4]byte
	DstAddr   [4]byte
	SrcPort   uint16
	DstPort   uint16
	_         [4]byte // padding
}

// NetworkTracerConfig 配置
type NetworkTracerConfig struct {
	Tracer  trace.Tracer
	Meter   metric.Meter
	Enabled bool
	// TargetPID 目标进程ID（0表示所有进程）
	TargetPID uint32
	// BufferSize perf buffer 大小
	BufferSize int
	// TrackInbound 是否追踪入站连接
	TrackInbound bool
	// TrackOutbound 是否追踪出站连接
	TrackOutbound bool
}

// NewNetworkTracer 创建网络追踪器
//
// 注意：
// 1. 需要 root 权限或 CAP_BPF + CAP_NET_ADMIN
// 2. 需要 Linux 5.2+ 内核（支持 BPF_PROG_TYPE_SOCK_OPS）
// 3. 需要先编译 eBPF 程序：make generate-ebpf
func NewNetworkTracer(cfg NetworkTracerConfig) (*NetworkTracer, error) {
	if !cfg.Enabled {
		return &NetworkTracer{enabled: false}, nil
	}

	// 非 Linux 系统不支持 eBPF
	if runtime.GOOS != "linux" {
		return &NetworkTracer{enabled: false}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	tracer := &NetworkTracer{
		tracer:  cfg.Tracer,
		meter:   cfg.Meter,
		ctx:     ctx,
		cancel:  cancel,
		enabled: true,
		links:   make([]link.Link, 0, 6),
	}

	// 初始化指标
	if cfg.Meter != nil {
		var err error

		tracer.tcpConnectionCounter, err = cfg.Meter.Int64Counter(
			"ebpf.tcp.connections",
			metric.WithDescription("Total number of TCP connections"),
			metric.WithUnit("{connections}"),
		)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create connection counter: %w", err)
		}

		tracer.tcpBytesCounter, err = cfg.Meter.Int64Counter(
			"ebpf.tcp.bytes",
			metric.WithDescription("Total bytes transferred over TCP"),
			metric.WithUnit("By"),
		)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create bytes counter: %w", err)
		}

		tracer.tcpLatencyHistogram, err = cfg.Meter.Float64Histogram(
			"ebpf.tcp.connection.duration",
			metric.WithDescription("TCP connection duration"),
			metric.WithUnit("ms"),
		)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create latency histogram: %w", err)
		}

		tracer.activeConnections, err = cfg.Meter.Int64UpDownCounter(
			"ebpf.tcp.connections.active",
			metric.WithDescription("Number of active TCP connections"),
			metric.WithUnit("{connections}"),
		)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create active connections counter: %w", err)
		}
	}

	// 移除资源限制
	if err := rlimit.RemoveMemlock(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to remove memlock limit: %w", err)
	}

	// 加载 eBPF 程序
	if err := tracer.loadPrograms(cfg); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to load eBPF programs: %w", err)
	}

	return tracer, nil
}

// loadPrograms 加载 eBPF 程序
func (nt *NetworkTracer) loadPrograms(cfg NetworkTracerConfig) error {
	objs := &networkObjects{}

	// 加载 eBPF 对象
	if err := nt.loadObjects(objs); err != nil {
		return fmt.Errorf("failed to load network objects: %w", err)
	}
	nt.objs = objs

	// 附加 kprobe: tcp_connect (出站连接)
	if cfg.TrackOutbound {
		connectLink, err := link.Kprobe("tcp_connect", objs.Programs.TraceTCPConnect, nil)
		if err != nil {
			objs.Close()
			return fmt.Errorf("failed to attach tcp_connect kprobe: %w", err)
		}
		nt.links = append(nt.links, connectLink)

		// 附加 kretprobe: tcp_v4_connect
		connectRetLink, err := link.Kretprobe("tcp_v4_connect", objs.Programs.TraceTCPConnect, nil)
		if err != nil {
			// 某些内核可能没有这个函数，忽略错误
			_ = connectRetLink
		}
	}

	// 附加 kprobe: inet_csk_accept (入站连接)
	if cfg.TrackInbound {
		acceptLink, err := link.Kprobe("inet_csk_accept", objs.Programs.TraceTCPAccept, nil)
		if err != nil {
			objs.Close()
			return fmt.Errorf("failed to attach inet_csk_accept kprobe: %w", err)
		}
		nt.links = append(nt.links, acceptLink)
	}

	// 附加 kprobe: tcp_close (连接关闭)
	closeLink, err := link.Kprobe("tcp_close", objs.Programs.TraceTCPClose, nil)
	if err != nil {
		objs.Close()
		return fmt.Errorf("failed to attach tcp_close kprobe: %w", err)
	}
	nt.links = append(nt.links, closeLink)

	// 附加 kprobe: tcp_sendmsg (发送数据)
	sendLink, err := link.Kprobe("tcp_sendmsg", objs.Programs.TraceTCPSendMsg, nil)
	if err != nil {
		objs.Close()
		return fmt.Errorf("failed to attach tcp_sendmsg kprobe: %w", err)
	}
	nt.links = append(nt.links, sendLink)

	// 附加 kprobe/kretprobe: tcp_recvmsg (接收数据)
	recvLink, err := link.Kprobe("tcp_recvmsg", objs.Programs.TraceTCPRecvMsg, nil)
	if err != nil {
		objs.Close()
		return fmt.Errorf("failed to attach tcp_recvmsg kprobe: %w", err)
	}
	nt.links = append(nt.links, recvLink)

	recvRetLink, err := link.Kretprobe("tcp_recvmsg", objs.Programs.TraceTCPRecvMsgRet, nil)
	if err != nil {
		objs.Close()
		return fmt.Errorf("failed to attach tcp_recvmsg kretprobe: %w", err)
	}
	nt.links = append(nt.links, recvRetLink)

	// 创建 perf reader
	bufferSize := cfg.BufferSize
	if bufferSize == 0 {
		bufferSize = 4096 * 16 // 默认 64KB
	}

	reader, err := perf.NewReader(objs.Maps.TCPEvents, bufferSize)
	if err != nil {
		objs.Close()
		return fmt.Errorf("failed to create perf reader: %w", err)
	}
	nt.perfReader = reader

	return nil
}

// loadObjects 加载 eBPF 对象
func (nt *NetworkTracer) loadObjects(objs *networkObjects) error {
	// 这里应该调用 bpf2go 生成的 loadNetworkObjects 函数
	return errors.New("eBPF objects not generated. Run: go generate ./pkg/observability/ebpf")
}

// Start 启动追踪器
func (nt *NetworkTracer) Start() error {
	if !nt.enabled {
		return nil
	}

	// 启动 perf reader goroutine
	go nt.readEvents()

	return nil
}

// Stop 停止追踪器
func (nt *NetworkTracer) Stop() error {
	if !nt.enabled {
		return nil
	}

	// 取消上下文
	if nt.cancel != nil {
		nt.cancel()
	}

	// 清理资源
	if nt.perfReader != nil {
		nt.perfReader.Close()
	}

	for _, l := range nt.links {
		l.Close()
	}

	if nt.objs != nil {
		nt.objs.Close()
	}

	return nil
}

// readEvents 读取网络事件
func (nt *NetworkTracer) readEvents() {
	if nt.perfReader == nil {
		return
	}

	for {
		select {
		case <-nt.ctx.Done():
			return
		default:
			// 从 perf buffer 读取事件
			record, err := nt.perfReader.Read()
			if err != nil {
				if errors.Is(err, perf.ErrClosed) {
					return
				}
				continue
			}

			// 解析事件
			var event TCPEvent
			if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
				continue
			}

			// 处理事件
			nt.handleTCPEvent(nt.ctx, &event)
		}
	}
}

// handleTCPEvent 处理 TCP 事件
func (nt *NetworkTracer) handleTCPEvent(ctx context.Context, event *TCPEvent) {
	eventType := "unknown"
	switch event.EventType {
	case 0:
		eventType = "connect"
		if nt.activeConnections != nil {
			nt.activeConnections.Add(ctx, 1)
		}
	case 1:
		eventType = "accept"
		if nt.activeConnections != nil {
			nt.activeConnections.Add(ctx, 1)
		}
	case 2:
		eventType = "close"
		if nt.activeConnections != nil {
			nt.activeConnections.Add(ctx, -1)
		}
	}

	// 记录连接数
	if nt.tcpConnectionCounter != nil {
		nt.tcpConnectionCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("tcp.event_type", eventType),
				attribute.Int("process.pid", int(event.PID)),
			),
		)
	}

	// 记录字节数
	if nt.tcpBytesCounter != nil && event.EventType == 2 { // close 时记录
		totalBytes := int64(event.BytesSent + event.BytesRecv)
		if totalBytes > 0 {
			nt.tcpBytesCounter.Add(ctx, totalBytes,
				metric.WithAttributes(
					attribute.String("tcp.direction", "total"),
					attribute.Int("process.pid", int(event.PID)),
				),
			)
		}
	}

	// 记录连接延迟
	if nt.tcpLatencyHistogram != nil && event.Duration > 0 {
		duration := float64(event.Duration) / 1000000.0 // ns to ms
		nt.tcpLatencyHistogram.Record(ctx, duration,
			metric.WithAttributes(
				attribute.String("tcp.event_type", eventType),
			),
		)
	}

	// 创建追踪 span
	if nt.tracer != nil {
		srcIP := net.IPv4(event.SrcAddr[0], event.SrcAddr[1], event.SrcAddr[2], event.SrcAddr[3])
		dstIP := net.IPv4(event.DstAddr[0], event.DstAddr[1], event.DstAddr[2], event.DstAddr[3])

		_, span := nt.tracer.Start(ctx, "tcp."+eventType,
			trace.WithTimestamp(time.Unix(0, int64(event.Timestamp))),
			trace.WithAttributes(
				attribute.String("net.peer.ip", dstIP.String()),
				attribute.Int("net.peer.port", int(event.DstPort)),
				attribute.String("net.host.ip", srcIP.String()),
				attribute.Int("net.host.port", int(event.SrcPort)),
				attribute.Int("process.pid", int(event.PID)),
				attribute.Int64("tcp.bytes.sent", int64(event.BytesSent)),
				attribute.Int64("tcp.bytes.recv", int64(event.BytesRecv)),
				attribute.String("tcp.event_type", eventType),
			),
		)
		if event.Duration > 0 {
			span.End(trace.WithTimestamp(time.Unix(0, int64(event.Timestamp+event.Duration))))
		} else {
			span.End()
		}
	}
}

// IsEnabled 检查是否启用
func (nt *NetworkTracer) IsEnabled() bool {
	return nt.enabled
}

// GetActiveConnections 获取活跃连接数
func (nt *NetworkTracer) GetActiveConnections(ctx context.Context) (int64, error) {
	if !nt.enabled || nt.objs == nil {
		return 0, nil
	}

	// 从 eBPF map 读取活跃连接数
	var count int64

	// 使用 Iterate 遍历 map
	var key uint64
	var value TCPConnInfo

	iter := nt.objs.Maps.TCPConnections.Iterate()
	for iter.Next(&key, &value) {
		count++
	}

	if err := iter.Err(); err != nil {
		return 0, fmt.Errorf("failed to iterate tcp connections: %w", err)
	}

	return count, nil
}

// GetConnectionStats 获取连接统计信息
func (nt *NetworkTracer) GetConnectionStats(ctx context.Context) (map[uint32]uint64, error) {
	if !nt.enabled || nt.objs == nil {
		return nil, nil
	}

	stats := make(map[uint32]uint64)

	// 从 eBPF map 读取统计信息
	var key uint32
	var value uint64

	iter := nt.objs.Maps.TCPStats.Iterate()
	for iter.Next(&key, &value) {
		stats[key] = value
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tcp stats: %w", err)
	}

	return stats, nil
}

// GetConnectionDetails 获取连接详细信息
func (nt *NetworkTracer) GetConnectionDetails(ctx context.Context) ([]ConnectionDetail, error) {
	if !nt.enabled || nt.objs == nil {
		return nil, nil
	}

	var details []ConnectionDetail

	var key uint64
	var value TCPConnInfo

	iter := nt.objs.Maps.TCPConnections.Iterate()
	for iter.Next(&key, &value) {
		detail := ConnectionDetail{
			SocketFD:  key,
			StartTime: value.StartTime,
			BytesSent: value.BytesSent,
			BytesRecv: value.BytesRecv,
			SrcIP:     net.IPv4(value.SrcAddr[0], value.SrcAddr[1], value.SrcAddr[2], value.SrcAddr[3]),
			DstIP:     net.IPv4(value.DstAddr[0], value.DstAddr[1], value.DstAddr[2], value.DstAddr[3]),
			SrcPort:   value.SrcPort,
			DstPort:   value.DstPort,
		}
		details = append(details, detail)
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tcp connections: %w", err)
	}

	return details, nil
}

// ConnectionDetail 连接详细信息
type ConnectionDetail struct {
	SocketFD  uint64
	StartTime uint64
	BytesSent uint64
	BytesRecv uint64
	SrcIP     net.IP
	DstIP     net.IP
	SrcPort   uint16
	DstPort   uint16
}

// Close 关闭追踪器
func (nt *NetworkTracer) Close() error {
	return nt.Stop()
}
