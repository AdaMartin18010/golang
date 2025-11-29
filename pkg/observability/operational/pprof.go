package operational

import (
	_ "net/http/pprof" // 导入 pprof
	"net/http"
)

// EnablePprof 启用性能分析端点
// 需要在 HTTP 服务器中注册
func EnablePprof(mux *http.ServeMux, prefix string) {
	if prefix == "" {
		prefix = "/debug/pprof"
	}

	// pprof 端点已经通过导入 _ "net/http/pprof" 自动注册到 http.DefaultServeMux
	// 这里只需要将请求转发到 DefaultServeMux
	mux.HandleFunc(prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	})
}

// PprofEndpoints 返回所有 pprof 端点列表
func PprofEndpoints() []string {
	return []string{
		"/debug/pprof/",
		"/debug/pprof/goroutine",
		"/debug/pprof/heap",
		"/debug/pprof/allocs",
		"/debug/pprof/block",
		"/debug/pprof/mutex",
		"/debug/pprof/profile",
		"/debug/pprof/trace",
	}
}
