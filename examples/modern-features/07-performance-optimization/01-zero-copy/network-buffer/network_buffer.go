package network_buffer

import (
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Buffer 网络缓冲区
type Buffer struct {
	data   []byte
	offset int
	length int
	pool   *BufferPool
}

// NewBuffer 创建新缓冲区
func NewBuffer(size int) *Buffer {
	return &Buffer{
		data:   make([]byte, size),
		offset: 0,
		length: 0,
	}
}

// Data 获取缓冲区数据
func (b *Buffer) Data() []byte {
	return b.data[b.offset : b.offset+b.length]
}

// Capacity 获取缓冲区容量
func (b *Buffer) Capacity() int {
	return len(b.data)
}

// Length 获取数据长度
func (b *Buffer) Length() int {
	return b.length
}

// Available 获取可用空间
func (b *Buffer) Available() int {
	return len(b.data) - b.offset - b.length
}

// Reset 重置缓冲区
func (b *Buffer) Reset() {
	b.offset = 0
	b.length = 0
}

// ReadFrom 从Reader读取数据
func (b *Buffer) ReadFrom(r io.Reader) (int, error) {
	if b.Available() == 0 {
		return 0, io.ErrShortBuffer
	}

	n, err := r.Read(b.data[b.offset+b.length:])
	if n > 0 {
		b.length += n
	}
	return n, err
}

// WriteTo 写入到Writer
func (b *Buffer) WriteTo(w io.Writer) (int, error) {
	if b.length == 0 {
		return 0, nil
	}

	n, err := w.Write(b.data[b.offset : b.offset+b.length])
	if n > 0 {
		b.offset += n
		b.length -= n
	}
	return n, err
}

// Consume 消费数据
func (b *Buffer) Consume(n int) {
	if n > b.length {
		n = b.length
	}
	b.offset += n
	b.length -= n
}

// Release 释放缓冲区
func (b *Buffer) Release() {
	if b.pool != nil {
		b.pool.Put(b)
	}
}

// BufferPool 缓冲区池
type BufferPool struct {
	pool sync.Pool
	size int
}

// NewBufferPool 创建新的缓冲区池
func NewBufferPool(bufferSize int) *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return NewBuffer(bufferSize)
			},
		},
		size: bufferSize,
	}
}

// Get 获取缓冲区
func (bp *BufferPool) Get() *Buffer {
	buffer := bp.pool.Get().(*Buffer)
	buffer.pool = bp
	buffer.Reset()
	return buffer
}

// Put 归还缓冲区
func (bp *BufferPool) Put(buffer *Buffer) {
	if buffer.pool == bp {
		buffer.Reset()
		buffer.pool = nil
		bp.pool.Put(buffer)
	}
}

// BufferChain 缓冲区链
type BufferChain struct {
	head *BufferNode
	tail *BufferNode
	size int64
	mu   sync.RWMutex
}

// BufferNode 缓冲区节点
type BufferNode struct {
	buffer *Buffer
	next   *BufferNode
}

// NewBufferChain 创建新的缓冲区链
func NewBufferChain() *BufferChain {
	return &BufferChain{}
}

// Append 添加缓冲区
func (bc *BufferChain) Append(buffer *Buffer) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	node := &BufferNode{buffer: buffer}

	if bc.head == nil {
		bc.head = node
		bc.tail = node
	} else {
		bc.tail.next = node
		bc.tail = node
	}

	atomic.AddInt64(&bc.size, int64(buffer.Length()))
}

// Prepend 前置添加缓冲区
func (bc *BufferChain) Prepend(buffer *Buffer) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	node := &BufferNode{buffer: buffer, next: bc.head}
	bc.head = node

	if bc.tail == nil {
		bc.tail = node
	}

	atomic.AddInt64(&bc.size, int64(buffer.Length()))
}

// Remove 移除缓冲区
func (bc *BufferChain) Remove(buffer *Buffer) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if bc.head == nil {
		return
	}

	if bc.head.buffer == buffer {
		bc.head = bc.head.next
		if bc.head == nil {
			bc.tail = nil
		}
		atomic.AddInt64(&bc.size, -int64(buffer.Length()))
		return
	}

	current := bc.head
	for current.next != nil {
		if current.next.buffer == buffer {
			current.next = current.next.next
			if current.next == nil {
				bc.tail = current
			}
			atomic.AddInt64(&bc.size, -int64(buffer.Length()))
			return
		}
		current = current.next
	}
}

// Size 获取总大小
func (bc *BufferChain) Size() int64 {
	return atomic.LoadInt64(&bc.size)
}

// Clear 清空缓冲区链
func (bc *BufferChain) Clear() {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.head = nil
	bc.tail = nil
	atomic.StoreInt64(&bc.size, 0)
}

// ReadFrom 从Reader读取数据到缓冲区链
func (bc *BufferChain) ReadFrom(r io.Reader, pool *BufferPool) (int64, error) {
	var total int64

	for {
		buffer := pool.Get()
		n, err := buffer.ReadFrom(r)
		total += int64(n)

		if n > 0 {
			bc.Append(buffer)
		} else {
			buffer.Release()
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return total, err
		}
	}

	return total, nil
}

