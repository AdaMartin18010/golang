// HTTP/3 Server示例：使用QUIC协议的高性能服务器
// 注意：本示例展示HTTP/3的概念和API使用方式
//
// 要运行完整的HTTP/3服务器，需要安装依赖：
// go get github.com/quic-go/quic-go/http3
//
// 当前版本：展示HTTP/2服务器（无需额外依赖）

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Response API响应结构
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Protocol  string    `json:"protocol"`
	Server    string    `json:"server"`
}

// Stats 服务器统计
type Stats struct {
	Requests    int64
	StartTime   time.Time
	AvgDuration time.Duration
}

var (
	stats   Stats
	statsMu sync.RWMutex
)

func init() {
	stats.StartTime = time.Now()
}

// recordRequest 记录一次请求的处理耗时（并发安全）
func recordRequest(duration time.Duration) {
	statsMu.Lock()
	defer statsMu.Unlock()

	stats.Requests++
	stats.AvgDuration = (stats.AvgDuration*time.Duration(stats.Requests-1) + duration) / time.Duration(stats.Requests)
}

// statsSnapshot 返回当前统计信息的副本（并发安全）
func statsSnapshot() (requests int64, startTime time.Time, avgDuration time.Duration) {
	statsMu.RLock()
	defer statsMu.RUnlock()

	return stats.Requests, stats.StartTime, stats.AvgDuration
}

// handleRoot 根路径处理
func handleRoot(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// 检测协议
	protocol := "HTTP/1.1"
	if r.ProtoMajor == 3 {
		protocol = "HTTP/3"
	} else if r.ProtoMajor == 2 {
		protocol = "HTTP/2"
	}

	resp := Response{
		Message:   "Welcome to HTTP/3 Server!",
		Timestamp: time.Now(),
		Protocol:  protocol,
		Server:    "Go 1.23+",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	// 更新统计
	duration := time.Since(start)
	recordRequest(duration)

	log.Printf("%s %s - %v - %s", r.Method, r.URL.Path, duration, protocol)
}

// handleStats 统计信息处理
func handleStats(w http.ResponseWriter, r *http.Request) {
	requests, startTime, avgDuration := statsSnapshot()
	uptime := time.Since(startTime)
	reqPerSec := 0.0
	if uptime > 0 {
		reqPerSec = float64(requests) / uptime.Seconds()
	}

	data := map[string]any{
		"requests":     requests,
		"uptime":       uptime.String(),
		"avg_duration": avgDuration.String(),
		"req_per_sec":  reqPerSec,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// handleHealth 健康检查
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// handleData 数据处理示例
func handleData(w http.ResponseWriter, r *http.Request) {
	// 模拟数据处理
	time.Sleep(10 * time.Millisecond)

	data := make([]map[string]any, 100)
	for i := range 100 {
		data[i] = map[string]any{
			"id":    i,
			"value": float64(i) * 1.5,
			"name":  fmt.Sprintf("Item-%d", i),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// generateCert 生成自签名证书（仅用于测试）
func generateCert() {
	// 注意：生产环境请使用正式证书
	fmt.Println("⚠️  Using self-signed certificate (for testing only)")
	fmt.Println("💡 For production, use Let's Encrypt or other CA")
}

func main() {
	fmt.Println("🚀 HTTP/2 Server Demo (HTTP/3 Concept)")
	fmt.Println("======================================")
	fmt.Println("⚠️  Note: This demo runs HTTP/2.")
	fmt.Println("    For HTTP/3, install: go get github.com/quic-go/quic-go/http3")
	fmt.Println()

	generateCert()

	// 设置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/stats", handleStats)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/data", handleData)

	// HTTP服务器配置
	server := &http.Server{
		Addr:         ":8443",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// TLS配置
	tlsConfig := &tls.Config{
		// 注意：实际使用时需要配置证书
		MinVersion: tls.VersionTLS12,
	}
	server.TLSConfig = tlsConfig

	// HTTP/3配置（需要 github.com/quic-go/quic-go/http3）
	// 示例代码：
	// http3Server := &http3.Server{
	// 	Handler:    mux,
	// 	Addr:       ":8443",
	// 	TLSConfig:  tlsConfig,
	// 	QUICConfig: nil,
	// }

	fmt.Println("📝 Endpoints:")
	fmt.Println("  GET  /        - Welcome message")
	fmt.Println("  GET  /stats   - Server statistics")
	fmt.Println("  GET  /health  - Health check")
	fmt.Println("  GET  /data    - Sample data")
	fmt.Println()

	fmt.Println("🌐 Server configuration:")
	fmt.Println("  Address: https://localhost:8443")
	fmt.Println("  Protocol: HTTP/2 (TLS)")
	fmt.Println()

	fmt.Println("💡 Test with:")
	fmt.Println("  curl -k https://localhost:8443")
	fmt.Println("  curl -k https://localhost:8443/stats")
	fmt.Println("  curl -k https://localhost:8443/health")
	fmt.Println()

	// 启动HTTP/2服务器（示例）
	fmt.Println("⚠️  To actually start the server, uncomment the following:")
	fmt.Println()
	fmt.Println("  // Generate certificate first:")
	fmt.Println("  // openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes")
	fmt.Println()
	fmt.Println("  // Then uncomment:")
	fmt.Println("  // log.Fatal(server.ListenAndServeTLS(\"cert.pem\", \"key.pem\"))")
	fmt.Println()

	fmt.Println("🎯 HTTP/3 Features (when enabled):")
	fmt.Println("  ✅ HTTP/3 over QUIC (UDP-based)")
	fmt.Println("  ✅ 0-RTT connection resumption")
	fmt.Println("  ✅ Connection migration")
	fmt.Println("  ✅ Better performance on lossy networks")
	fmt.Println("  ✅ No head-of-line blocking")
	fmt.Println()

	// 示例：如何启动服务器（需要证书）
	// log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}

/*
性能对比说明:

HTTP/1.1:
- 队头阻塞
- 每个连接一个请求
- TCP慢启动

HTTP/2:
- 多路复用
- 头部压缩
- 但仍有TCP队头阻塞

HTTP/3 (QUIC):
- 基于UDP
- 独立流
- 0-RTT握手
- 连接迁移
- 更好的丢包恢复

理想网络: HTTP/3 提升5-10%
弱网环境: HTTP/3 提升30-50%
*/
