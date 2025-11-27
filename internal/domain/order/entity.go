package order

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待支付
	OrderStatusPaid      OrderStatus = "paid"      // 已支付
	OrderStatusShipped   OrderStatus = "shipped"   // 已发货
	OrderStatusDelivered OrderStatus = "delivered" // 已送达
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
	OrderStatusRefunded  OrderStatus = "refunded"  // 已退款
)

// OrderItem 订单项
type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
	Name      string
}

// Order 订单实体
type Order struct {
	ID          string
	UserID      string
	Items       []OrderItem
	TotalAmount float64
	Status      OrderStatus
	Address     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PaidAt      *time.Time
	ShippedAt   *time.Time
	DeliveredAt *time.Time
}

// NewOrder 创建新订单
func NewOrder(userID string, items []OrderItem, address string) *Order {
	now := time.Now()
	totalAmount := calculateTotalAmount(items)

	return &Order{
		ID:          uuid.New().String(),
		UserID:      userID,
		Items:       items,
		TotalAmount: totalAmount,
		Status:      OrderStatusPending,
		Address:     address,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// calculateTotalAmount 计算订单总金额
func calculateTotalAmount(items []OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// Pay 支付订单
func (o *Order) Pay() error {
	if o.Status != OrderStatusPending {
		return ErrInvalidOrderStatus
	}

	now := time.Now()
	o.Status = OrderStatusPaid
	o.PaidAt = &now
	o.UpdatedAt = now
	return nil
}

// Ship 发货
func (o *Order) Ship() error {
	if o.Status != OrderStatusPaid {
		return ErrInvalidOrderStatus
	}

	now := time.Now()
	o.Status = OrderStatusShipped
	o.ShippedAt = &now
	o.UpdatedAt = now
	return nil
}

// Deliver 送达
func (o *Order) Deliver() error {
	if o.Status != OrderStatusShipped {
		return ErrInvalidOrderStatus
	}

	now := time.Now()
	o.Status = OrderStatusDelivered
	o.DeliveredAt = &now
	o.UpdatedAt = now
	return nil
}

// Cancel 取消订单
func (o *Order) Cancel() error {
	if o.Status == OrderStatusDelivered || o.Status == OrderStatusShipped {
		return ErrCannotCancelOrder
	}

	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

// Refund 退款
func (o *Order) Refund() error {
	if o.Status != OrderStatusPaid && o.Status != OrderStatusShipped {
		return ErrCannotRefundOrder
	}

	o.Status = OrderStatusRefunded
	o.UpdatedAt = time.Now()
	return nil
}

// AddItem 添加订单项
func (o *Order) AddItem(item OrderItem) error {
	if o.Status != OrderStatusPending {
		return ErrCannotModifyOrder
	}

	o.Items = append(o.Items, item)
	o.TotalAmount = calculateTotalAmount(o.Items)
	o.UpdatedAt = time.Now()
	return nil
}

// RemoveItem 移除订单项
func (o *Order) RemoveItem(productID string) error {
	if o.Status != OrderStatusPending {
		return ErrCannotModifyOrder
	}

	for i, item := range o.Items {
		if item.ProductID == productID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			o.TotalAmount = calculateTotalAmount(o.Items)
			o.UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrItemNotFound
}

// IsValid 验证订单是否有效
func (o *Order) IsValid() bool {
	if o.ID == "" || o.UserID == "" || o.Address == "" {
		return false
	}

	if len(o.Items) == 0 {
		return false
	}

	if o.TotalAmount <= 0 {
		return false
	}

	return true
}

// CanBeCancelled 检查订单是否可以取消
func (o *Order) CanBeCancelled() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusPaid
}

// CanBeRefunded 检查订单是否可以退款
func (o *Order) CanBeRefunded() bool {
	return o.Status == OrderStatusPaid || o.Status == OrderStatusShipped
}
