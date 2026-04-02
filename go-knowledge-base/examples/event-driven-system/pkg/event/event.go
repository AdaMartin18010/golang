package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event represents a domain event in the event sourcing system
type Event struct {
	ID            uuid.UUID       `json:"id"`
	AggregateID   uuid.UUID       `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	Type          string          `json:"type"`
	Version       int             `json:"version"`
	Data          json.RawMessage `json:"data"`
	Metadata      Metadata        `json:"metadata"`
	Timestamp     time.Time       `json:"timestamp"`
}

// Metadata contains additional event metadata
type Metadata struct {
	CorrelationID string            `json:"correlation_id"`
	CausationID   string            `json:"causation_id"`
	UserID        string            `json:"user_id,omitempty"`
	Service       string            `json:"service"`
	Extra         map[string]string `json:"extra,omitempty"`
}

// NewEvent creates a new domain event
func NewEvent(aggregateID uuid.UUID, aggregateType, eventType string, version int, data interface{}) (*Event, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:            uuid.New(),
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		Type:          eventType,
		Version:       version,
		Data:          dataBytes,
		Metadata: Metadata{
			CorrelationID: uuid.New().String(),
			Service:       "event-driven-system",
		},
		Timestamp: time.Now().UTC(),
	}, nil
}

// UnmarshalData unmarshals event data into the provided struct
func (e *Event) UnmarshalData(v interface{}) error {
	return json.Unmarshal(e.Data, v)
}

// EventStore defines the interface for event storage
type EventStore interface {
	// Append persists events to the store
	Append(ctx context.Context, events ...*Event) error
	
	// GetEvents retrieves events for a specific aggregate
	GetEvents(ctx context.Context, aggregateID uuid.UUID, fromVersion int) ([]*Event, error)
	
	// GetAllEvents retrieves all events after a specific position (for projections)
	GetAllEvents(ctx context.Context, afterPosition int64, batchSize int) ([]*Event, error)
	
	// GetAggregateVersion returns the current version of an aggregate
	GetAggregateVersion(ctx context.Context, aggregateID uuid.UUID) (int, error)
	
	// Subscribe subscribes to events of specific types
	Subscribe(ctx context.Context, eventTypes []string) (<-chan *Event, error)
}

// EventBus defines the interface for event publishing/subscribing
type EventBus interface {
	// Publish publishes an event to the bus
	Publish(ctx context.Context, event *Event) error
	
	// Subscribe subscribes to events matching the filter
	Subscribe(ctx context.Context, filter EventFilter) (<-chan *Event, error)
	
	// Close closes the event bus connection
	Close() error
}

// EventFilter defines criteria for event subscription
type EventFilter struct {
	AggregateTypes []string
	EventTypes     []string
	AggregateIDs   []uuid.UUID
}

// Handler is a function that handles events
type Handler func(ctx context.Context, event *Event) error

// Registry manages event handlers
type Registry struct {
	handlers map[string][]Handler
}

// NewRegistry creates a new event handler registry
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string][]Handler),
	}
}

// Register registers a handler for a specific event type
func (r *Registry) Register(eventType string, handler Handler) {
	r.handlers[eventType] = append(r.handlers[eventType], handler)
}

// Dispatch dispatches an event to all registered handlers
func (r *Registry) Dispatch(ctx context.Context, event *Event) error {
	handlers, ok := r.handlers[event.Type]
	if !ok {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
