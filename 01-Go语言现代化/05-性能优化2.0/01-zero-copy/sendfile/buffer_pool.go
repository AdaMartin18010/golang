package sendfile

import (
	"sync"
)

// BufferPool 缓冲区池
type BufferPool struct {
	pool sync.Pool
	size int
}

// NewBufferPool 创建新的缓冲区池
func NewBufferPool(maxSize, bufferSize int) *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, bufferSize)
			},
		},
		size: bufferSize,
	}
}

// Get 获取缓冲区
func (bp *BufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

// Put 归还缓冲区
func (bp *BufferPool) Put(buffer []byte) {
	// 重置缓冲区
	for i := range buffer {
		buffer[i] = 0
	}
	bp.pool.Put(buffer)
}

// GetSize 获取缓冲区大小
func (bp *BufferPool) GetSize() int {
	return bp.size
}
