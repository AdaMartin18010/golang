package interfaces

import "context"

// Repository 通用仓储接口（框架抽象）
//
// 设计原理：
// 1. 仓储模式（Repository Pattern）是 DDD 中的核心模式，用于抽象数据访问层
// 2. 仓储接口定义在领域层，实现放在基础设施层，符合依赖倒置原则
// 3. 领域层不依赖具体的数据访问技术（如 Ent、GORM），只依赖接口
// 4. 这样可以轻松替换数据访问实现，提高可测试性和可维护性
//
// 架构位置：
// - 接口定义：Domain Layer (internal/domain/interfaces/)
// - 接口实现：Infrastructure Layer (internal/infrastructure/database/)
//
// 依赖关系：
// Domain Layer (接口定义) ← Infrastructure Layer (接口实现)
// Domain Layer (接口定义) ← Application Layer (使用接口)
//
// 使用场景：
// 1. 应用层通过仓储接口访问领域实体
// 2. 基础设施层实现仓储接口，封装数据访问细节
// 3. 测试时可以使用 Mock 仓储实现
//
// 示例：
//   // 领域层定义接口
//   type UserRepository interface {
//       Create(ctx context.Context, user *User) error
//       FindByID(ctx context.Context, id string) (*User, error)
//   }
//
//   // 基础设施层实现接口
//   type EntUserRepository struct {
//       client *ent.Client
//   }
//
//   func (r *EntUserRepository) Create(ctx context.Context, user *User) error {
//       // 使用 Ent 实现
//   }
//
//   // 应用层使用接口
//   type UserService struct {
//       repo UserRepository  // 依赖接口，不依赖具体实现
//   }
//
// 注意事项：
// - 仓储接口应该表达业务需求，而不是技术细节
// - 接口方法应该使用领域对象，而不是数据库模型
// - 复杂查询应该通过领域服务或应用服务处理
//
// 用户需要根据业务需求定义具体的仓储接口，例如：
//   type UserRepository interface {
//       Repository[*User]  // 继承通用接口
//       FindByEmail(ctx context.Context, email string) (*User, error)  // 业务特定方法
//   }
type Repository[T any] interface {
	// Create 创建实体
	//
	// 参数：
	//   - ctx: 上下文，用于传递请求信息、超时控制等
	//   - entity: 要创建的实体指针
	//
	// 返回：
	//   - error: 创建失败时返回错误
	//
	// 业务规则：
	//   - 实体应该已经通过业务规则验证
	//   - 如果实体已存在（通过唯一标识判断），应该返回错误
	//   - 创建成功后，实体的 ID 和时间戳应该被设置
	//
	// 实现要求：
	//   - 基础设施层实现应该处理数据库事务
	//   - 应该处理并发冲突（如唯一约束冲突）
	//   - 应该记录操作日志（可选）
	Create(ctx context.Context, entity *T) error

	// FindByID 根据ID查找实体
	//
	// 参数：
	//   - ctx: 上下文
	//   - id: 实体的唯一标识
	//
	// 返回：
	//   - *T: 找到的实体指针，如果不存在返回 nil
	//   - error: 查找失败时返回错误
	//
	// 业务规则：
	//   - ID 应该是有效的唯一标识
	//   - 如果实体不存在，应该返回 nil 和 nil error（不是错误）
	//   - 如果实体被软删除，应该根据业务规则决定是否返回
	//
	// 实现要求：
	//   - 应该使用索引优化查询性能
	//   - 应该处理查询超时
	//   - 应该处理数据库连接错误
	FindByID(ctx context.Context, id string) (*T, error)

	// Update 更新实体
	//
	// 参数：
	//   - ctx: 上下文
	//   - entity: 要更新的实体指针
	//
	// 返回：
	//   - error: 更新失败时返回错误
	//
	// 业务规则：
	//   - 实体应该已经通过业务规则验证
	//   - 实体必须存在（通过 ID 判断）
	//   - 更新成功后，实体的 UpdatedAt 应该被更新
	//
	// 实现要求：
	//   - 应该使用乐观锁或悲观锁处理并发更新
	//   - 应该只更新变更的字段（可选）
	//   - 应该处理版本冲突
	Update(ctx context.Context, entity *T) error

	// Delete 删除实体
	//
	// 参数：
	//   - ctx: 上下文
	//   - id: 要删除的实体的唯一标识
	//
	// 返回：
	//   - error: 删除失败时返回错误
	//
	// 业务规则：
	//   - 实体必须存在
	//   - 应该根据业务规则决定是硬删除还是软删除
	//   - 如果实体有关联数据，应该根据业务规则处理
	//
	// 实现要求：
	//   - 应该处理级联删除（如果需要）
	//   - 应该处理外键约束
	//   - 软删除时应该更新 DeletedAt 字段
	Delete(ctx context.Context, id string) error

	// List 列出实体（支持分页）
	//
	// 参数：
	//   - ctx: 上下文
	//   - limit: 每页数量，用于分页
	//   - offset: 偏移量，用于分页
	//
	// 返回：
	//   - []*T: 实体列表
	//   - error: 查询失败时返回错误
	//
	// 业务规则：
	//   - limit 应该大于 0，建议设置最大值限制
	//   - offset 应该大于等于 0
	//   - 如果没有任何实体，应该返回空切片，不是 nil
	//   - 应该根据业务规则过滤已删除的实体
	//
	// 实现要求：
	//   - 应该使用分页查询优化性能
	//   - 应该处理大数据量的情况
	//   - 应该考虑使用游标分页（可选）
	List(ctx context.Context, limit, offset int) ([]*T, error)
}

