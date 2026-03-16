//go:build linux
// +build linux

package ebpf

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type syscall_event syscall ./programs/syscall.bpf.c -- -I/usr/include -I/usr/include/x86_64-linux-gnu

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
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

// SyscallTracer 系统调用追踪器
// 使用 Cilium eBPF 库实现真正的系统调用追踪
type SyscallTracer struct {
	tracer trace.Tracer
	meter  metric.Meter

	// eBPF 对象
	objs       *syscallObjects // bpf2go 生成的对象
	links      []link.Link     // tracepoint links
	perfReader *perf.Reader    // perf event reader
	ctx        context.Context
	cancel     context.CancelFunc
	enabled    bool

	// 指标
	syscallCounter metric.Int64Counter
	syscallLatency metric.Float64Histogram
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
	SyscallEvents    *ebpf.Map
	SyscallStats     *ebpf.Map
	SyscallStartTime *ebpf.Map
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

	// 非 Linux 系统不支持 eBPF
	if runtime.GOOS != "linux" {
		return &SyscallTracer{enabled: false}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	tracer := &SyscallTracer{
		tracer:  cfg.Tracer,
		meter:   cfg.Meter,
		ctx:     ctx,
		cancel:  cancel,
		enabled: true,
		links:   make([]link.Link, 0, 2),
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
func (st *SyscallTracer) loadPrograms(cfg SyscallTracerConfig) error {
	// 创建 eBPF 对象集合
	objs := &syscallObjects{}

	// 加载 eBPF 程序
	// 注意：这里使用 loadSyscallObjects 函数，它由 bpf2go 生成
	// 如果 bpf2go 生成的文件不存在，使用手动加载方式
	if err := st.loadObjects(objs); err != nil {
		return fmt.Errorf("failed to load syscall objects: %w", err)
	}
	st.objs = objs

	// 附加到 tracepoint: sys_enter
	enterLink, err := link.Tracepoint("raw_syscalls", "sys_enter", objs.Programs.TraceSyscallEnter, nil)
	if err != nil {
		objs.Close()
		return fmt.Errorf("failed to attach sys_enter tracepoint: %w", err)
	}
	st.links = append(st.links, enterLink)

	// 附加到 tracepoint: sys_exit
	exitLink, err := link.Tracepoint("raw_syscalls", "sys_exit", objs.Programs.TraceSyscallExit, nil)
	if err != nil {
		objs.Close()
		return fmt.Errorf("failed to attach sys_exit tracepoint: %w", err)
	}
	st.links = append(st.links, exitLink)

	// 创建 perf reader
	bufferSize := cfg.BufferSize
	if bufferSize == 0 {
		bufferSize = 4096 * 16 // 默认 64KB
	}

	reader, err := perf.NewReader(objs.Maps.SyscallEvents, bufferSize)
	if err != nil {
		objs.Close()
		return fmt.Errorf("failed to create perf reader: %w", err)
	}
	st.perfReader = reader

	return nil
}

// loadObjects 手动加载 eBPF 对象（当 bpf2go 生成文件不可用时使用）
func (st *SyscallTracer) loadObjects(objs *syscallObjects) error {
	// 这里应该调用 bpf2go 生成的 loadSyscallObjects 函数
	// 作为示例，返回错误提示用户需要生成代码
	// 实际使用时，取消下面注释并使用 bpf2go 生成的代码
	// return loadSyscallObjects(objs, nil)

	// 临时实现：返回错误，提示需要生成 eBPF 代码
	return errors.New("eBPF objects not generated. Run: go generate ./pkg/observability/ebpf")
}

// Start 启动追踪器
func (st *SyscallTracer) Start() error {
	if !st.enabled {
		return nil
	}

	// 启动 perf reader goroutine
	go st.readEvents()

	return nil
}

// Stop 停止追踪器
func (st *SyscallTracer) Stop() error {
	if !st.enabled {
		return nil
	}

	// 取消上下文
	if st.cancel != nil {
		st.cancel()
	}

	// 清理资源
	if st.perfReader != nil {
		st.perfReader.Close()
	}

	for _, l := range st.links {
		l.Close()
	}

	if st.objs != nil {
		st.objs.Close()
	}

	return nil
}

// readEvents 读取 eBPF 事件
// 这是实际的事件处理循环
func (st *SyscallTracer) readEvents() {
	if st.perfReader == nil {
		return
	}

	for {
		select {
		case <-st.ctx.Done():
			return
		default:
			// 从 perf buffer 读取事件
			record, err := st.perfReader.Read()
			if err != nil {
				if errors.Is(err, perf.ErrClosed) {
					return
				}
				// 记录错误但继续读取
				continue
			}

			// 解析事件
			var event SyscallEvent
			if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
				continue
			}

			// 处理事件
			st.handleSyscallEvent(st.ctx, &event)
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

// GetSyscallStats 获取系统调用统计
func (st *SyscallTracer) GetSyscallStats(ctx context.Context) (map[uint64]uint64, error) {
	if !st.enabled || st.objs == nil {
		return nil, nil
	}

	stats := make(map[uint64]uint64)

	// 从 eBPF map 读取统计信息
	var key uint64
	var value uint64

	iter := st.objs.Maps.SyscallStats.Iterate()
	for iter.Next(&key, &value) {
		stats[key] = value
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate syscall stats: %w", err)
	}

	return stats, nil
}

// IsEnabled 检查是否启用
func (st *SyscallTracer) IsEnabled() bool {
	return st.enabled
}

// Close 关闭追踪器（别名方法）
func (st *SyscallTracer) Close() error {
	return st.Stop()
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
