package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// =============================================================================
// 优化的处理器 - Optimized Handlers
// =============================================================================

// handleRootOptimized 优化后的根路径处理（使用对象池）
func handleRootOptimized(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// 从对象池获取Response
	resp := GetResponse()
	defer PutResponse(resp)

	// 检测协议
	protocol := "HTTP/1.1"
	if r.ProtoMajor == 3 {
		protocol = "HTTP/3"
	} else if r.ProtoMajor == 2 {
		protocol = "HTTP/2"
	}

	// 设置Response字段
	resp.Message = "Welcome to HTTP/3 Server!"
	resp.Timestamp = time.Now()
	resp.Protocol = protocol
	resp.Server = "Go 1.23+"

	// 使用buffer池进行JSON编码
	buf := GetBuffer()
	defer PutBuffer(buf)

	if err := json.NewEncoder(buf).Encode(resp); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf.Bytes())

	// 更新统计
	stats.Requests++
	duration := time.Since(start)
	stats.AvgDuration = (stats.AvgDuration*time.Duration(stats.Requests-1) + duration) / time.Duration(stats.Requests)
}

// handleStatsOptimized 优化后的统计信息处理（减少分配）
func handleStatsOptimized(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(stats.StartTime)

	// 使用buffer直接构建JSON，避免map分配
	buf := GetBuffer()
	defer PutBuffer(buf)

	buf.WriteString(`{"requests":`)
	buf.WriteString(strconv.FormatInt(stats.Requests, 10))
	buf.WriteString(`,"uptime":"`)
	buf.WriteString(uptime.String())
	buf.WriteString(`","avg_duration":"`)
	buf.WriteString(stats.AvgDuration.String())
	buf.WriteString(`","req_per_sec":`)
	buf.WriteString(strconv.FormatFloat(float64(stats.Requests)/uptime.Seconds(), 'f', 2, 64))
	buf.WriteString(`}`)

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf.Bytes())
}

// handleHealthOptimized 优化后的健康检查（静态响应）
func handleHealthOptimized(w http.ResponseWriter, r *http.Request) {
	// 预先分配的响应，避免每次都创建map
	buf := GetBuffer()
	defer PutBuffer(buf)

	// 手动构建JSON，避免反射开销
	buf.WriteString(`{"status":"healthy","time":"`)
	buf.WriteString(time.Now().Format(time.RFC3339))
	buf.WriteString(`"}`)

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf.Bytes())
}

// handleDataOptimized 优化后的数据处理（使用对象池和预分配）
func handleDataOptimized(w http.ResponseWriter, r *http.Request) {
	// 移除模拟延迟，专注于数据生成和编码性能
	// 如果需要延迟，应该异步处理或使用channel

	// 从对象池获取数据切片
	data := GetDataSlice()
	defer PutDataSlice(data)

	// 预分配所有数据项
	for i := 0; i < 100; i++ {
		item := GetDataItem()
		item["id"] = i
		item["value"] = float64(i) * 1.5
		item["name"] = fmt.Sprintf("Item-%d", i)
		data = append(data, item)
	}

	// 使用buffer池进行JSON编码
	buf := GetBuffer()
	defer PutBuffer(buf)

	if err := json.NewEncoder(buf).Encode(data); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		// 清理数据项
		for i := range data {
			PutDataItem(data[i])
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf.Bytes())

	// 将所有数据项放回池中
	for i := range data {
		PutDataItem(data[i])
	}
}

// handleDataOptimizedV2 进一步优化的数据处理（字符串拼接）
func handleDataOptimizedV2(w http.ResponseWriter, r *http.Request) {
	// 使用字符串拼接代替JSON编码，进一步减少开销
	buf := GetBuffer()
	defer PutBuffer(buf)

	buf.WriteString(`[`)
	for i := 0; i < 100; i++ {
		if i > 0 {
			buf.WriteString(`,`)
		}
		buf.WriteString(`{"id":`)
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString(`,"value":`)
		buf.WriteString(strconv.FormatFloat(float64(i)*1.5, 'f', 1, 64))
		buf.WriteString(`,"name":"Item-`)
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString(`"}`)
	}
	buf.WriteString(`]`)

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf.Bytes())
}

// =============================================================================
// 批量响应缓存 - Response Caching
// =============================================================================

var (
	healthResponseCache []byte
	cacheInitOnce       sync.Once
)

// initResponseCache 初始化响应缓存
func initResponseCache() {
	// 预生成健康检查响应（静态部分）
	// 注意：时间戳部分仍需动态生成
	healthResponseCache = []byte(`{"status":"healthy","time":"`)
}

// handleHealthCached 使用缓存的健康检查处理器
func handleHealthCached(w http.ResponseWriter, r *http.Request) {
	cacheInitOnce.Do(initResponseCache)

	buf := GetBuffer()
	defer PutBuffer(buf)

	buf.Write(healthResponseCache)
	buf.WriteString(time.Now().Format(time.RFC3339))
	buf.WriteString(`"}`)

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf.Bytes())
}
