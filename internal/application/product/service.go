package product

import (
	"context"
	"fmt"

	domain "github.com/yourusername/golang/internal/domain/product"
)

// Service 产品应用服务接口
type Service interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductDTO, error)
	GetProduct(ctx context.Context, id string) (*ProductDTO, error)
	GetProductBySKU(ctx context.Context, sku string) (*ProductDTO, error)
	UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*ProductDTO, error)
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, limit, offset int) ([]*ProductDTO, error)
	SearchProducts(ctx context.Context, keyword string, limit, offset int) ([]*ProductDTO, error)
	GetProductsByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*ProductDTO, error)
	ActivateProduct(ctx context.Context, id string) (*ProductDTO, error)
	DeactivateProduct(ctx context.Context, id string) (*ProductDTO, error)
	UpdateStock(ctx context.Context, id string, stock int) (*ProductDTO, error)
	IncreaseStock(ctx context.Context, id string, quantity int) (*ProductDTO, error)
	DecreaseStock(ctx context.Context, id string, quantity int) (*ProductDTO, error)
}

// service 产品应用服务实现
type service struct {
	productRepo domain.Repository
	domainSvc   domain.DomainService
}

// NewService 创建产品应用服务
func NewService(productRepo domain.Repository, domainSvc domain.DomainService) Service {
	return &service{
		productRepo: productRepo,
		domainSvc:   domainSvc,
	}
}

// CreateProduct 创建产品
func (s *service) CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductDTO, error) {
	// 验证 SKU 是否唯一
	unique, err := s.domainSvc.ValidateSKU(ctx, req.SKU)
	if err != nil {
		return nil, fmt.Errorf("failed to validate SKU: %w", err)
	}
	if !unique {
		return nil, domain.ErrDuplicateSKU
	}

	// 验证分类
	valid, err := s.domainSvc.ValidateCategory(ctx, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate category: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid category")
	}

	// 创建领域实体
	product := domain.NewProduct(req.Name, req.Description, req.Price, req.Stock, req.CategoryID, req.SKU)

	// 验证产品有效性
	if !product.IsValid() {
		return nil, fmt.Errorf("invalid product")
	}

	// 保存到仓储
	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return ToDTO(product), nil
}

// GetProduct 获取产品
func (s *service) GetProduct(ctx context.Context, id string) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return ToDTO(product), nil
}

// GetProductBySKU 根据 SKU 获取产品
func (s *service) GetProductBySKU(ctx context.Context, sku string) (*ProductDTO, error) {
	product, err := s.productRepo.FindBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}
	return ToDTO(product), nil
}

// UpdateProduct 更新产品
func (s *service) UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		product.Name = *req.Name
		product.UpdatedAt = product.UpdatedAt
	}
	if req.Description != nil {
		product.Description = *req.Description
		product.UpdatedAt = product.UpdatedAt
	}
	if req.Price != nil {
		if err := product.UpdatePrice(*req.Price); err != nil {
			return nil, fmt.Errorf("failed to update price: %w", err)
		}
	}
	if req.Stock != nil {
		if err := product.UpdateStock(*req.Stock); err != nil {
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}
	}
	if req.CategoryID != nil {
		valid, err := s.domainSvc.ValidateCategory(ctx, *req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate category: %w", err)
		}
		if !valid {
			return nil, fmt.Errorf("invalid category")
		}
		product.CategoryID = *req.CategoryID
	}
	if req.SKU != nil {
		if *req.SKU != product.SKU {
			unique, err := s.domainSvc.ValidateSKU(ctx, *req.SKU)
			if err != nil {
				return nil, fmt.Errorf("failed to validate SKU: %w", err)
			}
			if !unique {
				return nil, domain.ErrDuplicateSKU
			}
		}
		product.SKU = *req.SKU
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return ToDTO(product), nil
}

// DeleteProduct 删除产品
func (s *service) DeleteProduct(ctx context.Context, id string) error {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	product.Delete()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// ListProducts 列出产品
func (s *service) ListProducts(ctx context.Context, limit, offset int) ([]*ProductDTO, error) {
	products, err := s.productRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	dtos := make([]*ProductDTO, len(products))
	for i, p := range products {
		dtos[i] = ToDTO(p)
	}

	return dtos, nil
}

// SearchProducts 搜索产品
func (s *service) SearchProducts(ctx context.Context, keyword string, limit, offset int) ([]*ProductDTO, error) {
	products, err := s.productRepo.Search(ctx, keyword, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	dtos := make([]*ProductDTO, len(products))
	for i, p := range products {
		dtos[i] = ToDTO(p)
	}

	return dtos, nil
}

// GetProductsByCategory 根据分类获取产品
func (s *service) GetProductsByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*ProductDTO, error) {
	products, err := s.productRepo.FindByCategory(ctx, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	dtos := make([]*ProductDTO, len(products))
	for i, p := range products {
		dtos[i] = ToDTO(p)
	}

	return dtos, nil
}

// ActivateProduct 上架产品
func (s *service) ActivateProduct(ctx context.Context, id string) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	product.Activate()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to activate product: %w", err)
	}

	return ToDTO(product), nil
}

// DeactivateProduct 下架产品
func (s *service) DeactivateProduct(ctx context.Context, id string) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	product.Deactivate()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to deactivate product: %w", err)
	}

	return ToDTO(product), nil
}

// UpdateStock 更新库存
func (s *service) UpdateStock(ctx context.Context, id string, stock int) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if err := product.UpdateStock(stock); err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return ToDTO(product), nil
}

// IncreaseStock 增加库存
func (s *service) IncreaseStock(ctx context.Context, id string, quantity int) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if err := product.IncreaseStock(quantity); err != nil {
		return nil, fmt.Errorf("failed to increase stock: %w", err)
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return ToDTO(product), nil
}

// DecreaseStock 减少库存
func (s *service) DecreaseStock(ctx context.Context, id string, quantity int) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if err := product.DecreaseStock(quantity); err != nil {
		return nil, fmt.Errorf("failed to decrease stock: %w", err)
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return ToDTO(product), nil
}
