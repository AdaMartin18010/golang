package eventstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"event-driven-system/pkg/event"
)

// PostgresEventStore implements EventStore using PostgreSQL
type PostgresEventStore struct {
	db *sql.DB
}

// NewPostgresEventStore creates a new PostgreSQL event store
func NewPostgresEventStore(db *sql.DB) *PostgresEventStore {
	return &PostgresEventStore{db: db}
}

// Initialize creates the required tables
func (s *PostgresEventStore) Initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS events (
		id UUID PRIMARY KEY,
		aggregate_id UUID NOT NULL,
		aggregate_type VARCHAR(100) NOT NULL,
		type VARCHAR(100) NOT NULL,
		version INTEGER NOT NULL,
		data JSONB NOT NULL,
		metadata JSONB,
		timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		position BIGSERIAL
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_events_aggregate_version 
		ON events(aggregate_id, version);
	
	CREATE INDEX IF NOT EXISTS idx_events_aggregate_id 
		ON events(aggregate_id);
	
	CREATE INDEX IF NOT EXISTS idx_events_type 
		ON events(type);
	
	CREATE INDEX IF NOT EXISTS idx_events_position 
		ON events(position);

	CREATE TABLE IF NOT EXISTS aggregates (
		id UUID PRIMARY KEY,
		type VARCHAR(100) NOT NULL,
		version INTEGER NOT NULL DEFAULT 0,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

// Append persists events to the store
func (s *PostgresEventStore) Append(ctx context.Context, events ...*event.Event) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, evt := range events {
		// Insert event
		metadataJSON, _ := json.Marshal(evt.Metadata)
		
		_, err := tx.ExecContext(ctx, `
			INSERT INTO events (id, aggregate_id, aggregate_type, type, version, data, metadata, timestamp)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, evt.ID, evt.AggregateID, evt.AggregateType, evt.Type, evt.Version, evt.Data, metadataJSON, evt.Timestamp)
		
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				return fmt.Errorf("concurrency conflict: %w", err)
			}
			return err
		}

		// Update aggregate version
		_, err = tx.ExecContext(ctx, `
			INSERT INTO aggregates (id, type, version, updated_at)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET
				version = EXCLUDED.version,
				updated_at = EXCLUDED.updated_at
		`, evt.AggregateID, evt.AggregateType, evt.Version, time.Now())
		
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetEvents retrieves events for a specific aggregate
func (s *PostgresEventStore) GetEvents(ctx context.Context, aggregateID uuid.UUID, fromVersion int) ([]*event.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, aggregate_id, aggregate_type, type, version, data, metadata, timestamp
		FROM events
		WHERE aggregate_id = $1 AND version > $2
		ORDER BY version ASC
	`, aggregateID, fromVersion)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanEvents(rows)
}

// GetAllEvents retrieves all events after a specific position
func (s *PostgresEventStore) GetAllEvents(ctx context.Context, afterPosition int64, batchSize int) ([]*event.Event, error) {
	if batchSize <= 0 {
		batchSize = 100
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, aggregate_id, aggregate_type, type, version, data, metadata, timestamp
		FROM events
		WHERE position > $1
		ORDER BY position ASC
		LIMIT $2
	`, afterPosition, batchSize)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanEvents(rows)
}

// GetAggregateVersion returns the current version of an aggregate
func (s *PostgresEventStore) GetAggregateVersion(ctx context.Context, aggregateID uuid.UUID) (int, error) {
	var version int
	err := s.db.QueryRowContext(ctx, `
		SELECT version FROM aggregates WHERE id = $1
	`, aggregateID).Scan(&version)
	
	if err == sql.ErrNoRows {
		return 0, nil
	}
	
	return version, err
}

// Subscribe subscribes to events (simplified implementation)
func (s *PostgresEventStore) Subscribe(ctx context.Context, eventTypes []string) (<-chan *event.Event, error) {
	// In production, this would use LISTEN/NOTIFY or a message bus
	return nil, fmt.Errorf("not implemented: use event bus for subscriptions")
}

func (s *PostgresEventStore) scanEvents(rows *sql.Rows) ([]*event.Event, error) {
	var events []*event.Event

	for rows.Next() {
		var evt event.Event
		var metadataJSON []byte

		err := rows.Scan(
			&evt.ID,
			&evt.AggregateID,
			&evt.AggregateType,
			&evt.Type,
			&evt.Version,
			&evt.Data,
			&metadataJSON,
			&evt.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &evt.Metadata)
		}

		events = append(events, &evt)
	}

	return events, rows.Err()
}
