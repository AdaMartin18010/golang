package order

import (
	"time"

	domain "github.com/yourusername/golang/internal/domain/order"
)

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserID  string                `json:"user_id" validate:"required"`
	Items   []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
	Address string                `json:"address" validate:"required"`
}

// CreateOrderItemRequest 创建订单项请求
type CreateOrderItemRequest struct {
	ProductID string  `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,min=0"`
	Name      string  `json:"name" validate:"required"`
}

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
	Address *string `json:"address,omitempty" validate:"omitempty"`
}

// OrderDTO 订单数据传输对象
type OrderDTO struct {
	ID          string           `json:"id"`
	UserID      string           `json:"user_id"`
	Items       []OrderItemDTO   `json:"items"`
	TotalAmount float64          `json:"total_amount"`
	Status      string           `json:"status"`
	Address     string           `json:"address"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	PaidAt      *time.Time       `json:"paid_at,omitempty"`
	ShippedAt   *time.Time       `json:"shipped_at,omitempty"`
	DeliveredAt *time.Time       `json:"delivered_at,omitempty"`
}

// OrderItemDTO 订单项数据传输对象
type OrderItemDTO struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Name      string  `json:"name"`
}

// ToDTO 转换为 DTO
func ToDTO(o *domain.Order) *OrderDTO {
	if o == nil {
		return nil
	}

	items := make([]OrderItemDTO, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Name:      item.Name,
		}
	}

	return &OrderDTO{
		ID:          o.ID,
		UserID:      o.UserID,
		Items:       items,
		TotalAmount: o.TotalAmount,
		Status:      string(o.Status),
		Address:     o.Address,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
		PaidAt:      o.PaidAt,
		ShippedAt:   o.ShippedAt,
		DeliveredAt: o.DeliveredAt,
	}
}

// ToDomainItems 转换为领域订单项
func ToDomainItems(items []CreateOrderItemRequest) []domain.OrderItem {
	domainItems := make([]domain.OrderItem, len(items))
	for i, item := range items {
		domainItems[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Name:      item.Name,
		}
	}
	return domainItems
}
