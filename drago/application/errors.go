package application

import "errors"

const (
	errUnauthorized = "unauthorized"
)

var (
	// ErrUnauthorized :
	ErrUnauthorized = errors.New(errUnauthorized)
)
