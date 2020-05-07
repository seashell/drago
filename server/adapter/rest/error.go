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
	Detail    string    `json:"detail"`
	Timestamp time.Time `json:"timestamp"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s", e.Err)
}

func WrapControllerError(e error) *APIError {

	code := generateStatusCode(e)

	return &APIError{
		Err:       e,
		Code:      code,
		Text:      http.StatusText(code),
		Detail:    errors.Cause(e).Error(),
		Timestamp: time.Now(),
	}
}

func generateStatusCode(err error) int {

	dict := map[error]int{
		controller.ErrNotFound:       http.StatusNotFound,
		controller.ErrInvalidInput:   http.StatusBadRequest,
		controller.ErrNotImplemented: http.StatusNotImplemented,
	}

	if status, ok := dict[errors.Cause(err)]; ok {
		return status
	}

	return http.StatusInternalServerError
}
