package structs

import (
	"context"

	"github.com/pkg/errors"
	"github.com/seashell/drago/pkg/validator"
)

type Filters map[string][]string

func (f Filters) Get(k string) []string {
	if v, ok := f[k]; ok {
		return v
	}
	return []string{}
}

func (f Filters) Add(k, v string) {
	f[k] = append(f[k], []string{v}...)
}

// QueryOptions contains information that is common to all read requests.
type QueryOptions struct {
	AuthToken string
	Filters   Filters
}

// WriteRequest contains information that is common to all write requests.
type WriteRequest struct {
	AuthToken string
}

// Response contains information that is common to all responses.
type Response struct {
}

// GenericRequest is used to request where no
// specific information is needed.
type GenericRequest struct {
	QueryOptions
}

// GenericResponse is used to respond to a request where no
// specific response information is needed.
type GenericResponse struct {
	Response
}

// Validate validates a struct/DTO, returning an error in case its
// attributes are not compliant with the allowed values.
func Validate(s interface{}) error {

	ctx := context.Background()
	v, err := validator.New(ctx)
	if err != nil {
		return err
	}

	err = v.Validate(s)
	if err != nil {
		return errors.Wrap(errors.New("invalid struct"), err.Error())
	}

	return nil
}
