package order

import "context"

// Repository 订单仓储接口（领域层定义）
type Repository interface {
	// Create 创建订单
	Create(ctx context.Context, order *Order) error

	// FindByID 根据ID查找订单
	FindByID(ctx context.Context, id string) (*Order, error)

	// FindByUserID 根据用户ID查找订单列表
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*Order, error)

	// Update 更新订单
	Update(ctx context.Context, order *Order) error

	// Delete 删除订单
	Delete(ctx context.Context, id string) error

	// FindByStatus 根据状态查找订单
	FindByStatus(ctx context.Context, status OrderStatus, limit, offset int) ([]*Order, error)

	// CountByUserID 统计用户的订单数量
	CountByUserID(ctx context.Context, userID string) (int, error)
}
