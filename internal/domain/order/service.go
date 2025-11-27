package order

import "context"

// DomainService 订单领域服务接口
type DomainService interface {
	// ValidateOrderItems 验证订单项
	ValidateOrderItems(ctx context.Context, items []OrderItem) error

	// CheckStock 检查库存
	CheckStock(ctx context.Context, productID string, quantity int) (bool, error)

	// CalculateShippingFee 计算运费
	CalculateShippingFee(ctx context.Context, order *Order) (float64, error)

	// ValidateAddress 验证地址
	ValidateAddress(ctx context.Context, address string) (bool, error)
}
