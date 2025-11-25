package chi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/internal/interfaces/http/chi/handlers"
	temporalhandler "github.com/yourusername/golang/internal/interfaces/workflow/temporal"
)

// Router Chi 路由配置
type Router struct {
	router *chi.Mux
}

// NewRouter 创建 HTTP 路由器
func NewRouter(userService appuser.Service, temporalClient *temporalhandler.Handler) *Router {
	r := chi.NewRouter()

	// 中间件
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(TracingMiddleware) // OpenTelemetry 追踪
	r.Use(LoggingMiddleware)  // 日志
	r.Use(RecovererMiddleware) // 恢复
	r.Use(TimeoutMiddleware(60 * time.Second)) // 超时
	r.Use(CORSMiddleware) // CORS

	// 健康检查
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API 路由
	r.Route("/api/v1", func(r chi.Router) {
		userHandler := handlers.NewUserHandler(userService)
		r.Mount("/users", userRoutes(userHandler))

		// 工作流路由
		if temporalClient != nil {
			workflowHandler := handlers.NewWorkflowHandler(temporalClient)
			r.Mount("/workflows", workflowRoutes(workflowHandler))
		}
	})

	return &Router{router: r}
}

// Handler 返回 HTTP 处理器
func (r *Router) Handler() http.Handler {
	return r.router
}

// userRoutes 用户路由
func userRoutes(userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/", userHandler.CreateUser)
	r.Get("/", userHandler.ListUsers)
	r.Get("/{id}", userHandler.GetUser)
	r.Put("/{id}", userHandler.UpdateUser)
	r.Delete("/{id}", userHandler.DeleteUser)
	return r
}

// workflowRoutes 工作流路由
func workflowRoutes(workflowHandler *handlers.WorkflowHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/user", workflowHandler.StartUserWorkflow)
	r.Get("/user/{workflow_id}/result", workflowHandler.GetWorkflowResult)
	return r
}
