package structs

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// BaseInput contains information that is common to all operations.
type BaseInput struct {
	Subject string
}

// BaseOutput contains information that is common to all results.
type BaseOutput struct {
}

func dashedAlphanumValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^[A-Za-z0-9][A-Za-z0-9_-]*$")
	return re.MatchString(fl.Field().String())
}

// Validate : Validates a struct/DTO, returning an error in case its
// attributes are not compliant with the allowed values.
func Validate(s interface{}) error {

	v := validator.New()
	v.RegisterValidation("dashedalphanum", dashedAlphanumValidator)

	err := v.Struct(s)
	if err != nil {
		return errors.Wrap(errors.New("invalid struct"), err.Error())
	}

	return nil
}
