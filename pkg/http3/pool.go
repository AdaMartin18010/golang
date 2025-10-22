package main

import (
	"bytes"
	"sync"
)

// =============================================================================
// 对象池优化 - Object Pooling for Performance
// =============================================================================

// ResponsePool Response对象池
var ResponsePool = sync.Pool{
	New: func() interface{} {
		return &Response{}
	},
}

// GetResponse 从池中获取Response对象
func GetResponse() *Response {
	return ResponsePool.Get().(*Response)
}

// PutResponse 将Response对象放回池中
func PutResponse(resp *Response) {
	// 重置对象
	resp.Message = ""
	resp.Protocol = ""
	resp.Server = ""
	ResponsePool.Put(resp)
}

// BufferPool JSON编码buffer池
var BufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// GetBuffer 从池中获取buffer
func GetBuffer() *bytes.Buffer {
	return BufferPool.Get().(*bytes.Buffer)
}

// PutBuffer 将buffer放回池中
func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	BufferPool.Put(buf)
}

// DataItemPool 数据项对象池
var DataItemPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{})
	},
}

// GetDataItem 从池中获取数据项
func GetDataItem() map[string]interface{} {
	return DataItemPool.Get().(map[string]interface{})
}

// PutDataItem 将数据项放回池中
func PutDataItem(item map[string]interface{}) {
	// 清空map
	for k := range item {
		delete(item, k)
	}
	DataItemPool.Put(item)
}

// DataSlicePool 数据切片对象池
var DataSlicePool = sync.Pool{
	New: func() interface{} {
		// 预分配100个元素的切片
		return make([]map[string]interface{}, 0, 100)
	},
}

// GetDataSlice 从池中获取数据切片
func GetDataSlice() []map[string]interface{} {
	return DataSlicePool.Get().([]map[string]interface{})
}

// PutDataSlice 将数据切片放回池中
func PutDataSlice(slice []map[string]interface{}) {
	// 清空切片但保留容量
	slice = slice[:0]
	DataSlicePool.Put(slice)
}
