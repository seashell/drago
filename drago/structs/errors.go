package structs

import (
	"fmt"

	"errors"
)

const (
	errPermissionDenied       = "Permission denied"
	errTokenNotFound          = "ACL Token not found"
	errACLDisabled            = "ACL disabled"
	errACLAlreadyBootstrapped = "ACL already bootstrapped"
	errInvalidInput           = "Invalid input"
	errNotFound               = "Resource not found"
	errInternal               = "Internal error"
)

var (
	// ErrPermissionDenied :
	ErrPermissionDenied = errors.New(errPermissionDenied)

	// ErrACLAlreadyBootstrapped ...
	ErrACLAlreadyBootstrapped = errors.New(errACLAlreadyBootstrapped)

	// ErrACLDisabled ...
	ErrACLDisabled = errors.New(errACLDisabled)

	// ErrInternal ...
	ErrInternal = errors.New(errInternal)

	// ErrInvalidInput ...
	ErrInvalidInput = errors.New(errInvalidInput)

	// ErrNotFound ...
	ErrNotFound = errors.New(errNotFound)
)

// Error :
type Error struct {
	Message string
}

// NewError ...
func NewError(base error, extra ...interface{}) error {
	msg := base.Error()
	for _, v := range extra {
		msg = fmt.Sprintf("%s : %v", msg, v)
	}
	return &Error{
		Message: msg,
	}
}

func (e Error) Error() string {
	return e.Message
}
