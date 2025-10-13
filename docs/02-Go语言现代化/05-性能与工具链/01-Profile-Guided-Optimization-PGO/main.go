package main

import (
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof" // 关键：导入 pprof 包以自动注册其 HTTP handler
	"os"
	"strconv"
)

// hotFunction 是一个 CPU 密集型函数，我们将通过 PGO 优化它。
// 它计算斐波那契数，这是一个典型的递归计算例子。
// 注意：这个实现效率很低，目的是为了消耗 CPU，从而在 profile 中成为"热点"。
func hotFunction(n int) int {
	if n <= 1 {
		return n
	}
	return hotFunction(n-1) + hotFunction(n-2)
}

// coldFunction 是一个不常被调用的函数，PGO 不会对其进行优化。
func coldFunction() string {
	return "This is a cold function."
}

// httpHandler 是我们的主要 HTTP 处理器。
// 它会根据查询参数决定调用热点函数还是冷点函数。
func httpHandler(w http.ResponseWriter, r *http.Request) {
	// 默认调用 hotFunction，使其成为 profile 中的热点路径
	call := r.URL.Query().Get("call")
	if call == "cold" {
		fmt.Fprint(w, coldFunction())
		return
	}

	nStr := r.URL.Query().Get("n")
	n, err := strconv.Atoi(nStr)
	if err != nil || n <= 0 {
		n = 25 // 默认值，以保证一定的计算量
	}

	result := hotFunction(n)
	fmt.Fprintf(w, "Result of hotFunction(%d) is %d\n", n, result)
}

func main() {
	// 注册 HTTP handler
	http.HandleFunc("/", httpHandler)

	// 在一个单独的 goroutine 中启动 pprof 服务器
	// 这样它就不会阻塞主服务逻辑
	go func() {
		slog.Info("Starting pprof server on :6060")
		// 通常 pprof 服务会绑定到 localhost 以避免对外暴露
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			slog.Error("Pprof server failed", "error", err)
		}
	}()

	port := "8080"
	slog.Info("Main server starting...", "port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("Main server failed", "error", err)
		os.Exit(1)
	}
}
