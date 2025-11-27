package patterns

import "context"

// Query 查询接口（框架抽象）
//
// 设计原理：
// 1. 查询模式（Query Pattern）是 CQRS 的核心
// 2. 查询用于表示读操作（Read、List、Search），不改变系统状态
// 3. 查询应该是幂等的，多次执行相同查询应该返回相同结果
// 4. 查询应该包含查询所需的所有条件
//
// 架构位置：
// - 定义：Application Layer (internal/application/patterns/)
// - 使用：Application Layer 和 Interfaces Layer
//
// CQRS 原理：
// - Command（命令）：写操作，改变系统状态
// - Query（查询）：读操作，不改变系统状态
// - 分离读写操作，可以独立优化和扩展
//
// 使用场景：
// 1. 查询单个实体
// 2. 查询实体列表（支持分页、排序、过滤）
// 3. 搜索操作
// 4. 统计查询
//
// 示例：
//   // 定义查询
//   type GetUserQuery struct {
//       ID string
//   }
//
//   func (q GetUserQuery) Execute(ctx context.Context) (interface{}, error) {
//       // 查询本身不包含执行逻辑，只是数据载体
//       return nil, nil
//   }
//
//   // 定义查询处理器
//   type GetUserQueryHandler struct {
//       userRepo domain.UserRepository
//   }
//
//   func (h *GetUserQueryHandler) Handle(ctx context.Context, query GetUserQuery) (*UserDTO, error) {
//       // 1. 验证查询参数
//       if query.ID == "" {
//           return nil, ErrInvalidQuery
//       }
//
//       // 2. 查询领域实体
//       user, err := h.userRepo.FindByID(ctx, query.ID)
//       if err != nil {
//           return nil, err
//       }
//
//       // 3. 转换为 DTO
//       return toUserDTO(user), nil
//   }
//
// 注意事项：
// - 查询应该是不可变的（immutable）
// - 查询应该是幂等的
// - 查询不应该改变系统状态
// - 查询处理器应该处理缓存（可选）
//
// 与 Command 的区别：
// - Query：读操作，不改变状态，返回数据
// - Command：写操作，改变状态，不返回数据
//
// 用户需要根据业务需求定义具体的查询，例如：
//   type ListUsersQuery struct {
//       Page  int
//       Size  int
//       Email string  // 可选过滤条件
//   }
//
//   func (q ListUsersQuery) Execute(ctx context.Context) (interface{}, error) {
//       return nil, nil
//   }
type Query interface {
	// Execute 执行查询
	//
	// 注意：这个方法通常不包含实际的执行逻辑，只是满足接口要求
	// 实际的执行逻辑应该在 QueryHandler 中实现
	//
	// 参数：
	//   - ctx: 上下文，用于传递请求信息、超时控制等
	//
	// 返回：
	//   - interface{}: 查询结果（通常不使用，实际结果在 Handler 中返回）
	//   - error: 执行失败时返回错误
	//
	// 实现建议：
	//   - 可以在这里添加查询级别的验证
	//   - 可以返回 nil, nil，实际逻辑在 Handler 中
	Execute(ctx context.Context) (interface{}, error)
}

// QueryHandler 查询处理器接口（框架抽象）
//
// 设计原理：
// 1. 查询处理器负责执行查询的实际逻辑
// 2. 使用泛型 T 表示查询类型，R 表示返回结果类型，确保类型安全
// 3. 查询处理器应该处理缓存、错误处理和结果转换
//
// 职责：
// 1. 验证查询参数
// 2. 从仓储或缓存获取数据
// 3. 转换领域对象为 DTO
// 4. 处理错误和异常
//
// 示例：
//   type GetUserQueryHandler struct {
//       userRepo domain.UserRepository
//       cache    cache.Cache
//   }
//
//   func NewGetUserQueryHandler(
//       userRepo domain.UserRepository,
//       cache cache.Cache,
//   ) *GetUserQueryHandler {
//       return &GetUserQueryHandler{
//           userRepo: userRepo,
//           cache:    cache,
//       }
//   }
//
//   func (h *GetUserQueryHandler) Handle(ctx context.Context, query GetUserQuery) (*UserDTO, error) {
//       // 1. 检查缓存
//       if cached, err := h.cache.Get(ctx, "user:"+query.ID); err == nil {
//           return cached.(*UserDTO), nil
//       }
//
//       // 2. 查询数据库
//       user, err := h.userRepo.FindByID(ctx, query.ID)
//       if err != nil {
//           return nil, err
//       }
//
//       // 3. 转换为 DTO
//       dto := toUserDTO(user)
//
//       // 4. 更新缓存
//       h.cache.Set(ctx, "user:"+query.ID, dto, 5*time.Minute)
//
//       return dto, nil
//   }
//
// 注意事项：
// - 查询处理器应该是无状态的
// - 应该通过依赖注入获取依赖
// - 应该考虑缓存策略
// - 应该处理查询性能优化
type QueryHandler[T Query, R any] interface {
	// Handle 处理查询
	//
	// 参数：
	//   - ctx: 上下文
	//   - query: 要处理的查询，类型为 T
	//
	// 返回：
	//   - R: 查询结果，类型为 R
	//   - error: 处理失败时返回错误
	//
	// 处理流程：
	// 1. 验证查询参数
	// 2. 检查缓存（可选）
	// 3. 查询数据（通过仓储）
	// 4. 转换领域对象为 DTO
	// 5. 更新缓存（可选）
	// 6. 返回结果
	//
	// 错误处理：
	// - 应该返回明确的错误信息
	// - 应该区分业务错误和系统错误
	// - 应该记录错误日志
	Handle(ctx context.Context, query T) (R, error)
}

// QueryResult 查询结果（分页结果）
//
// 设计原理：
// 1. 提供标准化的分页查询结果格式
// 2. 使用泛型 T 表示结果数据类型
// 3. 包含分页信息和数据列表
//
// 使用场景：
// 1. 列表查询（支持分页）
// 2. 搜索查询（支持分页）
// 3. 统计查询（返回总数）
//
// 示例：
//   type ListUsersQueryHandler struct {
//       userRepo domain.UserRepository
//   }
//
//   func (h *ListUsersQueryHandler) Handle(ctx context.Context, query ListUsersQuery) (*QueryResult[*UserDTO], error) {
//       // 1. 查询数据
//       users, err := h.userRepo.List(ctx, query.Page, query.Size)
//       if err != nil {
//           return nil, err
//       }
//
//       // 2. 查询总数
//       total, err := h.userRepo.Count(ctx)
//       if err != nil {
//           return nil, err
//       }
//
//       // 3. 转换为 DTO
//       dtos := make([]*UserDTO, len(users))
//       for i, user := range users {
//           dtos[i] = toUserDTO(user)
//       }
//
//       // 4. 返回结果
//       return &QueryResult[*UserDTO]{
//           Data:  dtos,
//           Total: total,
//           Page:  query.Page,
//           Size:  query.Size,
//       }, nil
//   }
//
// 注意事项：
// - Data 应该是切片，即使只有一条数据
// - Total 表示总记录数，不是当前页的数量
// - Page 和 Size 应该与查询参数一致
type QueryResult[T any] struct {
	// Data 查询结果数据列表
	// 注意：即使只有一条数据，也应该是切片
	Data []T

	// Total 总记录数（用于分页）
	// 注意：这是所有符合条件的记录数，不是当前页的数量
	Total int

	// Page 当前页码（从 1 开始）
	Page int

	// Size 每页数量
	Size int
}
