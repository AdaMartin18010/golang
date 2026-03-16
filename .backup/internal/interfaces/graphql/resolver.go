// Package graphql provides GraphQL resolvers for the interfaces layer.
//
// GraphQL 解析器负责：
// 1. 解析 GraphQL 查询（Query）和变更（Mutation）
// 2. 将 GraphQL 请求转换为应用层调用
// 3. 将应用层结果转换为 GraphQL 响应
// 4. 处理字段解析和嵌套查询
//
// 设计原则：
// 1. 薄层设计：解析器只负责协议转换，不包含业务逻辑
// 2. 按需加载：只解析客户端请求的字段
// 3. 批量加载：使用 DataLoader 解决 N+1 查询问题
// 4. 错误处理：返回 GraphQL 错误格式
//
// GraphQL 核心概念：
// - Query: 查询数据（类似 REST GET）
// - Mutation: 修改数据（类似 REST POST/PUT/DELETE）
// - Subscription: 实时推送（WebSocket）
// - Resolver: 字段解析函数
// - DataLoader: 批量加载数据，避免 N+1 问题
//
// 工作流程：
// 1. 客户端发送 GraphQL 查询
// 2. GraphQL 引擎解析查询
// 3. 调用对应的 Resolver 方法
// 4. Resolver 调用应用层服务
// 5. 返回结果给 GraphQL 引擎
// 6. GraphQL 引擎组装响应
//
// 示例：
//
//	// 创建解析器
//	userService := appuser.NewService(userRepo)
//	resolver := graphql.NewResolver(userService)
//
//	// 注册到 GraphQL 服务器
//	schema := graphql.MustParseSchema(schemaString, resolver)
package graphql

import (
	"context"

	domainuser "github.com/yourusername/golang/internal/domain/user"
)

// UserService 定义用户服务接口。
//
// 功能说明：
// - 解耦 GraphQL 解析器与具体的应用层服务实现
// - 便于单元测试时使用 mock 实现
//
// 方法列表：
// - GetUser: 根据 ID 获取用户
// - CreateUser: 创建新用户
// - UpdateUserName: 更新用户名称
// - DeleteUser: 删除用户
// - ListUsers: 列出用户
type UserService interface {
	GetUser(ctx context.Context, id string) (*domainuser.User, error)
	CreateUser(ctx context.Context, email, name string) (*domainuser.User, error)
	UpdateUserName(ctx context.Context, id, name string) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*domainuser.User, error)
}

// Resolver 是 GraphQL 的根解析器。
//
// 功能说明：
// - 包含所有查询和变更解析器
// - 持有应用层服务接口（而非具体实现）
// - 作为 GraphQL Schema 的入口点
//
// 设计说明：
// - 依赖 UserService 接口处理业务逻辑
// - 通过依赖注入接收服务实例
// - 可以包含多个服务的引用
//
// 使用示例：
//
//	userService := appuser.NewService(userRepo)
//	resolver := graphql.NewResolver(userService)
//
// 注意事项：
// - Resolver 应该是无状态的
// - 可以在多个请求间共享
// - 可以使用 DataLoader 优化性能
type Resolver struct {
	userService UserService
	// 可以添加更多服务：
	// postService PostService
	// orderService OrderService
}

// NewResolver 创建并初始化 GraphQL 根解析器。
//
// 功能说明：
// - 接收应用层服务接口实例
// - 创建并返回配置好的解析器
//
// 参数：
// - userService: 应用层用户服务接口实例
//
// 返回：
// - *Resolver: 配置好的解析器实例
//
// 使用示例：
//
//	userService := appuser.NewService(userRepo)
//	resolver := graphql.NewResolver(userService)
//
// 注意事项：
// - 服务实例应通过依赖注入提供
// - 可以接受任何实现 UserService 接口的类型
// - 便于单元测试时使用 mock 实现
func NewResolver(userService UserService) *Resolver {
	return &Resolver{
		userService: userService,
	}
}

// Query 是 GraphQL 查询解析器。
//
// 功能说明：
// - 处理所有查询操作（Query）
// - 实现 GraphQL Schema 中定义的查询字段
// - 调用应用层服务获取数据
//
// 设计说明：
// - 每个查询字段对应一个方法
// - 方法签名必须匹配 GraphQL Schema
// - 支持参数、分页、过滤等
//
// 使用示例：
//
//	query {
//	  user(id: "123") {
//	    id
//	    name
//	    email
//	  }
//	}
type Query struct {
	resolver *Resolver
}

// Mutation 是 GraphQL 变更解析器。
//
// 功能说明：
// - 处理所有变更操作（Mutation）
// - 实现 GraphQL Schema 中定义的变更字段
// - 调用应用层服务修改数据
//
// 设计说明：
// - 每个变更字段对应一个方法
// - 方法签名必须匹配 GraphQL Schema
// - 支持输入参数和返回结果
//
// 使用示例：
//
//	mutation {
//	  createUser(input: {email: "user@example.com", name: "User"}) {
//	    id
//	    email
//	    name
//	  }
//	}
type Mutation struct {
	resolver *Resolver
}

// User 是 GraphQL User 类型。
//
// 功能说明：
// - 表示 GraphQL Schema 中的 User 类型
// - 包含用户的字段定义
// - 用于查询和变更的返回类型
//
// 字段说明：
// - ID: 用户唯一标识符
// - Email: 用户邮箱
// - Name: 用户名称
// - CreatedAt: 创建时间（字符串格式）
// - UpdatedAt: 更新时间（字符串格式）
//
// 注意事项：
// - 字段类型应与 GraphQL Schema 定义一致
// - 可以使用 GraphQL 标量类型（如 DateTime）
// - 可以添加更多字段（如 posts、comments 等关联数据）
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
	UpdatedAt string
}

// User 查询单个用户 - 实现在 resolver_impl.go

// Users 查询用户列表 - 实现在 resolver_impl.go

// CreateUser 创建用户 - 实现在 resolver_impl.go

// CreateUserInput 是创建用户的输入类型。
//
// 功能说明：
// - 定义创建用户所需的输入参数
// - 对应 GraphQL Schema 中的 input CreateUserInput
//
// 字段说明：
// - Email: 用户邮箱（必填）
// - Name: 用户名称（必填）
//
// 注意事项：
// - 字段类型应与 GraphQL Schema 定义一致
// - 可以添加验证标签（如 required、email 等）
// - 可以添加更多字段（如 password、phone 等）
type CreateUserInput struct {
	Email string
	Name  string
}
