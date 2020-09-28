package domain

import "errors"

const (
	errNotFound = "not found"
)

var (
	ErrNotFound = errors.New(errNotFound)
)
