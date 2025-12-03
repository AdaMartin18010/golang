package ebpf

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type syscall_event syscall ./programs/syscall.bpf.c -- -I/usr/include -I/usr/include/x86_64-linux-gnu

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// SyscallTracer 系统调用追踪器
// 使用 Cilium eBPF 库实现真正的系统调用追踪
type SyscallTracer struct {
	tracer trace.Tracer
	meter  metric.Meter

	// eBPF 对象
	objs         *syscallObjects // bpf2go 生成的对象
	link         link.Link       // tracepoint link
	perfReader   *perf.Reader    // perf event reader
	ctx          context.Context
	cancel       context.CancelFunc
	enabled      bool

	// 指标
	syscallCounter  metric.Int64Counter
	syscallLatency  metric.Float64Histogram
}

// syscallObjects 将由 bpf2go 生成
// 这里提供接口定义
type syscallObjects struct {
	Programs syscallPrograms
	Maps     syscallMaps
}

type syscallPrograms struct {
	TraceSyscallEnter *ebpf.Program
	TraceSyscallExit  *ebpf.Program
}

type syscallMaps struct {
	SyscallEvents *ebpf.Map
	SyscallStats  *ebpf.Map
}

// SyscallEvent 系统调用事件
// 与 eBPF C 程序中的结构体对应
type SyscallEvent struct {
	Timestamp uint64
	PID       uint32
	TID       uint32
	Syscall   uint64
	Duration  uint64
	RetVal    int64
}

// SyscallTracerConfig 配置
type SyscallTracerConfig struct {
	Tracer  trace.Tracer
	Meter   metric.Meter
	Enabled bool
	// TargetPID 目标进程ID（0表示所有进程）
	TargetPID uint32
	// BufferSize perf buffer 大小
	BufferSize int
}

// NewSyscallTracer 创建系统调用追踪器
//
// 注意：
// 1. 需要 root 权限或 CAP_BPF capability
// 2. 需要 Linux 4.18+ 内核
// 3. 需要先编译 eBPF 程序：make generate-ebpf
func NewSyscallTracer(cfg SyscallTracerConfig) (*SyscallTracer, error) {
	if !cfg.Enabled {
		return &SyscallTracer{enabled: false}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	tracer := &SyscallTracer{
		tracer:  cfg.Tracer,
		meter:   cfg.Meter,
		ctx:     ctx,
		cancel:  cancel,
		enabled: true,
	}

	// 初始化指标
	if cfg.Meter != nil {
		var err error
		tracer.syscallCounter, err = cfg.Meter.Int64Counter(
			"ebpf.syscall.count",
			metric.WithDescription("Total number of system calls"),
			metric.WithUnit("{calls}"),
		)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create syscall counter: %w", err)
		}

		tracer.syscallLatency, err = cfg.Meter.Float64Histogram(
			"ebpf.syscall.duration",
			metric.WithDescription("System call duration"),
			metric.WithUnit("ms"),
		)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create syscall latency histogram: %w", err)
		}
	}

	// TODO: 加载 eBPF 程序
	// 当前为框架实现，实际使用需要：
	// 1. 编译 eBPF C 程序：
	//    go generate ./pkg/observability/ebpf
	// 2. 加载程序：
	//    objs, err := loadSyscallObjects()
	// 3. 附加到 tracepoint：
	//    link, err := link.Tracepoint("raw_syscalls", "sys_enter", objs.TraceSyscallEnter)
	// 4. 创建 perf reader：
	//    reader, err := perf.NewReader(objs.SyscallEvents, bufferSize)

	return tracer, nil
}

// Start 启动追踪器
func (st *SyscallTracer) Start() error {
	if !st.enabled {
		return nil
	}

	// TODO: 实际实现需要启动 perf reader
	// go st.readEvents()

	return nil
}

// Stop 停止追踪器
func (st *SyscallTracer) Stop() error {
	if !st.enabled {
		return nil
	}

	if st.cancel != nil {
		st.cancel()
	}

	// 清理资源
	if st.perfReader != nil {
		st.perfReader.Close()
	}
	if st.link != nil {
		st.link.Close()
	}
	if st.objs != nil {
		// st.objs.Close() // bpf2go 生成的对象会有 Close 方法
	}

	return nil
}

// readEvents 读取 eBPF 事件
// 这是实际的事件处理循环
func (st *SyscallTracer) readEvents() {
	for {
		select {
		case <-st.ctx.Done():
			return
		default:
			// TODO: 从 perf buffer 读取事件
			// record, err := st.perfReader.Read()
			// if err != nil {
			//     if errors.Is(err, perf.ErrClosed) {
			//         return
			//     }
			//     continue
			// }
			//
			// // 解析事件
			// var event SyscallEvent
			// if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			//     continue
			// }
			//
			// // 处理事件
			// st.handleSyscallEvent(context.Background(), &event)
		}
	}
}

// handleSyscallEvent 处理系统调用事件
func (st *SyscallTracer) handleSyscallEvent(ctx context.Context, event *SyscallEvent) {
	// 记录指标
	if st.syscallCounter != nil {
		st.syscallCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.Int("syscall.id", int(event.Syscall)),
				attribute.Int("process.pid", int(event.PID)),
			),
		)
	}

	if st.syscallLatency != nil {
		duration := float64(event.Duration) / 1000000.0 // ns to ms
		st.syscallLatency.Record(ctx, duration,
			metric.WithAttributes(
				attribute.Int("syscall.id", int(event.Syscall)),
			),
		)
	}

	// 创建追踪 span
	if st.tracer != nil {
		_, span := st.tracer.Start(ctx, "syscall",
			trace.WithTimestamp(time.Unix(0, int64(event.Timestamp))),
			trace.WithAttributes(
				attribute.Int("syscall.id", int(event.Syscall)),
				attribute.Int("process.pid", int(event.PID)),
				attribute.Int("process.tid", int(event.TID)),
				attribute.Int64("syscall.return", event.RetVal),
				attribute.Float64("syscall.duration_ms", float64(event.Duration)/1000000.0),
			),
		)
		span.End(trace.WithTimestamp(time.Unix(0, int64(event.Timestamp+event.Duration))))
	}
}

// IsEnabled 检查是否启用
func (st *SyscallTracer) IsEnabled() bool {
	return st.enabled
}

// Note: 完整实现需要以下文件：
//
// 1. syscall.bpf.c - eBPF C 程序
//    - 附加到 sys_enter/sys_exit tracepoint
//    - 记录系统调用信息到 map
//
// 2. 使用 bpf2go 生成 Go 绑定：
//    //go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type syscall_event syscall ./programs/syscall.bpf.c
//
// 3. Makefile 添加生成命令：
//    generate-ebpf:
//        go generate ./pkg/observability/ebpf/...
//
// 参考实现：
// - https://github.com/cilium/ebpf/tree/main/examples
// - https://github.com/iovisor/bcc/blob/master/tools/syscount.py

