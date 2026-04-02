package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"event-driven-system/internal/domain/order"
	"event-driven-system/pkg/cqrs"
	"event-driven-system/pkg/event"
)

// Command definitions

type CreateOrderCommand struct {
	CustomerID uuid.UUID    `json:"customer_id"`
	Items      []order.Item `json:"items"`
}

func (c CreateOrderCommand) CommandName() string { return "CreateOrder" }

type PayOrderCommand struct {
	OrderID       uuid.UUID `json:"order_id"`
	PaymentID     uuid.UUID `json:"payment_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
}

func (c PayOrderCommand) CommandName() string { return "PayOrder" }

type ShipOrderCommand struct {
	OrderID        uuid.UUID `json:"order_id"`
	ShipmentID     uuid.UUID `json:"shipment_id"`
	Carrier        string    `json:"carrier"`
	TrackingNumber string    `json:"tracking_number"`
}

func (c ShipOrderCommand) CommandName() string { return "ShipOrder" }

type CancelOrderCommand struct {
	OrderID uuid.UUID `json:"order_id"`
	Reason  string    `json:"reason"`
}

func (c CancelOrderCommand) CommandName() string { return "CancelOrder" }

// OrderCommandHandler handles order commands
type OrderCommandHandler struct {
	eventStore event.EventStore
	eventBus   event.EventBus
}

// NewOrderCommandHandler creates a new order command handler
func NewOrderCommandHandler(eventStore event.EventStore, eventBus event.EventBus) *OrderCommandHandler {
	return &OrderCommandHandler{
		eventStore: eventStore,
		eventBus:   eventBus,
	}
}

// HandleCreateOrder handles the CreateOrder command
func (h *OrderCommandHandler) HandleCreateOrder(ctx context.Context, cmd CreateOrderCommand) error {
	// Create the order aggregate
	orderAggregate, err := order.NewOrder(cmd.CustomerID, cmd.Items)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Get uncommitted events
	events := orderAggregate.GetUncommittedEvents()

	// Append events to event store
	if err := h.eventStore.Append(ctx, events...); err != nil {
		return fmt.Errorf("failed to append events: %w", err)
	}

	// Publish events to event bus
	for _, evt := range events {
		if err := h.eventBus.Publish(ctx, evt); err != nil {
			// Log error but don't fail the command
			// Event will be published by outbox pattern or retry mechanism
			fmt.Printf("Failed to publish event: %v\n", err)
		}
	}

	orderAggregate.MarkCommitted()
	return nil
}

// HandlePayOrder handles the PayOrder command
func (h *OrderCommandHandler) HandlePayOrder(ctx context.Context, cmd PayOrderCommand) error {
	// Load order aggregate from event store
	orderAggregate, err := h.loadOrder(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	// Execute business logic
	evt, err := orderAggregate.Pay(cmd.PaymentID, cmd.Amount, cmd.PaymentMethod)
	if err != nil {
		return err
	}

	// Append and publish event
	if err := h.eventStore.Append(ctx, evt); err != nil {
		return err
	}

	if err := h.eventBus.Publish(ctx, evt); err != nil {
		fmt.Printf("Failed to publish event: %v\n", err)
	}

	return nil
}

// HandleShipOrder handles the ShipOrder command
func (h *OrderCommandHandler) HandleShipOrder(ctx context.Context, cmd ShipOrderCommand) error {
	orderAggregate, err := h.loadOrder(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	evt, err := orderAggregate.Ship(cmd.ShipmentID, cmd.Carrier, cmd.TrackingNumber)
	if err != nil {
		return err
	}

	if err := h.eventStore.Append(ctx, evt); err != nil {
		return err
	}

	if err := h.eventBus.Publish(ctx, evt); err != nil {
		fmt.Printf("Failed to publish event: %v\n", err)
	}

	return nil
}

// HandleCancelOrder handles the CancelOrder command
func (h *OrderCommandHandler) HandleCancelOrder(ctx context.Context, cmd CancelOrderCommand) error {
	orderAggregate, err := h.loadOrder(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	evt, err := orderAggregate.Cancel(cmd.Reason)
	if err != nil {
		return err
	}

	if err := h.eventStore.Append(ctx, evt); err != nil {
		return err
	}

	if err := h.eventBus.Publish(ctx, evt); err != nil {
		fmt.Printf("Failed to publish event: %v\n", err)
	}

	return nil
}

// loadOrder reconstructs an order aggregate from events
func (h *OrderCommandHandler) loadOrder(ctx context.Context, orderID uuid.UUID) (*order.Aggregate, error) {
	events, err := h.eventStore.GetEvents(ctx, orderID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to load events: %w", err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	orderAggregate := &order.Aggregate{}

	for _, evt := range events {
		switch evt.Type {
		case "OrderCreated":
			var data order.Created
			if err := json.Unmarshal(evt.Data, &data); err != nil {
				return nil, err
			}
			orderAggregate.ApplyCreated(&data)

		case "OrderPaid":
			var data order.Paid
			if err := json.Unmarshal(evt.Data, &data); err != nil {
				return nil, err
			}
			orderAggregate.ApplyPaid(&data)

		case "OrderShipped":
			var data order.Shipped
			if err := json.Unmarshal(evt.Data, &data); err != nil {
				return nil, err
			}
			orderAggregate.ApplyShipped(&data)

		case "OrderDelivered":
			var data order.Delivered
			if err := json.Unmarshal(evt.Data, &data); err != nil {
				return nil, err
			}
			orderAggregate.ApplyDelivered(&data)

		case "OrderCancelled":
			var data order.Cancelled
			if err := json.Unmarshal(evt.Data, &data); err != nil {
				return nil, err
			}
			orderAggregate.ApplyCancelled(&data)
		}
		orderAggregate.Version = evt.Version
	}

	orderAggregate.ID = orderID
	return orderAggregate, nil
}

// RegisterHandlers registers all command handlers
func (h *OrderCommandHandler) RegisterHandlers(bus *cqrs.CommandBus) {
	bus.Register("CreateOrder", func(ctx context.Context, cmd cqrs.Command) error {
		return h.HandleCreateOrder(ctx, cmd.(CreateOrderCommand))
	})
	bus.Register("PayOrder", func(ctx context.Context, cmd cqrs.Command) error {
		return h.HandlePayOrder(ctx, cmd.(PayOrderCommand))
	})
	bus.Register("ShipOrder", func(ctx context.Context, cmd cqrs.Command) error {
		return h.HandleShipOrder(ctx, cmd.(ShipOrderCommand))
	})
	bus.Register("CancelOrder", func(ctx context.Context, cmd cqrs.Command) error {
		return h.HandleCancelOrder(ctx, cmd.(CancelOrderCommand))
	})
}
