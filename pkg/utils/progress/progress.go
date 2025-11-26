package progress

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// ProgressBar 进度条
type ProgressBar struct {
	total      int64
	current    int64
	width      int
	showPercent bool
	showSpeed   bool
	showETA     bool
	startTime   time.Time
	lastUpdate  time.Time
	lastCurrent int64
	mu          sync.Mutex
	writer      io.Writer
	format      string
	prefix      string
	suffix      string
}

// ProgressBarOption 进度条选项
type ProgressBarOption func(*ProgressBar)

// WithWidth 设置进度条宽度
func WithWidth(width int) ProgressBarOption {
	return func(pb *ProgressBar) {
		pb.width = width
	}
}

// WithShowPercent 设置是否显示百分比
func WithShowPercent(show bool) ProgressBarOption {
	return func(pb *ProgressBar) {
		pb.showPercent = show
	}
}

// WithShowSpeed 设置是否显示速度
func WithShowSpeed(show bool) ProgressBarOption {
	return func(pb *ProgressBar) {
		pb.showSpeed = show
	}
}

// WithShowETA 设置是否显示预计剩余时间
func WithShowETA(show bool) ProgressBarOption {
	return func(pb *ProgressBar) {
		pb.showETA = show
	}
}

// WithWriter 设置输出写入器
func WithWriter(w io.Writer) ProgressBarOption {
	return func(pb *ProgressBar) {
		pb.writer = w
	}
}

// WithPrefix 设置前缀
func WithPrefix(prefix string) ProgressBarOption {
	return func(pb *ProgressBar) {
		pb.prefix = prefix
	}
}

// WithSuffix 设置后缀
func WithSuffix(suffix string) ProgressBarOption {
	return func(pb *ProgressBar) {
		pb.suffix = suffix
	}
}

// NewProgressBar 创建进度条
func NewProgressBar(total int64, opts ...ProgressBarOption) *ProgressBar {
	pb := &ProgressBar{
		total:       total,
		current:     0,
		width:       50,
		showPercent: true,
		showSpeed:   false,
		showETA:     false,
		startTime:   time.Now(),
		lastUpdate:  time.Now(),
		writer:      nil, // 默认使用标准输出
		format:      "",
		prefix:      "",
		suffix:      "",
	}

	for _, opt := range opts {
		opt(pb)
	}

	return pb
}

// Add 增加进度
func (pb *ProgressBar) Add(n int64) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.current += n
	if pb.current > pb.total {
		pb.current = pb.total
	}
	pb.update()
}

// Set 设置进度
func (pb *ProgressBar) Set(n int64) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.current = n
	if pb.current > pb.total {
		pb.current = pb.total
	}
	if pb.current < 0 {
		pb.current = 0
	}
	pb.update()
}

// Increment 增加1
func (pb *ProgressBar) Increment() {
	pb.Add(1)
}

// SetTotal 设置总数
func (pb *ProgressBar) SetTotal(total int64) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.total = total
	pb.update()
}

// Current 获取当前进度
func (pb *ProgressBar) Current() int64 {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	return pb.current
}

// Total 获取总数
func (pb *ProgressBar) Total() int64 {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	return pb.total
}

// Percent 获取百分比
func (pb *ProgressBar) Percent() float64 {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	if pb.total == 0 {
		return 0
	}
	return float64(pb.current) / float64(pb.total) * 100
}

