package timer

import (
	"sync"
	"time"
)

// Timer 定时器接口
type Timer interface {
	Start()
	Stop()
	Reset(duration time.Duration)
	IsRunning() bool
}

// SimpleTimer 简单定时器
type SimpleTimer struct {
	duration time.Duration
	ticker   *time.Ticker
	callback func()
	mu       sync.RWMutex
	running  bool
	stop     chan struct{}
}

// NewSimpleTimer 创建简单定时器
func NewSimpleTimer(duration time.Duration, callback func()) *SimpleTimer {
	return &SimpleTimer{
		duration: duration,
		callback: callback,
		stop:     make(chan struct{}),
	}
}

// Start 启动定时器
func (st *SimpleTimer) Start() {
	st.mu.Lock()
	defer st.mu.Unlock()
	
	if st.running {
		return
	}
	
	st.running = true
	st.ticker = time.NewTicker(st.duration)
	
	go func() {
		for {
			select {
			case <-st.ticker.C:
				st.callback()
			case <-st.stop:
				return
			}
		}
	}()
}

// Stop 停止定时器
func (st *SimpleTimer) Stop() {
	st.mu.Lock()
	defer st.mu.Unlock()
	
	if !st.running {
		return
	}
	
	st.running = false
	if st.ticker != nil {
		st.ticker.Stop()
	}
	close(st.stop)
	st.stop = make(chan struct{})
}

// Reset 重置定时器
func (st *SimpleTimer) Reset(duration time.Duration) {
	st.mu.Lock()
	defer st.mu.Unlock()
	
	st.duration = duration
	
	if st.running {
		if st.ticker != nil {
			st.ticker.Stop()
		}
		st.ticker = time.NewTicker(duration)
	}
}

// IsRunning 检查是否运行中
func (st *SimpleTimer) IsRunning() bool {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.running
}

// OneShotTimer 一次性定时器
type OneShotTimer struct {
	duration time.Duration
	timer    *time.Timer
	callback func()
	mu       sync.RWMutex
	running  bool
}

// NewOneShotTimer 创建一次性定时器
func NewOneShotTimer(duration time.Duration, callback func()) *OneShotTimer {
	return &OneShotTimer{
		duration: duration,
		callback: callback,
	}
}

// Start 启动定时器
func (ot *OneShotTimer) Start() {
	ot.mu.Lock()
	defer ot.mu.Unlock()
	
	if ot.running {
		return
	}
	
	ot.running = true
	ot.timer = time.AfterFunc(ot.duration, func() {
		ot.mu.Lock()
		ot.running = false
		ot.mu.Unlock()
		ot.callback()
	})
}

// Stop 停止定时器
func (ot *OneShotTimer) Stop() {
	ot.mu.Lock()
	defer ot.mu.Unlock()
	
	if !ot.running {
		return
	}
	
	ot.running = false
	if ot.timer != nil {
		ot.timer.Stop()
	}
}

// Reset 重置定时器
func (ot *OneShotTimer) Reset(duration time.Duration) {
	ot.mu.Lock()
	defer ot.mu.Unlock()
	
	ot.duration = duration
	
	if ot.running && ot.timer != nil {
		ot.timer.Stop()
	}
	
	ot.running = true
	ot.timer = time.AfterFunc(duration, func() {
		ot.mu.Lock()
		ot.running = false
		ot.mu.Unlock()
		ot.callback()
	})
}

// IsRunning 检查是否运行中
func (ot *OneShotTimer) IsRunning() bool {
	ot.mu.RLock()
	defer ot.mu.RUnlock()
	return ot.running
}

// DebounceTimer 防抖定时器
type DebounceTimer struct {
	duration time.Duration
	timer    *time.Timer
	callback func()
	mu       sync.Mutex
}

// NewDebounceTimer 创建防抖定时器
func NewDebounceTimer(duration time.Duration, callback func()) *DebounceTimer {
	return &DebounceTimer{
		duration: duration,
		callback: callback,
	}
}

