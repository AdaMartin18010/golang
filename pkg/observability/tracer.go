package observability

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// 分布式追踪 - Distributed Tracing
// =============================================================================

// Span 追踪跨度
type Span struct {
	TraceID   string
	SpanID    string
	ParentID  string
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Tags      map[string]string
	Logs      []SpanLog
	Status    SpanStatus
	mu        sync.RWMutex
}

// SpanLog 跨度日志
type SpanLog struct {
	Timestamp time.Time
	Fields    map[string]interface{}
}

// SpanStatus 跨度状态
type SpanStatus struct {
	Code    StatusCode
	Message string
}

// StatusCode 状态码
type StatusCode int

const (
	StatusOK StatusCode = iota
	StatusError
	StatusUnknown
)

// Tracer 追踪器
type Tracer struct {
	serviceName string
	spans       sync.Map
	recorder    SpanRecorder
	sampler     Sampler
	idGenerator IDGenerator
}

// SpanRecorder 跨度记录器接口
type SpanRecorder interface {
	RecordSpan(*Span)
}

// Sampler 采样器接口
type Sampler interface {
	ShouldSample(string) bool
}

// IDGenerator ID生成器接口
type IDGenerator interface {
	GenerateTraceID() string
	GenerateSpanID() string
}

// NewTracer 创建追踪器
func NewTracer(serviceName string, recorder SpanRecorder, sampler Sampler) *Tracer {
	return &Tracer{
		serviceName: serviceName,
		recorder:    recorder,
		sampler:     sampler,
		idGenerator: &defaultIDGenerator{},
	}
}

// StartSpan 开始一个新跨度
func (t *Tracer) StartSpan(ctx context.Context, name string) (*Span, context.Context) {
	// 检查是否应该采样
	if t.sampler != nil && !t.sampler.ShouldSample(name) {
		return nil, ctx
	}

	span := &Span{
		TraceID:   t.idGenerator.GenerateTraceID(),
		SpanID:    t.idGenerator.GenerateSpanID(),
		Name:      name,
		StartTime: time.Now(),
		Tags:      make(map[string]string),
		Logs:      make([]SpanLog, 0),
		Status: SpanStatus{
			Code: StatusUnknown,
		},
	}

	// 如果有父跨度，继承TraceID
	if parentSpan := SpanFromContext(ctx); parentSpan != nil {
		span.TraceID = parentSpan.TraceID
		span.ParentID = parentSpan.SpanID
	}

	t.spans.Store(span.SpanID, span)

	// 将span放入context
	ctx = ContextWithSpan(ctx, span)

	return span, ctx
}

// Finish 结束跨度
func (s *Span) Finish() {
	if s == nil {
		return
	}

	s.mu.Lock()
	s.EndTime = time.Now()
	s.Duration = s.EndTime.Sub(s.StartTime)
	s.mu.Unlock()

	// 记录跨度（这里可以发送到后端）
	// 实际应用中会发送到Jaeger/Zipkin等
}

// SetTag 设置标签
func (s *Span) SetTag(key, value string) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.Tags[key] = value
}

// LogFields 记录日志字段
func (s *Span) LogFields(fields map[string]interface{}) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	log := SpanLog{
		Timestamp: time.Now(),
		Fields:    fields,
	}
	s.Logs = append(s.Logs, log)
}

// SetStatus 设置状态
func (s *Span) SetStatus(code StatusCode, message string) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Status = SpanStatus{
		Code:    code,
		Message: message,
	}
}

// SetError 设置错误
func (s *Span) SetError(err error) {
	if s == nil || err == nil {
		return
	}

	s.SetStatus(StatusError, err.Error())
	s.SetTag("error", "true")
	s.LogFields(map[string]interface{}{
		"event":        "error",
		"error.object": err.Error(),
	})
}

// =============================================================================
// Context支持
// =============================================================================

type spanContextKey struct{}

// ContextWithSpan 将Span放入Context
func ContextWithSpan(ctx context.Context, span *Span) context.Context {
	return context.WithValue(ctx, spanContextKey{}, span)
}

// SpanFromContext 从Context获取Span
func SpanFromContext(ctx context.Context) *Span {
	span, _ := ctx.Value(spanContextKey{}).(*Span)
	return span
}

// =============================================================================
// 默认实现
// =============================================================================

// defaultIDGenerator 默认ID生成器
type defaultIDGenerator struct {
	counter uint64
	mu      sync.Mutex
}

func (g *defaultIDGenerator) GenerateTraceID() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counter++
	return fmt.Sprintf("trace-%d-%d", time.Now().UnixNano(), g.counter)
}

func (g *defaultIDGenerator) GenerateSpanID() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counter++
	return fmt.Sprintf("span-%d-%d", time.Now().UnixNano(), g.counter)
}

// AlwaysSampler 总是采样
type AlwaysSampler struct{}

func (s *AlwaysSampler) ShouldSample(name string) bool {
	return true
}

// ProbabilitySampler 概率采样器
type ProbabilitySampler struct {
	probability float64
	counter     uint64
}

func NewProbabilitySampler(probability float64) *ProbabilitySampler {
	if probability < 0 {
		probability = 0
	}
	if probability > 1 {
		probability = 1
	}
	return &ProbabilitySampler{probability: probability}
}

func (s *ProbabilitySampler) ShouldSample(name string) bool {
	// 使用计数器来实现确定性采样（用于测试）
	// 实际生产中应使用crypto/rand
	s.counter++
	// 每N次采样一次，其中N = 1/probability
	if s.probability == 0 {
		return false
	}
	if s.probability == 1 {
		return true
	}
	threshold := uint64(1.0 / s.probability)
	return (s.counter % threshold) == 0
}

// InMemoryRecorder 内存记录器（用于测试）
type InMemoryRecorder struct {
	spans []Span
	mu    sync.Mutex
}

func NewInMemoryRecorder() *InMemoryRecorder {
	return &InMemoryRecorder{
		spans: make([]Span, 0),
	}
}

func (r *InMemoryRecorder) RecordSpan(span *Span) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.spans = append(r.spans, *span)
}

func (r *InMemoryRecorder) GetSpans() []Span {
	r.mu.Lock()
	defer r.mu.Unlock()
	spans := make([]Span, len(r.spans))
	copy(spans, r.spans)
	return spans
}

func (r *InMemoryRecorder) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.spans = r.spans[:0]
}

// =============================================================================
// 全局追踪器
// =============================================================================

var (
	globalTracer *Tracer
	tracerMu     sync.RWMutex
)

// SetGlobalTracer 设置全局追踪器
func SetGlobalTracer(tracer *Tracer) {
	tracerMu.Lock()
	defer tracerMu.Unlock()
	globalTracer = tracer
}

// GlobalTracer 获取全局追踪器
func GlobalTracer() *Tracer {
	tracerMu.RLock()
	defer tracerMu.RUnlock()
	return globalTracer
}

// StartSpan 从全局追踪器开始一个跨度
func StartSpan(ctx context.Context, name string) (*Span, context.Context) {
	tracer := GlobalTracer()
	if tracer == nil {
		// 没有全局追踪器，返回空跨度
		return nil, ctx
	}
	return tracer.StartSpan(ctx, name)
}