// update 更新显示
func (pb *ProgressBar) update() {
	now := time.Now()
	pb.lastUpdate = now

	var buf strings.Builder

	// 前缀
	if pb.prefix != "" {
		buf.WriteString(pb.prefix)
		buf.WriteString(" ")
	}

	// 进度条
	percent := float64(pb.current) / float64(pb.total)
	filled := int(percent * float64(pb.width))
	empty := pb.width - filled

	buf.WriteString("[")
	buf.WriteString(strings.Repeat("=", filled))
	buf.WriteString(strings.Repeat(" ", empty))
	buf.WriteString("]")

	// 百分比
	if pb.showPercent {
		buf.WriteString(fmt.Sprintf(" %6.2f%%", percent*100))
	}

	// 当前/总数
	buf.WriteString(fmt.Sprintf(" %d/%d", pb.current, pb.total))

	// 速度
	if pb.showSpeed {
		elapsed := now.Sub(pb.startTime).Seconds()
		if elapsed > 0 {
			speed := float64(pb.current) / elapsed
			buf.WriteString(fmt.Sprintf(" %.2f/s", speed))
		}
	}

	// ETA
	if pb.showETA && pb.current > 0 && pb.current < pb.total {
		elapsed := now.Sub(pb.startTime).Seconds()
		if elapsed > 0 {
			rate := float64(pb.current) / elapsed
			if rate > 0 {
				remaining := float64(pb.total-pb.current) / rate
				buf.WriteString(fmt.Sprintf(" ETA: %.0fs", remaining))
			}
		}
	}

	// 后缀
	if pb.suffix != "" {
		buf.WriteString(" ")
		buf.WriteString(pb.suffix)
	}

	// 输出
	output := buf.String()
	if pb.writer != nil {
		fmt.Fprint(pb.writer, "\r"+output)
	} else {
		fmt.Print("\r" + output)
	}
}

// Finish 完成
func (pb *ProgressBar) Finish() {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.current = pb.total
	pb.update()
	if pb.writer != nil {
		fmt.Fprintln(pb.writer)
	} else {
		fmt.Println()
	}
}

// Reset 重置
func (pb *ProgressBar) Reset() {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.current = 0
	pb.startTime = time.Now()
	pb.lastUpdate = time.Now()
	pb.lastCurrent = 0
}

// SimpleProgressBar 简单进度条
type SimpleProgressBar struct {
	total   int64
	current int64
	mu      sync.Mutex
}

// NewSimpleProgressBar 创建简单进度条
func NewSimpleProgressBar(total int64) *SimpleProgressBar {
	return &SimpleProgressBar{
		total:   total,
		current: 0,
	}
}

// Add 增加进度
func (spb *SimpleProgressBar) Add(n int64) {
	spb.mu.Lock()
	defer spb.mu.Unlock()
	spb.current += n
	if spb.current > spb.total {
		spb.current = spb.total
	}
	spb.print()
}

// Set 设置进度
func (spb *SimpleProgressBar) Set(n int64) {
	spb.mu.Lock()
	defer spb.mu.Unlock()
	spb.current = n
	if spb.current > spb.total {
		spb.current = spb.total
	}
	if spb.current < 0 {
		spb.current = 0
	}
	spb.print()
}

// Increment 增加1
func (spb *SimpleProgressBar) Increment() {
	spb.Add(1)
}

// print 打印进度
func (spb *SimpleProgressBar) print() {
	percent := float64(spb.current) / float64(spb.total) * 100
	fmt.Printf("\rProgress: %d/%d (%.2f%%)", spb.current, spb.total, percent)
	if spb.current >= spb.total {
		fmt.Println()
	}
}

// Finish 完成
func (spb *SimpleProgressBar) Finish() {
	spb.mu.Lock()
	defer spb.mu.Unlock()
	spb.current = spb.total
	spb.print()
}

// Spinner 旋转器
type Spinner struct {
	chars    []string
	index    int
	message  string
	mu       sync.Mutex
	stop     chan struct{}
	stopped  bool
}

// NewSpinner 创建旋转器
func NewSpinner(message string) *Spinner {
	return &Spinner{
		chars:   []string{"|", "/", "-", "\\"},
		index:   0,
		message: message,
		stop:    make(chan struct{}),
		stopped: false,
	}
}

// Start 启动旋转器
func (s *Spinner) Start() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.mu.Lock()
				if s.stopped {
					s.mu.Unlock()
					return
				}
				char := s.chars[s.index]
				s.index = (s.index + 1) % len(s.chars)
				fmt.Printf("\r%s %s", s.message, char)
				s.mu.Unlock()
			case <-s.stop:
				return
			}
		}
	}()
}

// Stop 停止旋转器
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.stopped {
		s.stopped = true
		close(s.stop)
		fmt.Print("\r" + strings.Repeat(" ", len(s.message)+3) + "\r")
	}
}

// StopWithMessage 停止并显示消息
func (s *Spinner) StopWithMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.stopped {
		s.stopped = true
		close(s.stop)
		fmt.Printf("\r%s %s\n", s.message, message)
	}
}
