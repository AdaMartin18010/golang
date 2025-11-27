package product

import (
	"time"

	"github.com/google/uuid"
)

// ProductStatus 产品状态
type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"   // 上架
	ProductStatusInactive ProductStatus = "inactive" // 下架
	ProductStatusDeleted  ProductStatus = "deleted"  // 已删除
)

// Product 产品实体
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int
	Status      ProductStatus
	CategoryID  string
	SKU         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProduct 创建新产品
func NewProduct(name, description string, price float64, stock int, categoryID, sku string) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		Status:      ProductStatusActive,
		CategoryID:  categoryID,
		SKU:         sku,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdatePrice 更新价格
func (p *Product) UpdatePrice(price float64) error {
	if price <= 0 {
		return ErrInvalidPrice
	}

	p.Price = price
	p.UpdatedAt = time.Now()
	return nil
}

// UpdateStock 更新库存
func (p *Product) UpdateStock(stock int) error {
	if stock < 0 {
		return ErrInvalidStock
	}

	p.Stock = stock
	p.UpdatedAt = time.Now()
	return nil
}

// IncreaseStock 增加库存
func (p *Product) IncreaseStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	p.Stock += quantity
	p.UpdatedAt = time.Now()
	return nil
}

// DecreaseStock 减少库存
func (p *Product) DecreaseStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if p.Stock < quantity {
		return ErrInsufficientStock
	}

	p.Stock -= quantity
	p.UpdatedAt = time.Now()
	return nil
}

// Activate 上架产品
func (p *Product) Activate() {
	p.Status = ProductStatusActive
	p.UpdatedAt = time.Now()
}

// Deactivate 下架产品
func (p *Product) Deactivate() {
	p.Status = ProductStatusInactive
	p.UpdatedAt = time.Now()
}

// Delete 删除产品（软删除）
func (p *Product) Delete() {
	p.Status = ProductStatusDeleted
	p.UpdatedAt = time.Now()
}

// IsAvailable 检查产品是否可用
func (p *Product) IsAvailable() bool {
	return p.Status == ProductStatusActive && p.Stock > 0
}

// IsValid 验证产品是否有效
func (p *Product) IsValid() bool {
	if p.ID == "" || p.Name == "" || p.SKU == "" {
		return false
	}

	if p.Price <= 0 {
		return false
	}

	return true
}

// HasStock 检查是否有库存
func (p *Product) HasStock(quantity int) bool {
	return p.Stock >= quantity
}
