package sampling

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// Sampler 采样器接口
// 提供可配置的采样策略，用于控制数据收集和处理的频率
type Sampler interface {
	// ShouldSample 判断是否应该采样
	ShouldSample(ctx context.Context) bool

	// SampleRate 返回采样率
	SampleRate() float64

	// UpdateRate 更新采样率
	UpdateRate(rate float64) error
}

// AlwaysSampler 总是采样
type AlwaysSampler struct{}

func NewAlwaysSampler() Sampler {
	return &AlwaysSampler{}
}

func (s *AlwaysSampler) ShouldSample(ctx context.Context) bool {
	return true
}

func (s *AlwaysSampler) SampleRate() float64 {
	return 1.0
}

func (s *AlwaysSampler) UpdateRate(rate float64) error {
	return nil // 总是采样，忽略更新
}

// NeverSampler 从不采样
type NeverSampler struct{}

func NewNeverSampler() Sampler {
	return &NeverSampler{}
}

func (s *NeverSampler) ShouldSample(ctx context.Context) bool {
	return false
}

func (s *NeverSampler) SampleRate() float64 {
	return 0.0
}

func (s *NeverSampler) UpdateRate(rate float64) error {
	return nil // 从不采样，忽略更新
}

// ProbabilisticSampler 概率采样器
// 根据配置的概率决定是否采样
type ProbabilisticSampler struct {
	mu        sync.RWMutex
	rate      float64
	rand      *rand.Rand
	threshold float64
}

// NewProbabilisticSampler 创建概率采样器
// rate: 采样率，范围 [0.0, 1.0]
func NewProbabilisticSampler(rate float64) (Sampler, error) {
	if rate < 0.0 || rate > 1.0 {
		return nil, ErrInvalidSampleRate
	}

	return &ProbabilisticSampler{
		rate:      rate,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
		threshold: rate,
	}, nil
}

func (s *ProbabilisticSampler) ShouldSample(ctx context.Context) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.rand.Float64() < s.threshold
}

func (s *ProbabilisticSampler) SampleRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rate
}

func (s *ProbabilisticSampler) UpdateRate(rate float64) error {
	if rate < 0.0 || rate > 1.0 {
		return ErrInvalidSampleRate
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.rate = rate
	s.threshold = rate
	return nil
}

// RateLimitingSampler 速率限制采样器
// 在指定时间窗口内最多采样 N 次
type RateLimitingSampler struct {
	mu          sync.RWMutex
	maxPerSecond float64
	window      time.Duration
	lastReset   time.Time
	count       int64
}

// NewRateLimitingSampler 创建速率限制采样器
// maxPerSecond: 每秒最大采样次数
func NewRateLimitingSampler(maxPerSecond float64) (Sampler, error) {
	if maxPerSecond <= 0 {
		return nil, ErrInvalidSampleRate
	}

	return &RateLimitingSampler{
		maxPerSecond: maxPerSecond,
		window:       time.Second,
		lastReset:    time.Now(),
		count:        0,
	}, nil
}

func (s *RateLimitingSampler) ShouldSample(ctx context.Context) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if now.Sub(s.lastReset) >= s.window {
		s.count = 0
		s.lastReset = now
	}

	if float64(s.count) >= s.maxPerSecond {
		return false
	}

	s.count++
	return true
}

func (s *RateLimitingSampler) SampleRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.maxPerSecond
}

func (s *RateLimitingSampler) UpdateRate(rate float64) error {
	if rate <= 0 {
		return ErrInvalidSampleRate
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.maxPerSecond = rate
	return nil
}

// AdaptiveSampler 自适应采样器
// 根据系统负载动态调整采样率
type AdaptiveSampler struct {
	mu          sync.RWMutex
	baseRate    float64
	currentRate float64
	minRate     float64
	maxRate     float64
	rand        *rand.Rand
}

// NewAdaptiveSampler 创建自适应采样器
func NewAdaptiveSampler(baseRate, minRate, maxRate float64) (Sampler, error) {
	if baseRate < 0.0 || baseRate > 1.0 {
		return nil, ErrInvalidSampleRate
	}
	if minRate < 0.0 || minRate > 1.0 {
		return nil, ErrInvalidSampleRate
	}
	if maxRate < 0.0 || maxRate > 1.0 {
		return nil, ErrInvalidSampleRate
	}
	if minRate > maxRate {
		return nil, ErrInvalidSampleRate
	}

	return &AdaptiveSampler{
		baseRate:    baseRate,
		currentRate: baseRate,
		minRate:     minRate,
		maxRate:     maxRate,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func (s *AdaptiveSampler) ShouldSample(ctx context.Context) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.rand.Float64() < s.currentRate
}

func (s *AdaptiveSampler) SampleRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentRate
}

func (s *AdaptiveSampler) UpdateRate(rate float64) error {
	if rate < s.minRate || rate > s.maxRate {
		return ErrInvalidSampleRate
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentRate = rate
	return nil
}

// AdjustForLoad 根据负载调整采样率
func (s *AdaptiveSampler) AdjustForLoad(load float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 负载高时降低采样率，负载低时提高采样率
	if load > 0.8 {
		s.currentRate = s.minRate
	} else if load < 0.3 {
		s.currentRate = s.maxRate
	} else {
		// 线性插值
		factor := (0.8 - load) / 0.5
		s.currentRate = s.minRate + (s.maxRate-s.minRate)*factor
	}
}
