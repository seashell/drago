package validator

import (
	"context"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func dashedAlphanumValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^[a-z0-9][a-z0-9_-]*$")

	return re.MatchString(fl.Field().String())
}

type Validator struct {
	v *validator.Validate
}

func New(ctx context.Context) (*Validator, error) {
	v := validator.New()
	v.RegisterValidation("dashed-alphanumeric", dashedAlphanumValidator)

	return &Validator{
		v: v,
	}, nil
}

func (v Validator) Validate(i interface{}) error {
	return v.v.Struct(i)
}
