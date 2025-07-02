package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// taskStore 是一个用于演示的内存任务存储。
var taskStore = map[string]string{
	"1": "Learn new ServeMux features",
	"2": "Build a sample API",
	"3": "Write documentation",
}

// getTaskHandler 处理获取单个任务的请求。
// 演示了如何从路径中提取参数。
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, ok := taskStore[id]
	if !ok {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Task %s: %s\n", id, task)
}

// getAllTasksHandler 处理获取所有任务的请求。
// 这是一个精确匹配的例子。
func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	for id, task := range taskStore {
		fmt.Fprintf(w, "- Task %s: %s\n", id, task)
	}
}

// createTaskHandler 处理创建新任务的请求。
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 在真实的app中，你会从 r.Body 中解析数据。
	newTaskID := fmt.Sprintf("%d", time.Now().UnixNano())
	taskStore[newTaskID] = "New task created at " + time.Now().String()
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Created new task with ID: %s\n", newTaskID)
}

// staticAssetHandler 演示了 `/{...}` 通配符的用法。
func staticAssetHandler(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	fmt.Fprintf(w, "Serving static asset at: %s\n", path)
}

// adminHandler 演示了按 Host 匹配。
func adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Admin Panel!\n")
}

func main() {
	mux := http.NewServeMux()

	// 1. 注册带路径参数的 GET 请求
	mux.HandleFunc("GET /tasks/{id}", getTaskHandler)
	slog.Info("Registered route: GET /tasks/{id}")

	// 2. 注册精确匹配的 GET 请求，它会优先于上面的通配符路由
	mux.HandleFunc("GET /tasks/all", getAllTasksHandler)
	slog.Info("Registered route: GET /tasks/all (exact match)")

	// 3. 注册 POST 请求
	mux.HandleFunc("POST /tasks", createTaskHandler)
	slog.Info("Registered route: POST /tasks")

	// 4. 注册一个服务静态资源的 `{...}` 通配符路由
	mux.HandleFunc("GET /static/{path...}", staticAssetHandler)
	slog.Info("Registered route: GET /static/{path...}")

	// 5. 注册一个只在特定 Hostname 下生效的路由
	//    要测试这个，你需要修改你的 hosts 文件 (e.g., `127.0.0.1 admin.example.com`)
	//    然后用 curl -H "Host: admin.example.com" http://localhost:8080/
	mux.HandleFunc("admin.example.com/", adminHandler)
	slog.Info("Registered host-specific route: admin.example.com/")

	// 设置一个简单的日志中间件来观察请求
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("Request received", "method", r.Method, "host", r.Host, "path", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	port := "8080"
	slog.Info("Server starting...", "port", port)

	// 启动服务器
	if err := http.ListenAndServe(":"+port, loggingMiddleware(mux)); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
