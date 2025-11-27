package interfaces

import "context"

// Repository 通用仓储接口（框架抽象）
// 用户需要根据业务需求定义具体的仓储接口
type Repository[T any] interface {
	// Create 创建实体
	Create(ctx context.Context, entity *T) error

	// FindByID 根据ID查找实体
	FindByID(ctx context.Context, id string) (*T, error)

	// Update 更新实体
	Update(ctx context.Context, entity *T) error

	// Delete 删除实体
	Delete(ctx context.Context, id string) error

	// List 列出实体（支持分页）
	List(ctx context.Context, limit, offset int) ([]*T, error)
}

// RepositoryWithQuery 支持查询的仓储接口（扩展接口）
type RepositoryWithQuery[T any, Q any] interface {
	Repository[T]

	// FindByQuery 根据查询条件查找实体
	FindByQuery(ctx context.Context, query Q) ([]*T, error)

	// Count 统计实体数量
	Count(ctx context.Context) (int, error)

	// CountByQuery 根据查询条件统计实体数量
	CountByQuery(ctx context.Context, query Q) (int, error)
}
