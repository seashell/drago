package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/seashell/drago/server/controller"
)

type APIError struct {
	Err       error     `json:"-"`
	Code      int       `json:"-"`
	Text      string    `json:"-"`
	Message   string    `json:"message"`
	Detail    string    `json:"detail"`
	Timestamp time.Time `json:"timestamp"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s", e.Err)
}

func WrapControllerError(e error) *APIError {

	if e == nil {
		return nil
	}

	code, err := generateStatusCode(e)
	if err != nil {
		fmt.Println("Unhandled controller error")
	}

	return &APIError{
		Err:       e,
		Code:      code,
		Text:      http.StatusText(code),
		Message:   errors.Cause(e).Error(),
		Detail:    errors.Unwrap(e).Error(),
		Timestamp: time.Now(),
	}
}

func generateStatusCode(err error) (int, error) {

	dict := map[error]int{
		controller.ErrInternal:       http.StatusInternalServerError,
		controller.ErrNotFound:       http.StatusNotFound,
		controller.ErrInvalidInput:   http.StatusBadRequest,
		controller.ErrNotImplemented: http.StatusNotImplemented,
	}

	if status, ok := dict[errors.Cause(err)]; ok {
		return status, nil
	}

	return http.StatusInternalServerError, errors.New("Unknown controler error")
}