// WriteTo 将缓冲区链写入Writer
func (bc *BufferChain) WriteTo(w io.Writer) (int64, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	var total int64
	current := bc.head

	for current != nil {
		n, err := current.buffer.WriteTo(w)
		total += int64(n)

		if err != nil {
			return total, err
		}

		current = current.next
	}

	return total, nil
}

// ZeroCopyBuffer 零拷贝缓冲区
type ZeroCopyBuffer struct {
	data   []byte
	offset int
	length int
	refs   int32
	pool   *ZeroCopyBufferPool
}

// NewZeroCopyBuffer 创建新的零拷贝缓冲区
func NewZeroCopyBuffer(size int) *ZeroCopyBuffer {
	return &ZeroCopyBuffer{
		data:   make([]byte, size),
		offset: 0,
		length: 0,
		refs:   1,
	}
}

// Data 获取数据
func (zcb *ZeroCopyBuffer) Data() []byte {
	return zcb.data[zcb.offset : zcb.offset+zcb.length]
}

// Slice 切片操作（零拷贝）
func (zcb *ZeroCopyBuffer) Slice(start, end int) *ZeroCopyBuffer {
	if start < 0 || end > zcb.length || start > end {
		return nil
	}

	atomic.AddInt32(&zcb.refs, 1)

	return &ZeroCopyBuffer{
		data:   zcb.data,
		offset: zcb.offset + start,
		length: end - start,
		refs:   1,
		pool:   zcb.pool,
	}
}

// Release 释放缓冲区
func (zcb *ZeroCopyBuffer) Release() {
	if atomic.AddInt32(&zcb.refs, -1) == 0 {
		if zcb.pool != nil {
			zcb.pool.Put(zcb)
		}
	}
}

// ZeroCopyBufferPool 零拷贝缓冲区池
type ZeroCopyBufferPool struct {
	pool sync.Pool
	size int
}

// NewZeroCopyBufferPool 创建新的零拷贝缓冲区池
func NewZeroCopyBufferPool(bufferSize int) *ZeroCopyBufferPool {
	return &ZeroCopyBufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return NewZeroCopyBuffer(bufferSize)
			},
		},
		size: bufferSize,
	}
}

// Get 获取零拷贝缓冲区
func (zcbp *ZeroCopyBufferPool) Get() *ZeroCopyBuffer {
	buffer := zcbp.pool.Get().(*ZeroCopyBuffer)
	buffer.pool = zcbp
	buffer.offset = 0
	buffer.length = 0
	buffer.refs = 1
	return buffer
}

// Put 归还零拷贝缓冲区
func (zcbp *ZeroCopyBufferPool) Put(buffer *ZeroCopyBuffer) {
	if buffer.pool == zcbp && atomic.LoadInt32(&buffer.refs) == 0 {
		buffer.pool = nil
		zcbp.pool.Put(buffer)
	}
}

// RingBuffer 环形缓冲区
type RingBuffer struct {
	data   []byte
	read   int
	write  int
	size   int
	length int
	mu     sync.RWMutex
}

// NewRingBuffer 创建新的环形缓冲区
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data:   make([]byte, size),
		read:   0,
		write:  0,
		size:   size,
		length: 0,
	}
}

// Read 读取数据
func (rb *RingBuffer) Read(p []byte) (int, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.length == 0 {
		return 0, io.EOF
	}

	n := len(p)
	if n > rb.length {
		n = rb.length
	}

	// 计算可读取的数据
	available := rb.size - rb.read
	if n > available {
		n = available
	}

	// 复制数据
	copy(p, rb.data[rb.read:rb.read+n])
	rb.read = (rb.read + n) % rb.size
	rb.length -= n

	return n, nil
}

// Write 写入数据
func (rb *RingBuffer) Write(p []byte) (int, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	n := len(p)
	if n > rb.size-rb.length {
		n = rb.size - rb.length
	}

	if n == 0 {
		return 0, io.ErrShortWrite
	}

	// 计算可写入的空间
	available := rb.size - rb.write
	if n > available {
		n = available
	}

	// 复制数据
	copy(rb.data[rb.write:], p[:n])
	rb.write = (rb.write + n) % rb.size
	rb.length += n

	return n, nil
}

// Length 获取数据长度
func (rb *RingBuffer) Length() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.length
}

// Available 获取可用空间
func (rb *RingBuffer) Available() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.size - rb.length
}

// Reset 重置环形缓冲区
func (rb *RingBuffer) Reset() {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.read = 0
	rb.write = 0
	rb.length = 0
}

// CPU特性检测
func hasSSE2() bool {
	return runtime.GOARCH == "amd64"
}

func hasAVX2() bool {
	return runtime.GOARCH == "amd64"
}

// 内存对齐辅助函数
func AlignedBuffer(size int) []byte {
	// 确保内存对齐到32字节边界
	aligned := make([]byte, size+32)

	// 找到对齐的起始位置
	ptr := uintptr(unsafe.Pointer(&aligned[0]))
	offset := (32 - ptr%32)

	return aligned[offset : offset+uintptr(size)]
}
