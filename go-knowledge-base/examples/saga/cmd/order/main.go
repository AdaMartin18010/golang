package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// OrderService 订单服务 - Saga 编排器
type OrderService struct {
	paymentClient   *PaymentClient
	inventoryClient *InventoryClient
	eventBus        EventBus
	orders          map[string]*Order
}

type Order struct {
	ID        string
	UserID    string
	Items     []OrderItem
	Status    OrderStatus
	CreatedAt time.Time
}

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusReserved  OrderStatus = "INVENTORY_RESERVED"
	StatusPaid      OrderStatus = "PAID"
	StatusCompleted OrderStatus = "COMPLETED"
	StatusFailed    OrderStatus = "FAILED"
	StatusCancelled OrderStatus = "CANCELLED"
)

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserID string      `json:"user_id"`
	Items  []OrderItem `json:"items"`
}

// CreateOrder 创建订单 (Saga 开始)
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
	order := &Order{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Items:     req.Items,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
	s.orders[order.ID] = order

	// Saga Step 1: 预留库存
	err := s.inventoryClient.ReserveInventory(ctx, ReserveInventoryRequest{
		OrderID: order.ID,
		Items:   req.Items,
	})
	if err != nil {
		order.Status = StatusFailed
		return nil, fmt.Errorf("reserve inventory failed: %w", err)
	}
	order.Status = StatusReserved

	// Saga Step 2: 处理支付
	total := calculateTotal(req.Items)
	err = s.paymentClient.ProcessPayment(ctx, ProcessPaymentRequest{
		OrderID: order.ID,
		UserID:  req.UserID,
		Amount:  total,
	})
	if err != nil {
		// 补偿: 释放库存
		order.Status = StatusFailed
		s.compensateInventory(ctx, order.ID)
		return nil, fmt.Errorf("payment failed: %w", err)
	}
	order.Status = StatusPaid

	// Saga Step 3: 确认库存
	err = s.inventoryClient.ConfirmReservation(ctx, order.ID)
	if err != nil {
		// 补偿: 退款 + 释放库存
		order.Status = StatusFailed
		s.compensatePayment(ctx, order.ID)
		s.compensateInventory(ctx, order.ID)
		return nil, fmt.Errorf("confirm inventory failed: %w", err)
	}

	order.Status = StatusCompleted

	// 发布订单完成事件
	s.eventBus.Publish(ctx, "order.completed", OrderCompletedEvent{
		OrderID: order.ID,
		UserID:  order.UserID,
	})

	return order, nil
}

// compensateInventory 补偿: 释放库存
func (s *OrderService) compensateInventory(ctx context.Context, orderID string) {
	if err := s.inventoryClient.ReleaseReservation(ctx, orderID); err != nil {
		log.Printf("Failed to release inventory for order %s: %v", orderID, err)
		// 记录到待处理队列，人工介入
	}
}

// compensatePayment 补偿: 退款
func (s *OrderService) compensatePayment(ctx context.Context, orderID string) {
	if err := s.paymentClient.RefundPayment(ctx, orderID); err != nil {
		log.Printf("Failed to refund payment for order %s: %v", orderID, err)
		// 记录到待处理队列，人工介入
	}
}

func calculateTotal(items []OrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// HTTP Handlers

func (s *OrderService) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := s.CreateOrder(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func main() {
	service := &OrderService{
		orders: make(map[string]*Order),
	}

	http.HandleFunc("/orders", service.handleCreateOrder)

	log.Println("Order service starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Client interfaces
type PaymentClient struct{}
type InventoryClient struct{}
type EventBus struct{}
type ReserveInventoryRequest struct {
	OrderID string
	Items   []OrderItem
}
type ProcessPaymentRequest struct {
	OrderID, UserID string
	Amount          float64
}
type OrderCompletedEvent struct{ OrderID, UserID string }

func (c *PaymentClient) ProcessPayment(ctx context.Context, req ProcessPaymentRequest) error {
	// 调用支付服务
	return nil
}

func (c *PaymentClient) RefundPayment(ctx context.Context, orderID string) error {
	return nil
}

func (c *InventoryClient) ReserveInventory(ctx context.Context, req ReserveInventoryRequest) error {
	// 调用库存服务
	return nil
}

func (c *InventoryClient) ConfirmReservation(ctx context.Context, orderID string) error {
	return nil
}

func (c *InventoryClient) ReleaseReservation(ctx context.Context, orderID string) error {
	return nil
}

func (e *EventBus) Publish(ctx context.Context, event string, data interface{}) error {
	return nil
}
