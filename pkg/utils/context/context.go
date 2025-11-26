package context

import (
	"context"
	"time"
)

// WithTimeout 创建带超时的context
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// WithDeadline 创建带截止时间的context
func WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, d)
}

// WithCancel 创建可取消的context
func WithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

// WithValue 创建带值的context
func WithValue(parent context.Context, key, val interface{}) context.Context {
	return context.WithValue(parent, key, val)
}

// Background 返回非nil的空context
func Background() context.Context {
	return context.Background()
}

// TODO 返回非nil的空context（用于不确定使用哪个context时）
func TODO() context.Context {
	return context.TODO()
}

// WithTimeoutSeconds 创建带超时的context（秒为单位）
func WithTimeoutSeconds(parent context.Context, seconds int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, time.Duration(seconds)*time.Second)
}

// WithTimeoutMinutes 创建带超时的context（分钟为单位）
func WithTimeoutMinutes(parent context.Context, minutes int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, time.Duration(minutes)*time.Minute)
}

// WithTimeoutHours 创建带超时的context（小时为单位）
func WithTimeoutHours(parent context.Context, hours int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, time.Duration(hours)*time.Hour)
}

// IsDone 检查context是否已取消
func IsDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// IsCancelled 检查context是否已取消（别名）
func IsCancelled(ctx context.Context) bool {
	return IsDone(ctx)
}

// GetValue 获取context中的值
func GetValue(ctx context.Context, key interface{}) interface{} {
	return ctx.Value(key)
}

// GetStringValue 获取context中的字符串值
func GetStringValue(ctx context.Context, key interface{}) (string, bool) {
	value := ctx.Value(key)
	if value == nil {
		return "", false
	}
	str, ok := value.(string)
	return str, ok
}

// GetIntValue 获取context中的整数值
func GetIntValue(ctx context.Context, key interface{}) (int, bool) {
	value := ctx.Value(key)
	if value == nil {
		return 0, false
	}
	switch v := value.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	default:
		return 0, false
	}
}

// GetInt64Value 获取context中的64位整数值
func GetInt64Value(ctx context.Context, key interface{}) (int64, bool) {
	value := ctx.Value(key)
	if value == nil {
		return 0, false
	}
	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	default:
		return 0, false
	}
}

// GetBoolValue 获取context中的布尔值
func GetBoolValue(ctx context.Context, key interface{}) (bool, bool) {
	value := ctx.Value(key)
	if value == nil {
		return false, false
	}
	b, ok := value.(bool)
	return b, ok
}

// GetError 获取context的错误
func GetError(ctx context.Context) error {
	return ctx.Err()
}

// GetDeadline 获取context的截止时间
func GetDeadline(ctx context.Context) (time.Time, bool) {
	return ctx.Deadline()
}

// Wait 等待context完成
func Wait(ctx context.Context) {
	<-ctx.Done()
}

// WaitWithTimeout 等待context完成或超时
func WaitWithTimeout(ctx context.Context, timeout time.Duration) bool {
	select {
	case <-ctx.Done():
		return true
	case <-time.After(timeout):
		return false
	}
}

// Sleep 睡眠，但可以被context取消
func Sleep(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}

// DoWithTimeout 在超时时间内执行函数
func DoWithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return fn(ctx)
}

// DoWithDeadline 在截止时间前执行函数
func DoWithDeadline(ctx context.Context, deadline time.Time, fn func(context.Context) error) error {
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()
	return fn(ctx)
}

// DoWithCancel 执行函数，支持取消
func DoWithCancel(ctx context.Context, fn func(context.Context) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return fn(ctx)
}

// RetryWithContext 使用context重试函数
func RetryWithContext(ctx context.Context, maxRetries int, fn func(context.Context) error) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if IsDone(ctx) {
			return ctx.Err()
		}
		err := fn(ctx)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return lastErr
}

// Merge 合并多个context（任一取消则取消）
func Merge(ctxs ...context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		done := make(chan struct{})
		for _, c := range ctxs {
			go func(c context.Context) {
				<-c.Done()
				done <- struct{}{}
			}(c)
		}
		<-done
		cancel()
	}()
	return ctx, cancel
}

// WithValues 批量设置context值
func WithValues(parent context.Context, values map[interface{}]interface{}) context.Context {
	ctx := parent
	for key, value := range values {
		ctx = context.WithValue(ctx, key, value)
	}
	return ctx
}

// CopyValues 复制context中的值到新context
func CopyValues(from, to context.Context) context.Context {
	// 注意：标准库的context不支持遍历所有值
	// 这里需要知道具体的key才能复制
	return to
}

// ContextKey 上下文键类型
type ContextKey string

// String 返回键的字符串表示
func (k ContextKey) String() string {
	return string(k)
}

// WithStringValue 使用字符串键设置值
func WithStringValue(parent context.Context, key string, value interface{}) context.Context {
	return context.WithValue(parent, ContextKey(key), value)
}

