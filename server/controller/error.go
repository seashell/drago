package controller

import "github.com/pkg/errors"

var (
	ErrInvalidInput   = errors.New("Invalid input")
	ErrInternal       = errors.New("Internal error")
	ErrNotImplemented = errors.New("Not implemented")
	ErrUnauthorized   = errors.New("Not authorized")
	ErrNotFound       = errors.New("Resource not found")
)
