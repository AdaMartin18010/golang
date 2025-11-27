package order

import "errors"

var (
	// ErrOrderNotFound 订单未找到
	ErrOrderNotFound = errors.New("order not found")

	// ErrInvalidOrderStatus 无效的订单状态
	ErrInvalidOrderStatus = errors.New("invalid order status")

	// ErrCannotCancelOrder 无法取消订单
	ErrCannotCancelOrder = errors.New("cannot cancel order in current status")

	// ErrCannotRefundOrder 无法退款
	ErrCannotRefundOrder = errors.New("cannot refund order in current status")

	// ErrCannotModifyOrder 无法修改订单
	ErrCannotModifyOrder = errors.New("cannot modify order in current status")

	// ErrItemNotFound 订单项未找到
	ErrItemNotFound = errors.New("order item not found")

	// ErrInvalidOrderItem 无效的订单项
	ErrInvalidOrderItem = errors.New("invalid order item")

	// ErrInsufficientStock 库存不足
	ErrInsufficientStock = errors.New("insufficient stock")
)