// RepositoryWithQuery 支持查询的仓储接口（扩展接口）
//
// 设计原理：
// 1. 这是 Repository 接口的扩展，提供更强大的查询能力
// 2. 使用泛型 Q 表示查询条件类型，提供类型安全
// 3. 查询条件应该在领域层定义，表达业务需求
//
// 使用场景：
// 1. 需要复杂查询条件的场景
// 2. 需要统计功能的场景
// 3. 需要动态查询的场景
//
// 示例：
//   // 定义查询条件
//   type UserQuery struct {
//       Email    string
//       Status   string
//       CreatedAfter time.Time
//   }
//
//   // 定义仓储接口
//   type UserRepository interface {
//       Repository[*User]
//       RepositoryWithQuery[*User, UserQuery]
//   }
//
//   // 使用
//   query := UserQuery{Email: "test@example.com"}
//   users, err := repo.FindByQuery(ctx, query)
//
// 注意事项：
// - 查询条件应该表达业务需求，而不是 SQL 查询
// - 应该避免在查询条件中暴露数据库细节
// - 复杂查询应该考虑性能影响
type RepositoryWithQuery[T any, Q any] interface {
	Repository[T] // 继承基础仓储接口

	// FindByQuery 根据查询条件查找实体
	//
	// 参数：
	//   - ctx: 上下文
	//   - query: 查询条件，类型为 Q
	//
	// 返回：
	//   - []*T: 匹配的实体列表
	//   - error: 查询失败时返回错误
	//
	// 业务规则：
	//   - 查询条件应该经过验证
	//   - 如果没有任何匹配的实体，应该返回空切片
	//   - 应该考虑查询性能，避免全表扫描
	//
	// 实现要求：
	//   - 应该使用索引优化查询
	//   - 应该处理查询超时
	//   - 应该考虑查询结果数量限制
	FindByQuery(ctx context.Context, query Q) ([]*T, error)

	// Count 统计实体数量
	//
	// 参数：
	//   - ctx: 上下文
	//
	// 返回：
	//   - int: 实体总数
	//   - error: 统计失败时返回错误
	//
	// 业务规则：
	//   - 应该根据业务规则过滤已删除的实体
	//   - 应该考虑性能，避免在大数据量时使用
	//
	// 实现要求：
	//   - 应该使用 COUNT 查询优化性能
	//   - 应该处理查询超时
	Count(ctx context.Context) (int, error)

	// CountByQuery 根据查询条件统计实体数量
	//
	// 参数：
	//   - ctx: 上下文
	//   - query: 查询条件
	//
	// 返回：
	//   - int: 匹配的实体数量
	//   - error: 统计失败时返回错误
	//
	// 业务规则：
	//   - 查询条件应该与 FindByQuery 保持一致
	//   - 应该用于分页时的总数统计
	//
	// 实现要求：
	//   - 应该使用 COUNT 查询优化性能
	//   - 应该复用 FindByQuery 的查询逻辑
	CountByQuery(ctx context.Context, query Q) (int, error)
}
