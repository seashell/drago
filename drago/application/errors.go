package application

import (
	"errors"
	"fmt"
	"strings"
)

const (
	errUnknownMethod    = "unknown method"
	errPermissionDenied = "permission denied"
)

var (
	ErrUnknownMethod    = errors.New(errUnknownMethod)
	ErrPermissionDenied = errors.New(errPermissionDenied)
)

func NewError(baseError string, extra ...interface{}) error {
	return fmt.Errorf("%s %q", baseError, extra)
}

func IsErrorType(err error, e error) bool {
	return err != nil && strings.Contains(err.Error(), e.Error())
}
