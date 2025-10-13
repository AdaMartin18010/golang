package memorypool

import (
	"sync"
	"time"
)

// Resetter 定义可重置对象的接口
type Resetter interface {
	Reset()
}

// ObjectPool 通用对象池实现
type ObjectPool struct {
	pool     chan interface{}
	new      func() interface{}
	reset    func(interface{})
	capacity int
	created  int64
	reused   int64
	mu       sync.RWMutex
}

// NewObjectPool 创建新的对象池
func NewObjectPool(capacity int, new func() interface{}, reset func(interface{})) *ObjectPool {
	return &ObjectPool{
		pool:     make(chan interface{}, capacity),
		new:      new,
		reset:    reset,
		capacity: capacity,
	}
}

// Get 从池中获取对象
func (p *ObjectPool) Get() interface{} {
	select {
	case obj := <-p.pool:
		// 从池中获取到对象，重置并返回
		if p.reset != nil {
			p.reset(obj)
		}
		p.mu.Lock()
		p.reused++
		p.mu.Unlock()
		return obj
	default:
		// 池为空，创建新对象
		obj := p.new()
		p.mu.Lock()
		p.created++
		p.mu.Unlock()
		return obj
	}
}

// Put 将对象放回池中
func (p *ObjectPool) Put(obj interface{}) {
	if obj == nil {
		return
	}

	select {
	case p.pool <- obj:
		// 成功放回池中
	default:
		// 池已满，丢弃对象
		// 在实际应用中，可能需要记录这种情况
	}
}

// Stats 获取池的统计信息
func (p *ObjectPool) Stats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	total := p.created + p.reused
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(p.reused) / float64(total) * 100
	}

	return PoolStats{
		Capacity:  p.capacity,
		Created:   p.created,
		Reused:    p.reused,
		HitRate:   hitRate,
		PoolSize:  len(p.pool),
		Timestamp: time.Now(),
	}
}

// PoolStats 池统计信息
type PoolStats struct {
	Capacity  int       `json:"capacity"`
	Created   int64     `json:"created"`
	Reused    int64     `json:"reused"`
	HitRate   float64   `json:"hit_rate"`
	PoolSize  int       `json:"pool_size"`
	Timestamp time.Time `json:"timestamp"`
}

// BufferPool 字节缓冲区池
type BufferPool struct {
	*ObjectPool
}

// NewBufferPool 创建字节缓冲区池
func NewBufferPool(capacity int, bufferSize int) *BufferPool {
	newFunc := func() interface{} {
		return make([]byte, 0, bufferSize)
	}

	resetFunc := func(obj interface{}) {
		if buf, ok := obj.([]byte); ok {
			// 重置缓冲区，保持容量但清空内容
			buf = buf[:0]
		}
	}

	return &BufferPool{
		ObjectPool: NewObjectPool(capacity, newFunc, resetFunc),
	}
}

// GetBuffer 获取字节缓冲区
func (p *BufferPool) GetBuffer() []byte {
	obj := p.Get()
	if buf, ok := obj.([]byte); ok {
		return buf
	}
	// 这种情况不应该发生，但为了安全起见
	return make([]byte, 0)
}

// PutBuffer 归还字节缓冲区
func (p *BufferPool) PutBuffer(buf []byte) {
	p.Put(buf)
}

// RequestPool HTTP请求对象池
type RequestPool struct {
	*ObjectPool
}

// HTTPRequest 简化的HTTP请求结构
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

// NewRequestPool 创建HTTP请求池
func NewRequestPool(capacity int) *RequestPool {
	newFunc := func() interface{} {
		return &HTTPRequest{
			Headers: make(map[string]string),
			Body:    make([]byte, 0),
		}
	}

	resetFunc := func(obj interface{}) {
		if req, ok := obj.(*HTTPRequest); ok {
			req.Method = ""
			req.URL = ""
			// 清空headers但保持map结构
			for k := range req.Headers {
				delete(req.Headers, k)
			}
			req.Body = req.Body[:0]
		}
	}

	return &RequestPool{
		ObjectPool: NewObjectPool(capacity, newFunc, resetFunc),
	}
}

// GetRequest 获取HTTP请求对象
func (p *RequestPool) GetRequest() *HTTPRequest {
	obj := p.Get()
	if req, ok := obj.(*HTTPRequest); ok {
		return req
	}
	return &HTTPRequest{
		Headers: make(map[string]string),
		Body:    make([]byte, 0),
	}
}

// PutRequest 归还HTTP请求对象
func (p *RequestPool) PutRequest(req *HTTPRequest) {
	p.Put(req)
}

// ResponsePool HTTP响应对象池
type ResponsePool struct {
	*ObjectPool
}

// HTTPResponse 简化的HTTP响应结构
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// NewResponsePool 创建HTTP响应池
func NewResponsePool(capacity int) *ResponsePool {
	newFunc := func() interface{} {
		return &HTTPResponse{
			Headers: make(map[string]string),
			Body:    make([]byte, 0),
		}
	}

	resetFunc := func(obj interface{}) {
		if resp, ok := obj.(*HTTPResponse); ok {
			resp.StatusCode = 0
			// 清空headers但保持map结构
			for k := range resp.Headers {
				delete(resp.Headers, k)
			}
			resp.Body = resp.Body[:0]
		}
	}

	return &ResponsePool{
		ObjectPool: NewObjectPool(capacity, newFunc, resetFunc),
	}
}

// GetResponse 获取HTTP响应对象
func (p *ResponsePool) GetResponse() *HTTPResponse {
	obj := p.Get()
	if resp, ok := obj.(*HTTPResponse); ok {
		return resp
	}
	return &HTTPResponse{
		Headers: make(map[string]string),
		Body:    make([]byte, 0),
	}
}

// PutResponse 归还HTTP响应对象
func (p *ResponsePool) PutResponse(resp *HTTPResponse) {
	p.Put(resp)
}
