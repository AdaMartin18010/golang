package product

import (
	"time"

	domain "github.com/yourusername/golang/internal/domain/product"
)

// CreateProductRequest 创建产品请求
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=200"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"required,min=0"`
	CategoryID  string  `json:"category_id" validate:"required"`
	SKU         string  `json:"sku" validate:"required"`
}

// UpdateProductRequest 更新产品请求
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,min=0"`
	CategoryID  *string  `json:"category_id,omitempty"`
	SKU         *string  `json:"sku,omitempty"`
}

// ProductDTO 产品数据传输对象
type ProductDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status"`
	CategoryID  string    `json:"category_id"`
	SKU         string    `json:"sku"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToDTO 转换为 DTO
func ToDTO(p *domain.Product) *ProductDTO {
	if p == nil {
		return nil
	}

	return &ProductDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Status:      string(p.Status),
		CategoryID:  p.CategoryID,
		SKU:         p.SKU,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
