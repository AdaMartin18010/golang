package cqrs

import "errors"

var (
	ErrCommandHandlerNotFound = errors.New("command handler not found")
	ErrQueryHandlerNotFound   = errors.New("query handler not found")
	ErrConcurrencyConflict    = errors.New("concurrency conflict detected")
	ErrInvalidCommand         = errors.New("invalid command")
	ErrInvalidQuery           = errors.New("invalid query")
)
