// Package chi 提供基于 Chi Router 的 HTTP 接口实现
//
// 设计原理：
// 1. 这是 Interfaces Layer 的 HTTP 接口实现
// 2. 使用 Chi Router 进行路由管理
// 3. 调用 Application Layer 服务，不直接访问 Domain Layer
// 4. 负责协议适配、请求处理和响应格式化
//
// 架构位置：
// - 位置：Interfaces Layer (internal/interfaces/http/chi/)
// - 职责：HTTP 协议适配、路由管理、中间件配置
// - 依赖：Application Layer（调用应用服务）
//
// 中间件顺序：
// 1. RequestID - 生成请求ID
// 2. RealIP - 获取真实IP
// 3. Tracing - OpenTelemetry 追踪
// 4. Logging - 请求日志
// 5. Recovery - Panic 恢复
// 6. Timeout - 请求超时
// 7. CORS - 跨域支持
//
// 路由结构：
// - /health - 健康检查
// - /api/v1/users - 用户相关 API
// - /api/v1/workflows - 工作流相关 API
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
//
// 设计原理：
// 1. 封装 Chi Router，提供统一的接口
// 2. 管理所有 HTTP 路由和中间件
// 3. 通过依赖注入获取应用服务
//
// 职责：
// - 路由注册和管理
// - 中间件配置
// - 请求处理分发
type Router struct {
	// router Chi Router 实例
	router *chi.Mux
}

// NewRouter 创建 HTTP 路由器
//
// 设计原理：
// 1. 通过依赖注入获取应用服务
// 2. 配置中间件链
// 3. 注册所有路由
//
// 参数：
//   - userService: 用户应用服务（来自 Application Layer）
//   - temporalClient: Temporal 客户端处理器（可选）
//
// 返回：
//   - *Router: 创建的路由器实例
//
// 中间件说明：
// - RequestID: 为每个请求生成唯一ID，用于追踪
// - RealIP: 获取客户端真实IP（考虑代理）
// - Tracing: OpenTelemetry 分布式追踪
// - Logging: 记录请求日志
// - Recovery: 捕获 Panic，防止程序崩溃
// - Timeout: 设置请求超时时间（60秒）
// - CORS: 跨域资源共享支持
//
// 路由说明：
// - /health: 健康检查端点
// - /api/v1/users: 用户相关 REST API
// - /api/v1/workflows: 工作流相关 API（如果 Temporal 客户端可用）
//
// 使用示例：
//   userService := appuser.NewService(userRepo)
//   router := chiRouter.NewRouter(userService, temporalHandler)
//   httpServer := &http.Server{
//       Handler: router.Handler(),
//   }
func NewRouter(userService appuser.Service, temporalClient *temporalhandler.Handler) *Router {
	r := chi.NewRouter()

	// 中间件配置（按顺序执行）
	// 注意：中间件的执行顺序很重要，应该按照依赖关系排序
	r.Use(middleware.RequestID)                    // 1. 生成请求ID（最先执行，其他中间件可以使用）
	r.Use(middleware.RealIP)                       // 2. 获取真实IP
	r.Use(TracingMiddleware)                       // 3. OpenTelemetry 追踪（需要 RequestID）
	r.Use(LoggingMiddleware)                       // 4. 请求日志（需要 RequestID 和 Tracing）
	r.Use(RecovererMiddleware)                     // 5. Panic 恢复（保护所有后续处理）
	r.Use(TimeoutMiddleware(60 * time.Second))     // 6. 请求超时（60秒）
	r.Use(CORSMiddleware)                          // 7. CORS 支持（最后执行，处理响应头）

	// 健康检查端点
	// 用途：用于负载均衡器、监控系统等检查服务健康状态
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API 路由组
	// 路径前缀：/api/v1
	// 用途：版本化 API，便于后续版本升级
	r.Route("/api/v1", func(r chi.Router) {
		// 用户相关路由
		// 路径：/api/v1/users
		userHandler := handlers.NewUserHandler(userService)
		r.Mount("/users", userRoutes(userHandler))

		// 工作流相关路由（可选）
		// 路径：/api/v1/workflows
		// 只有当 Temporal 客户端可用时才注册
		if temporalClient != nil {
			workflowHandler := handlers.NewWorkflowHandler(temporalClient)
			r.Mount("/workflows", workflowRoutes(workflowHandler))
		}
	})

	return &Router{router: r}
}

// Handler 返回 HTTP 处理器
//
// 设计原理：
// 1. 返回 Chi Router 作为 HTTP Handler
// 2. 可以用于 http.Server 的 Handler 字段
//
// 返回：
//   - http.Handler: HTTP 处理器接口
//
// 使用示例：
//   router := chiRouter.NewRouter(userService, temporalHandler)
//   httpServer := &http.Server{
//       Handler: router.Handler(),
//   }
func (r *Router) Handler() http.Handler {
	return r.router
}

// userRoutes 用户路由
//
// 设计原理：
// 1. 定义用户相关的 REST API 路由
// 2. 遵循 RESTful 设计原则
// 3. 将 HTTP 请求路由到对应的处理器方法
//
// 路由定义：
// - POST   /users      - 创建用户
// - GET    /users      - 列出用户（支持分页、过滤）
// - GET    /users/{id} - 获取单个用户
// - PUT    /users/{id} - 更新用户
// - DELETE /users/{id} - 删除用户
//
// RESTful 设计：
// - 使用 HTTP 方法表示操作类型
// - 使用 URL 路径表示资源
// - 使用 HTTP 状态码表示结果
//
// 参数：
//   - userHandler: 用户处理器（包含所有用户相关的处理方法）
//
// 返回：
//   - http.Handler: 路由处理器
func userRoutes(userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/", userHandler.CreateUser)      // POST /api/v1/users
	r.Get("/", userHandler.ListUsers)        // GET /api/v1/users
	r.Get("/{id}", userHandler.GetUser)      // GET /api/v1/users/{id}
	r.Put("/{id}", userHandler.UpdateUser)   // PUT /api/v1/users/{id}
	r.Delete("/{id}", userHandler.DeleteUser) // DELETE /api/v1/users/{id}
	return r
}

// workflowRoutes 工作流路由
//
// 设计原理：
// 1. 定义工作流相关的 API 路由
// 2. 用于启动和查询 Temporal 工作流
// 3. 提供工作流管理的 HTTP 接口
//
// 路由定义：
// - POST /workflows/user - 启动用户工作流
// - GET  /workflows/user/{workflow_id}/result - 获取工作流结果
//
// 参数：
//   - workflowHandler: 工作流处理器
//
// 返回：
//   - http.Handler: 路由处理器
func workflowRoutes(workflowHandler *handlers.WorkflowHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/user", workflowHandler.StartUserWorkflow)                    // POST /api/v1/workflows/user
	r.Get("/user/{workflow_id}/result", workflowHandler.GetWorkflowResult) // GET /api/v1/workflows/user/{workflow_id}/result
	return r
}
