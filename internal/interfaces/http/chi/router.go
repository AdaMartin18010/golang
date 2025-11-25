package chi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/internal/interfaces/http/chi/handlers"
)

// Router Chi 路由配置
type Router struct {
	router *chi.Mux
}

// NewRouter 创建路由
func NewRouter(userService *appuser.Service) *Router {
	r := chi.NewRouter()

	// 中间件
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60))

	// 健康检查
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API 路由
	r.Route("/api/v1", func(r chi.Router) {
		// 用户路由
		userHandler := handlers.NewUserHandler(userService)
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.ListUsers)
			r.Get("/{id}", userHandler.GetUser)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})
	})

	return &Router{router: r}
}

// Handler 返回 HTTP 处理器
func (r *Router) Handler() http.Handler {
	return r.router
}
