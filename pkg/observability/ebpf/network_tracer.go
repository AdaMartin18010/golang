package ebpf

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type tcp_event network ./programs/network.bpf.c -- -I/usr/include -I/usr/include/x86_64-linux-gnu

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
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
	TraceTCPConnect *ebpf.Program
	TraceTCPAccept  *ebpf.Program
	TraceTCPSendMsg *ebpf.Program
	TraceTCPClose   *ebpf.Program
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

	ctx, cancel := context.WithCancel(context.Background())

	tracer := &NetworkTracer{
		tracer:  cfg.Tracer,
		meter:   cfg.Meter,
		ctx:     ctx,
		cancel:  cancel,
		enabled: true,
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

	// TODO: 加载 eBPF 程序
	// 实际实现需要：
	// 1. 加载 network.bpf.c 编译的对象
	// 2. 附加到 kprobe/tracepoint
	// 3. 创建 perf reader

	return tracer, nil
}

// Start 启动追踪器
func (nt *NetworkTracer) Start() error {
	if !nt.enabled {
		return nil
	}

	// TODO: 实际实现需要启动 perf reader
	// go nt.readEvents()

	return nil
}

// Stop 停止追踪器
func (nt *NetworkTracer) Stop() error {
	if !nt.enabled {
		return nil
	}

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
		// nt.objs.Close()
	}

	return nil
}

// readEvents 读取网络事件
func (nt *NetworkTracer) readEvents() {
	for {
		select {
		case <-nt.ctx.Done():
			return
		default:
			// TODO: 从 perf buffer 读取事件
			// record, err := nt.perfReader.Read()
			// ...
			// nt.handleTCPEvent(context.Background(), &event)
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
		nt.tcpBytesCounter.Add(ctx, totalBytes,
			metric.WithAttributes(
				attribute.String("tcp.direction", "sent"),
				attribute.Int("process.pid", int(event.PID)),
			),
		)
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

	// TODO: 从 eBPF map 读取活跃连接数
	// var count int64
	// iter := nt.objs.TCPConnections.Iterate()
	// ...
	// return count, nil

	return 0, nil
}
