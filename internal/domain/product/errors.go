package product

import "errors"

var (
	// ErrProductNotFound 产品未找到
	ErrProductNotFound = errors.New("product not found")

	// ErrInvalidPrice 无效的价格
	ErrInvalidPrice = errors.New("invalid price")

	// ErrInvalidStock 无效的库存
	ErrInvalidStock = errors.New("invalid stock")

	// ErrInvalidQuantity 无效的数量
	ErrInvalidQuantity = errors.New("invalid quantity")

	// ErrInsufficientStock 库存不足
	ErrInsufficientStock = errors.New("insufficient stock")

	// ErrProductNotAvailable 产品不可用
	ErrProductNotAvailable = errors.New("product not available")

	// ErrDuplicateSKU 重复的SKU
	ErrDuplicateSKU = errors.New("duplicate SKU")
)
