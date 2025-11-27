//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/internal/config"
	"github.com/yourusername/golang/internal/infrastructure/database/postgres"
	chirouter "github.com/yourusername/golang/internal/interfaces/http/chi"
)

// InitializeApp 初始化应用（Wire 自动生成）
//
// 设计原理：
// 1. Wire 是 Google 开源的 Go 语言依赖注入工具
// 2. Wire 在编译时生成依赖注入代码，而不是运行时反射
// 3. Wire 提供类型安全的依赖注入，编译时检查依赖关系
// 4. Wire 零反射，性能优秀
//
// 工作原理：
// 1. 定义 Provider 函数：每个 Provider 函数创建一个依赖
// 2. 使用 wire.Build 声明依赖关系
// 3. 运行 `go generate` 或 `wire` 命令生成代码
// 4. 生成的代码在 wire_gen.go 文件中
//
// 架构位置：
// - Wire 配置：scripts/wire/wire.go
// - 生成代码：scripts/wire/wire_gen.go
// - 使用位置：cmd/*/main.go
//
// 依赖注入流程：
// 1. 配置对象（Config）作为入口参数
// 2. Infrastructure Layer：创建数据库连接、仓储等
// 3. Application Layer：创建应用服务，注入仓储
// 4. Interfaces Layer：创建路由、处理器，注入应用服务
// 5. App：组装所有组件
//
// 依赖关系图：
//   Config
//     ↓
//   Infrastructure (数据库、消息队列等)
//     ↓
//   Application (应用服务)
//     ↓
//   Interfaces (HTTP 路由、gRPC 服务等)
//     ↓
//   App
//
// 使用步骤：
// 1. 定义 Provider 函数（在各个包中）
// 2. 在 wire.go 中使用 wire.Build 声明依赖
// 3. 运行 `go generate ./scripts/wire` 生成代码
// 4. 在 main.go 中调用 InitializeApp
//
// 示例：
//   // 1. 定义 Provider 函数
//   func NewUserRepository(client *ent.Client) domain.UserRepository {
//       return entrepo.NewUserRepository(client)
//   }
//
//   func NewUserService(repo domain.UserRepository) *appuser.Service {
//       return appuser.NewService(repo)
//   }
//
//   // 2. 在 wire.go 中声明依赖
//   func InitializeApp(cfg *config.Config) (*App, error) {
//       wire.Build(
//           NewEntClient,
//           NewUserRepository,
//           NewUserService,
//           NewRouter,
//           NewApp,
//       )
//       return &App{}, nil
//   }
//
//   // 3. 运行生成命令
//   // go generate ./scripts/wire
//
//   // 4. 在 main.go 中使用
//   app, err := wire.InitializeApp(cfg)
//   if err != nil {
//       log.Fatal(err)
//   }
//
// 优势：
// 1. 编译时检查：依赖关系在编译时检查，避免运行时错误
// 2. 类型安全：使用 Go 的类型系统，类型安全
// 3. 零反射：生成的代码不使用反射，性能优秀
// 4. 易于调试：生成的代码可以查看，易于调试
// 5. IDE 支持：IDE 可以理解依赖关系，提供代码补全
//
// 注意事项：
// - wire.go 文件需要 `//go:build wireinject` 构建标签
// - wire_gen.go 文件需要 `//go:build !wireinject` 构建标签
// - Provider 函数应该返回错误（如果可能失败）
// - Provider 函数应该遵循命名规范：NewXxx
// - 使用 wire.NewSet 组织 Provider 函数
//
// 最佳实践：
// 1. 按层次组织 Provider：Infrastructure、Application、Interfaces
// 2. 使用 wire.NewSet 创建 Provider 集合
// 3. Provider 函数应该简单，只负责创建依赖
// 4. 使用接口绑定（wire.Bind）绑定接口和实现
// 5. 错误处理应该在 Provider 函数中
func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		// Infrastructure Layer Providers
		// 基础设施层：创建数据库连接、仓储等
		postgres.NewUserRepository,

		// Application Layer Providers
		// 应用层：创建应用服务，注入仓储
		appuser.NewService,

		// Interfaces Layer Providers
		// 接口层：创建路由、处理器，注入应用服务
		chirouter.NewRouter,

		// App Provider
		// 应用组装：组装所有组件
		NewApp,
	)
	// 注意：这里的返回值不会被使用，Wire 会生成实际的实现
	return &App{}, nil
}

// App 应用结构
//
// 设计原理：
// 1. App 是应用的根对象，包含所有顶级组件
// 2. App 负责启动和关闭应用
// 3. App 是依赖注入的最终产物
//
// 职责：
// 1. 持有所有顶级组件（路由、服务等）
// 2. 提供启动和关闭方法
// 3. 管理应用生命周期
//
// 示例：
//   app, err := wire.InitializeApp(cfg)
//   if err != nil {
//       log.Fatal(err)
//   }
//
//   // 启动应用
//   if err := app.Start(); err != nil {
//       log.Fatal(err)
//   }
//
//   // 关闭应用
//   defer app.Close()
type App struct {
	// Router HTTP 路由
	// 注意：根据实际需求添加更多字段
	Router *chirouter.Router
}

// NewApp 创建应用
//
// 设计原理：
// 1. 这是 App 的 Provider 函数
// 2. Wire 会自动注入依赖（如 Router）
// 3. 可以在这里初始化应用状态
//
// 参数：
//   - router: HTTP 路由，由 Wire 自动注入
//
// 返回：
//   - *App: 创建的应用实例
//
// 示例：
//   func NewApp(router *chirouter.Router) *App {
//       return &App{
//           Router: router,
//       }
//   }
//
// 注意事项：
// - 可以添加更多参数，Wire 会自动注入
// - 可以在这里初始化应用状态
// - 可以在这里注册关闭钩子
func NewApp(router *chirouter.Router) *App {
	return &App{Router: router}
}
