package structs

import (
	"fmt"
)

// Error :
type Error struct {
	Message string `json:"message"`
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
