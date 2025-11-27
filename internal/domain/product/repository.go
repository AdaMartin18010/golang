package product

import "context"

// Repository 产品仓储接口（领域层定义）
type Repository interface {
	// Create 创建产品
	Create(ctx context.Context, product *Product) error

	// FindByID 根据ID查找产品
	FindByID(ctx context.Context, id string) (*Product, error)

	// FindBySKU 根据SKU查找产品
	FindBySKU(ctx context.Context, sku string) (*Product, error)

	// Update 更新产品
	Update(ctx context.Context, product *Product) error

	// Delete 删除产品
	Delete(ctx context.Context, id string) error

	// List 列出产品（支持分页）
	List(ctx context.Context, limit, offset int) ([]*Product, error)

	// FindByCategory 根据分类查找产品
	FindByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*Product, error)

	// FindByStatus 根据状态查找产品
	FindByStatus(ctx context.Context, status ProductStatus, limit, offset int) ([]*Product, error)

	// Search 搜索产品
	Search(ctx context.Context, keyword string, limit, offset int) ([]*Product, error)
}
