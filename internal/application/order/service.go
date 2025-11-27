package order

import (
	"context"
	"fmt"

	domain "github.com/yourusername/golang/internal/domain/order"
)

// Service 订单应用服务接口
type Service interface {
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderDTO, error)
	GetOrder(ctx context.Context, id string) (*OrderDTO, error)
	GetUserOrders(ctx context.Context, userID string, limit, offset int) ([]*OrderDTO, error)
	PayOrder(ctx context.Context, id string) (*OrderDTO, error)
	ShipOrder(ctx context.Context, id string) (*OrderDTO, error)
	DeliverOrder(ctx context.Context, id string) (*OrderDTO, error)
	CancelOrder(ctx context.Context, id string) (*OrderDTO, error)
	RefundOrder(ctx context.Context, id string) (*OrderDTO, error)
	UpdateOrder(ctx context.Context, id string, req UpdateOrderRequest) (*OrderDTO, error)
}

// service 订单应用服务实现
type service struct {
	orderRepo   domain.Repository
	domainSvc   domain.DomainService
}

// NewService 创建订单应用服务
func NewService(orderRepo domain.Repository, domainSvc domain.DomainService) Service {
	return &service{
		orderRepo: orderRepo,
		domainSvc: domainSvc,
	}
}

// CreateOrder 创建订单
func (s *service) CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderDTO, error) {
	// 验证订单项
	if err := s.domainSvc.ValidateOrderItems(ctx, ToDomainItems(req.Items)); err != nil {
		return nil, fmt.Errorf("failed to validate order items: %w", err)
	}

	// 验证地址
	valid, err := s.domainSvc.ValidateAddress(ctx, req.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to validate address: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid address")
	}

	// 创建领域实体
	order := domain.NewOrder(req.UserID, ToDomainItems(req.Items), req.Address)

	// 验证订单有效性
	if !order.IsValid() {
		return nil, fmt.Errorf("invalid order")
	}

	// 保存到仓储
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return ToDTO(order), nil
}

// GetOrder 获取订单
func (s *service) GetOrder(ctx context.Context, id string) (*OrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return ToDTO(order), nil
}

// GetUserOrders 获取用户订单列表
func (s *service) GetUserOrders(ctx context.Context, userID string, limit, offset int) ([]*OrderDTO, error) {
	orders, err := s.orderRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	dtos := make([]*OrderDTO, len(orders))
	for i, o := range orders {
		dtos[i] = ToDTO(o)
	}

	return dtos, nil
}

// PayOrder 支付订单
func (s *service) PayOrder(ctx context.Context, id string) (*OrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if err := order.Pay(); err != nil {
		return nil, fmt.Errorf("failed to pay order: %w", err)
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return ToDTO(order), nil
}

// ShipOrder 发货
func (s *service) ShipOrder(ctx context.Context, id string) (*OrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if err := order.Ship(); err != nil {
		return nil, fmt.Errorf("failed to ship order: %w", err)
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return ToDTO(order), nil
}

// DeliverOrder 送达
func (s *service) DeliverOrder(ctx context.Context, id string) (*OrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if err := order.Deliver(); err != nil {
		return nil, fmt.Errorf("failed to deliver order: %w", err)
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return ToDTO(order), nil
}

// CancelOrder 取消订单
func (s *service) CancelOrder(ctx context.Context, id string) (*OrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if err := order.Cancel(); err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return ToDTO(order), nil
}

// RefundOrder 退款
func (s *service) RefundOrder(ctx context.Context, id string) (*OrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if err := order.Refund(); err != nil {
		return nil, fmt.Errorf("failed to refund order: %w", err)
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return ToDTO(order), nil
}

// UpdateOrder 更新订单
func (s *service) UpdateOrder(ctx context.Context, id string, req UpdateOrderRequest) (*OrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// 只能更新待支付状态的订单
	if order.Status != domain.OrderStatusPending {
		return nil, domain.ErrCannotModifyOrder
	}

	// 更新地址
	if req.Address != nil {
		valid, err := s.domainSvc.ValidateAddress(ctx, *req.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to validate address: %w", err)
		}
		if !valid {
			return nil, fmt.Errorf("invalid address")
		}
		order.Address = *req.Address
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return ToDTO(order), nil
}
