package product

import "context"

// DomainService 产品领域服务接口
type DomainService interface {
	// ValidateSKU 验证SKU是否唯一
	ValidateSKU(ctx context.Context, sku string) (bool, error)

	// CheckAvailability 检查产品可用性
	CheckAvailability(ctx context.Context, productID string, quantity int) (bool, error)

	// CalculateDiscount 计算折扣
	CalculateDiscount(ctx context.Context, productID string, quantity int) (float64, error)

	// ValidateCategory 验证分类
	ValidateCategory(ctx context.Context, categoryID string) (bool, error)
}
