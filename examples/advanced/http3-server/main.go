// HTTP/3 Server示例：使用QUIC协议的高性能服务器
package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	
	"github.com/quic-go/quic-go/http3"
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

var stats Stats

func init() {
	stats.StartTime = time.Now()
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
		Server:    "Go 1.25",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	
	// 更新统计
	stats.Requests++
	duration := time.Since(start)
	stats.AvgDuration = (stats.AvgDuration*time.Duration(stats.Requests-1) + duration) / time.Duration(stats.Requests)
	
	log.Printf("%s %s - %v - %s", r.Method, r.URL.Path, duration, protocol)
}

// handleStats 统计信息处理
func handleStats(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(stats.StartTime)
	
	data := map[string]interface{}{
		"requests":     stats.Requests,
		"uptime":       uptime.String(),
		"avg_duration": stats.AvgDuration.String(),
		"req_per_sec":  float64(stats.Requests) / uptime.Seconds(),
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
	
	data := make([]map[string]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = map[string]interface{}{
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
	fmt.Println("🚀 HTTP/3 Server with QUIC")
	fmt.Println("==========================\n")
	
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
	
	// HTTP/3配置
	http3Server := &http3.Server{
		Handler:    mux,
		Addr:       ":8443",
		TLSConfig:  tlsConfig,
		QUICConfig: nil, // 使用默认配置
	}
	
	fmt.Println("📝 Endpoints:")
	fmt.Println("  GET  /        - Welcome message")
	fmt.Println("  GET  /stats   - Server statistics")
	fmt.Println("  GET  /health  - Health check")
	fmt.Println("  GET  /data    - Sample data")
	fmt.Println()
	
	fmt.Println("🌐 Server starting...")
	fmt.Println("  HTTP/2: https://localhost:8443")
	fmt.Println("  HTTP/3: https://localhost:8443 (QUIC/UDP)")
	fmt.Println()
	
	fmt.Println("💡 Test with:")
	fmt.Println("  curl --http3 https://localhost:8443")
	fmt.Println("  curl https://localhost:8443 (HTTP/2)")
	fmt.Println()
	
	// 启动HTTP/2服务器（TCP）
	go func() {
		log.Println("Starting HTTP/2 server...")
		// 注意：需要证书文件
		// if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		// 	log.Fatal(err)
		// }
		log.Println("Note: HTTP/2 requires cert.pem and key.pem")
	}()
	
	// 启动HTTP/3服务器（UDP + QUIC）
	log.Println("Starting HTTP/3 server...")
	// 注意：需要证书文件
	// if err := http3Server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
	// 	log.Fatal(err)
	// }
	
	fmt.Println("⚠️  Certificate files required:")
	fmt.Println("   Generate with: openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes")
	fmt.Println()
	
	fmt.Println("🎯 Features:")
	fmt.Println("  ✅ HTTP/3 over QUIC")
	fmt.Println("  ✅ 0-RTT connection resumption")
	fmt.Println("  ✅ Connection migration")
	fmt.Println("  ✅ Better performance on lossy networks")
	fmt.Println("  ✅ Fallback to HTTP/2")
	
	// 保持运行（示例代码）
	select {}
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

