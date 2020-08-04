package structs

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

func dashedAlphanumValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^[A-Za-z0-9][A-Za-z0-9_-]*$")
	return re.MatchString(fl.Field().String())
}

type dto struct{}

// Validate : Validates the DTO, returning an error in case its
// attributes are not compliant with the allowed values.
func (d *dto) Validate() error {

	v := validator.New()
	v.RegisterValidation("dashedalphanum", dashedAlphanumValidator)

	err := v.Struct(d)
	if err != nil {
		return errors.Wrap(errors.New("invalid dto"), err.Error())
	}

	return nil
}
