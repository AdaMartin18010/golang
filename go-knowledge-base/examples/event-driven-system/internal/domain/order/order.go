package order

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"event-driven-system/pkg/event"
)

var (
	ErrOrderAlreadyPaid    = errors.New("order is already paid")
	ErrOrderAlreadyShipped = errors.New("order is already shipped")
	ErrOrderCancelled      = errors.New("order is cancelled")
	ErrEmptyOrder          = errors.New("order must contain at least one item")
	ErrInvalidQuantity     = errors.New("quantity must be greater than 0")
)

// Aggregate represents the order aggregate root
type Aggregate struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Items      []Item
	Status     Status
	Total      float64
	Version    int
	Changes    []*event.Event
}

// Status represents order status
type Status string

const (
	StatusPending   Status = "pending"
	StatusPaid      Status = "paid"
	StatusShipped   Status = "shipped"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
)

// Item represents an order item
type Item struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
}

// Events

type Created struct {
	OrderID    uuid.UUID `json:"order_id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Items      []Item    `json:"items"`
	Total      float64   `json:"total"`
	CreatedAt  time.Time `json:"created_at"`
}

type Paid struct {
	OrderID       uuid.UUID `json:"order_id"`
	PaymentID     uuid.UUID `json:"payment_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	PaidAt        time.Time `json:"paid_at"`
}

type Shipped struct {
	OrderID        uuid.UUID `json:"order_id"`
	ShipmentID     uuid.UUID `json:"shipment_id"`
	Carrier        string    `json:"carrier"`
	TrackingNumber string    `json:"tracking_number"`
	ShippedAt      time.Time `json:"shipped_at"`
}

type Delivered struct {
	OrderID     uuid.UUID `json:"order_id"`
	DeliveredAt time.Time `json:"delivered_at"`
}

type Cancelled struct {
	OrderID     uuid.UUID `json:"order_id"`
	Reason      string    `json:"reason"`
	CancelledAt time.Time `json:"cancelled_at"`
}

// NewOrder creates a new order aggregate
func NewOrder(customerID uuid.UUID, items []Item) (*Aggregate, error) {
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}

	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
	}

	order := &Aggregate{
		ID:         uuid.New(),
		CustomerID: customerID,
		Status:     StatusPending,
	}

	// Calculate total
	var total float64
	for i := range items {
		items[i].TotalPrice = float64(items[i].Quantity) * items[i].UnitPrice
		total += items[i].TotalPrice
	}
	order.Total = total

	// Create event
	evt := &Created{
		OrderID:    order.ID,
		CustomerID: customerID,
		Items:      items,
		Total:      total,
		CreatedAt:  time.Now().UTC(),
	}

	eventData, err := event.NewEvent(order.ID, "order", "OrderCreated", 1, evt)
	if err != nil {
		return nil, err
	}

	order.Changes = append(order.Changes, eventData)
	order.Items = items

	return order, nil
}

// ApplyCreated applies the OrderCreated event
func (a *Aggregate) ApplyCreated(evt *Created) {
	a.ID = evt.OrderID
	a.CustomerID = evt.CustomerID
	a.Items = evt.Items
	a.Total = evt.Total
	a.Status = StatusPending
}

// Pay processes payment for the order
func (a *Aggregate) Pay(paymentID uuid.UUID, amount float64, paymentMethod string) (*event.Event, error) {
	if a.Status == StatusPaid {
		return nil, ErrOrderAlreadyPaid
	}
	if a.Status == StatusCancelled {
		return nil, ErrOrderCancelled
	}

	evt := &Paid{
		OrderID:       a.ID,
		PaymentID:     paymentID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		PaidAt:        time.Now().UTC(),
	}

	return event.NewEvent(a.ID, "order", "OrderPaid", a.Version+1, evt)
}

// ApplyPaid applies the OrderPaid event
func (a *Aggregate) ApplyPaid(evt *Paid) {
	a.Status = StatusPaid
}

// Ship marks the order as shipped
func (a *Aggregate) Ship(shipmentID uuid.UUID, carrier, trackingNumber string) (*event.Event, error) {
	if a.Status != StatusPaid {
		return nil, errors.New("order must be paid before shipping")
	}
	if a.Status == StatusShipped {
		return nil, ErrOrderAlreadyShipped
	}

	evt := &Shipped{
		OrderID:        a.ID,
		ShipmentID:     shipmentID,
		Carrier:        carrier,
		TrackingNumber: trackingNumber,
		ShippedAt:      time.Now().UTC(),
	}

	return event.NewEvent(a.ID, "order", "OrderShipped", a.Version+1, evt)
}

// ApplyShipped applies the OrderShipped event
func (a *Aggregate) ApplyShipped(evt *Shipped) {
	a.Status = StatusShipped
}

// Deliver marks the order as delivered
func (a *Aggregate) Deliver() (*event.Event, error) {
	if a.Status != StatusShipped {
		return nil, errors.New("order must be shipped before delivery")
	}

	evt := &Delivered{
		OrderID:     a.ID,
		DeliveredAt: time.Now().UTC(),
	}

	return event.NewEvent(a.ID, "order", "OrderDelivered", a.Version+1, evt)
}

// ApplyDelivered applies the OrderDelivered event
func (a *Aggregate) ApplyDelivered(evt *Delivered) {
	a.Status = StatusDelivered
}

// Cancel cancels the order
func (a *Aggregate) Cancel(reason string) (*event.Event, error) {
	if a.Status == StatusPaid || a.Status == StatusShipped {
		return nil, errors.New("cannot cancel paid or shipped order")
	}
	if a.Status == StatusCancelled {
		return nil, ErrOrderCancelled
	}

	evt := &Cancelled{
		OrderID:     a.ID,
		Reason:      reason,
		CancelledAt: time.Now().UTC(),
	}

	return event.NewEvent(a.ID, "order", "OrderCancelled", a.Version+1, evt)
}

// ApplyCancelled applies the OrderCancelled event
func (a *Aggregate) ApplyCancelled(evt *Cancelled) {
	a.Status = StatusCancelled
}

// GetUncommittedEvents returns uncommitted events
func (a *Aggregate) GetUncommittedEvents() []*event.Event {
	return a.Changes
}

// MarkCommitted marks events as committed
func (a *Aggregate) MarkCommitted() {
	a.Changes = nil
}
