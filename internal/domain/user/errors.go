package user

import "errors"

var (
	ErrEmailRequired = errors.New("email is required")
	ErrNameRequired  = errors.New("name is required")
	ErrNotFound      = errors.New("user not found")
	ErrAlreadyExists = errors.New("user already exists")
)
