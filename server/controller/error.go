package controller

import "github.com/pkg/errors"

var (
	// Invalid input error
	ErrInvalidInput = errors.New("Invalid input")
	// Internal error
	ErrInternal = errors.New("Internal error")
	// Not implemented error
	ErrNotImplemented = errors.New("Not implemented")
	// Unauthorized error
	ErrUnauthorized = errors.New("Not authorized")
	// Not found error
	ErrNotFound = errors.New("Resource not found")
)
