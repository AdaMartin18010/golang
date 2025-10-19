package sendfile

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FileServer 高性能文件服务器
type FileServer struct {
	rootDir    string
	bufferPool *BufferPool
	metrics    *Metrics
}

// NewFileServer 创建新的文件服务器
func NewFileServer(rootDir string) *FileServer {
	return &FileServer{
		rootDir:    rootDir,
		bufferPool: NewBufferPool(1024, 4096), // 1024个缓冲区，每个4KB
		metrics:    NewMetrics(),
	}
}

// ServeHTTP 处理HTTP请求
func (s *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		s.metrics.RecordRequest(time.Since(start))
	}()

	// 构建文件路径
	filePath := filepath.Join(s.rootDir, r.URL.Path)

	// 安全检查：防止目录遍历攻击
	if !strings.HasPrefix(filepath.Clean(filePath), s.rootDir) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// 如果是目录，返回403
	if fileInfo.IsDir() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 设置响应头
	w.Header().Set("Content-Type", getContentType(filePath))
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	w.Header().Set("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// 处理Range请求
	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		s.handleRangeRequest(w, r, file, fileInfo, rangeHeader)
		return
	}

	// 使用sendfile优化传输
	if err := s.sendFileOptimized(w, file, fileInfo.Size()); err != nil {
		// 如果sendfile失败，回退到标准方法
		s.sendFileStandard(w, file)
	}
}

// sendFileOptimized 使用sendfile系统调用优化文件传输
func (s *FileServer) sendFileOptimized(w http.ResponseWriter, file *os.File, size int64) error {
	// 获取底层连接
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		return err
	}
	defer conn.Close()

	// 发送HTTP响应头
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
		"Content-Length: %d\r\n"+
		"Content-Type: %s\r\n"+
		"Connection: close\r\n"+
		"\r\n", size, getContentType(file.Name()))

	if _, err := conn.Write([]byte(response)); err != nil {
		return err
	}

	// 注意：真正的sendfile优化在不同平台上的实现不同
	// 在Linux上使用syscall.Sendfile
	// 在Windows上需要使用TransmitFile API
	// 这里使用标准库的io.Copy作为跨平台的回退方案

	buffer := make([]byte, 32*1024)
	written, err := io.CopyBuffer(conn, file, buffer)
	if err != nil {
		return err
	}

	s.metrics.RecordBytesTransferred(int64(written))
	return nil
}

// sendFileStandard 标准文件传输方法（回退方案）
func (s *FileServer) sendFileStandard(w http.ResponseWriter, file *os.File) {
	// 获取缓冲区
	buffer := s.bufferPool.Get()
	defer s.bufferPool.Put(buffer)

	written, err := io.CopyBuffer(w, file, buffer)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.metrics.RecordBytesTransferred(written)
}

// handleRangeRequest 处理Range请求
func (s *FileServer) handleRangeRequest(w http.ResponseWriter, r *http.Request, file *os.File, fileInfo os.FileInfo, rangeHeader string) {
	// 解析Range头
	ranges, err := parseRange(rangeHeader, fileInfo.Size())
	if err != nil {
		http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// 目前只支持单个Range
	if len(ranges) != 1 {
		http.Error(w, "Multiple ranges not supported", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	rangeInfo := ranges[0]

	// 设置响应头
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeInfo.start, rangeInfo.end, fileInfo.Size()))
	w.Header().Set("Content-Length", strconv.FormatInt(rangeInfo.length, 10))
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusPartialContent)

	// 定位到指定位置
	if _, err := file.Seek(rangeInfo.start, io.SeekStart); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 传输指定范围的数据
	if err := s.sendFileRange(w, file, rangeInfo.length); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// sendFileRange 传输文件范围
func (s *FileServer) sendFileRange(w http.ResponseWriter, file *os.File, length int64) error {
	// 获取缓冲区
	buffer := s.bufferPool.Get()
	defer s.bufferPool.Put(buffer)

	// 限制传输长度
	limitedReader := io.LimitReader(file, length)
	written, err := io.CopyBuffer(w, limitedReader, buffer)
	if err != nil {
		return err
	}

	s.metrics.RecordBytesTransferred(written)
	return nil
}

// Range 范围信息
type Range struct {
	start  int64
	end    int64
	length int64
}

// parseRange 解析Range头
func parseRange(rangeHeader string, fileSize int64) ([]Range, error) {
	// 简单的Range解析实现
	// 格式: "bytes=start-end"
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return nil, fmt.Errorf("invalid range format")
	}

	rangeStr := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range format")
	}

	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}

	var end int64
	if parts[1] == "" {
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	// 验证范围
	if start < 0 || end >= fileSize || start > end {
		return nil, fmt.Errorf("invalid range")
	}

	return []Range{
		{
			start:  start,
			end:    end,
			length: end - start + 1,
		},
	}, nil
}

// getContentType 根据文件扩展名获取Content-Type
func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

// GetMetrics 获取性能指标
func (s *FileServer) GetMetrics() *Metrics {
	return s.metrics
}

// Start 启动文件服务器
func (s *FileServer) Start(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: s,
	}

	fmt.Printf("File server starting on %s, serving directory: %s\n", addr, s.rootDir)
	return server.ListenAndServe()
}