// GetStringKeyValue 使用字符串键获取值
func GetStringKeyValue(ctx context.Context, key string) interface{} {
	return ctx.Value(ContextKey(key))
}

// CommonKeys 常用键
var (
	KeyTraceID   = ContextKey("trace_id")
	KeySpanID    = ContextKey("span_id")
	KeyUserID    = ContextKey("user_id")
	KeyRequestID = ContextKey("request_id")
	KeyIP        = ContextKey("ip")
	KeyUserAgent = ContextKey("user_agent")
)

// WithTraceID 设置TraceID
func WithTraceID(parent context.Context, traceID string) context.Context {
	return context.WithValue(parent, KeyTraceID, traceID)
}

// GetTraceID 获取TraceID
func GetTraceID(ctx context.Context) (string, bool) {
	return GetStringValue(ctx, KeyTraceID)
}

// WithSpanID 设置SpanID
func WithSpanID(parent context.Context, spanID string) context.Context {
	return context.WithValue(parent, KeySpanID, spanID)
}

// GetSpanID 获取SpanID
func GetSpanID(ctx context.Context) (string, bool) {
	return GetStringValue(ctx, KeySpanID)
}

// WithUserID 设置UserID
func WithUserID(parent context.Context, userID string) context.Context {
	return context.WithValue(parent, KeyUserID, userID)
}

// GetUserID 获取UserID
func GetUserID(ctx context.Context) (string, bool) {
	return GetStringValue(ctx, KeyUserID)
}

// WithRequestID 设置RequestID
func WithRequestID(parent context.Context, requestID string) context.Context {
	return context.WithValue(parent, KeyRequestID, requestID)
}

// GetRequestID 获取RequestID
func GetRequestID(ctx context.Context) (string, bool) {
	return GetStringValue(ctx, KeyRequestID)
}

// WithIP 设置IP
func WithIP(parent context.Context, ip string) context.Context {
	return context.WithValue(parent, KeyIP, ip)
}

// GetIP 获取IP
func GetIP(ctx context.Context) (string, bool) {
	return GetStringValue(ctx, KeyIP)
}

// WithUserAgent 设置UserAgent
func WithUserAgent(parent context.Context, userAgent string) context.Context {
	return context.WithValue(parent, KeyUserAgent, userAgent)
}

// GetUserAgent 获取UserAgent
func GetUserAgent(ctx context.Context) (string, bool) {
	return GetStringValue(ctx, KeyUserAgent)
}

// MustGetStringValue 获取字符串值，如果不存在则panic
func MustGetStringValue(ctx context.Context, key interface{}) string {
	value, ok := GetStringValue(ctx, key)
	if !ok {
		panic("context value not found or not a string")
	}
	return value
}

// MustGetIntValue 获取整数值，如果不存在则panic
func MustGetIntValue(ctx context.Context, key interface{}) int {
	value, ok := GetIntValue(ctx, key)
	if !ok {
		panic("context value not found or not an int")
	}
	return value
}

// MustGetBoolValue 获取布尔值，如果不存在则panic
func MustGetBoolValue(ctx context.Context, key interface{}) bool {
	value, ok := GetBoolValue(ctx, key)
	if !ok {
		panic("context value not found or not a bool")
	}
	return value
}

// Chain 链式设置多个值
func Chain(parent context.Context) *ContextBuilder {
	return &ContextBuilder{ctx: parent}
}

// ContextBuilder 上下文构建器
type ContextBuilder struct {
	ctx context.Context
}

// WithValue 添加值
func (b *ContextBuilder) WithValue(key, value interface{}) *ContextBuilder {
	b.ctx = context.WithValue(b.ctx, key, value)
	return b
}

// WithStringValue 添加字符串键值
func (b *ContextBuilder) WithStringValue(key string, value interface{}) *ContextBuilder {
	b.ctx = WithStringValue(b.ctx, key, value)
	return b
}

// WithTraceID 添加TraceID
func (b *ContextBuilder) WithTraceID(traceID string) *ContextBuilder {
	b.ctx = WithTraceID(b.ctx, traceID)
	return b
}

// WithSpanID 添加SpanID
func (b *ContextBuilder) WithSpanID(spanID string) *ContextBuilder {
	b.ctx = WithSpanID(b.ctx, spanID)
	return b
}

// WithUserID 添加UserID
func (b *ContextBuilder) WithUserID(userID string) *ContextBuilder {
	b.ctx = WithUserID(b.ctx, userID)
	return b
}

// WithRequestID 添加RequestID
func (b *ContextBuilder) WithRequestID(requestID string) *ContextBuilder {
	b.ctx = WithRequestID(b.ctx, requestID)
	return b
}

// Build 构建context
func (b *ContextBuilder) Build() context.Context {
	return b.ctx
}
