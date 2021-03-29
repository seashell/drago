package api

import "fmt"

type CodedError struct {
	Message string
	Code    int
}

func (e CodedError) Error() string {
	return fmt.Sprintf("%d (%s)", e.Code, e.Message)
}
