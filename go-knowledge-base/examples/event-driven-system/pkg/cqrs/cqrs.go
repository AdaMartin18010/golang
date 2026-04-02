package cqrs

import (
	"context"
)

// Command represents a command in the CQRS pattern
type Command interface {
	CommandName() string
}

// CommandHandler handles commands
type CommandHandler[C Command] interface {
	Handle(ctx context.Context, cmd C) error
}

// Query represents a query in the CQRS pattern
type Query interface {
	QueryName() string
}

// QueryResult is the result of a query
type QueryResult interface{}

// QueryHandler handles queries
type QueryHandler[Q Query, R QueryResult] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

// CommandBus routes commands to their handlers
type CommandBus struct {
	handlers map[string]func(ctx context.Context, cmd Command) error
}

// NewCommandBus creates a new command bus
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]func(ctx context.Context, cmd Command) error),
	}
}

// Register registers a command handler
func (b *CommandBus) Register(commandName string, handler func(ctx context.Context, cmd Command) error) {
	b.handlers[commandName] = handler
}

// Dispatch dispatches a command to its handler
func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	handler, ok := b.handlers[cmd.CommandName()]
	if !ok {
		return ErrCommandHandlerNotFound
	}
	return handler(ctx, cmd)
}

// QueryBus routes queries to their handlers
type QueryBus struct {
	handlers map[string]func(ctx context.Context, query Query) (QueryResult, error)
}

// NewQueryBus creates a new query bus
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[string]func(ctx context.Context, query Query) (QueryResult, error)),
	}
}

// Register registers a query handler
func (b *QueryBus) Register(queryName string, handler func(ctx context.Context, query Query) (QueryResult, error)) {
	b.handlers[queryName] = handler
}

// Dispatch dispatches a query to its handler
func (b *QueryBus) Dispatch(ctx context.Context, query Query) (QueryResult, error) {
	handler, ok := b.handlers[query.QueryName()]
	if !ok {
		return nil, ErrQueryHandlerNotFound
	}
	return handler(ctx, query)
}