// Trigger 触发防抖
func (dt *DebounceTimer) Trigger() {
	dt.mu.Lock()
	defer dt.mu.Unlock()
	
	if dt.timer != nil {
		dt.timer.Stop()
	}
	
	dt.timer = time.AfterFunc(dt.duration, func() {
		dt.mu.Lock()
		dt.timer = nil
		dt.mu.Unlock()
		dt.callback()
	})
}

// Cancel 取消防抖
func (dt *DebounceTimer) Cancel() {
	dt.mu.Lock()
	defer dt.mu.Unlock()
	
	if dt.timer != nil {
		dt.timer.Stop()
		dt.timer = nil
	}
}

// ThrottleTimer 节流定时器
type ThrottleTimer struct {
	duration time.Duration
	lastTime time.Time
	callback func()
	mu       sync.Mutex
}

// NewThrottleTimer 创建节流定时器
func NewThrottleTimer(duration time.Duration, callback func()) *ThrottleTimer {
	return &ThrottleTimer{
		duration: duration,
		lastTime: time.Time{},
		callback: callback,
	}
}

// Trigger 触发节流
func (tt *ThrottleTimer) Trigger() {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	
	now := time.Now()
	if now.Sub(tt.lastTime) >= tt.duration {
		tt.lastTime = now
		tt.callback()
	}
}

// Reset 重置节流定时器
func (tt *ThrottleTimer) Reset() {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	tt.lastTime = time.Time{}
}

// IntervalTimer 间隔定时器
type IntervalTimer struct {
	interval  time.Duration
	callback  func()
	ticker    *time.Ticker
	mu        sync.RWMutex
	running   bool
	stop      chan struct{}
	execCount int64
}

// NewIntervalTimer 创建间隔定时器
func NewIntervalTimer(interval time.Duration, callback func()) *IntervalTimer {
	return &IntervalTimer{
		interval: interval,
		callback: callback,
		stop:     make(chan struct{}),
	}
}

// Start 启动定时器
func (it *IntervalTimer) Start() {
	it.mu.Lock()
	defer it.mu.Unlock()
	
	if it.running {
		return
	}
	
	it.running = true
	it.ticker = time.NewTicker(it.interval)
	
	go func() {
		for {
			select {
			case <-it.ticker.C:
				it.mu.Lock()
				it.execCount++
				it.mu.Unlock()
				it.callback()
			case <-it.stop:
				return
			}
		}
	}()
}

// Stop 停止定时器
func (it *IntervalTimer) Stop() {
	it.mu.Lock()
	defer it.mu.Unlock()
	
	if !it.running {
		return
	}
	
	it.running = false
	if it.ticker != nil {
		it.ticker.Stop()
	}
	close(it.stop)
	it.stop = make(chan struct{})
}

// Reset 重置定时器
func (it *IntervalTimer) Reset(interval time.Duration) {
	it.mu.Lock()
	defer it.mu.Unlock()
	
	it.interval = interval
	
	if it.running {
		if it.ticker != nil {
			it.ticker.Stop()
		}
		it.ticker = time.NewTicker(interval)
	}
}

// IsRunning 检查是否运行中
func (it *IntervalTimer) IsRunning() bool {
	it.mu.RLock()
	defer it.mu.RUnlock()
	return it.running
}

// ExecutionCount 获取执行次数
func (it *IntervalTimer) ExecutionCount() int64 {
	it.mu.RLock()
	defer it.mu.RUnlock()
	return it.execCount
}

// ResetCount 重置执行次数
func (it *IntervalTimer) ResetCount() {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.execCount = 0
}

// After 延迟执行
func After(duration time.Duration, callback func()) *time.Timer {
	return time.AfterFunc(duration, callback)
}

// Every 定期执行
func Every(interval time.Duration, callback func()) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			callback()
		}
	}()
	return ticker
}

// Schedule 调度执行
func Schedule(delay, interval time.Duration, callback func()) *time.Ticker {
	timer := time.NewTimer(delay)
	ticker := time.NewTicker(interval)
	
	go func() {
		<-timer.C
		callback() // 首次执行
		
		for range ticker.C {
			callback()
		}
	}()
	
	return ticker
}

