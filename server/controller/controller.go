package controller

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/seashell/drago/server/application"
)

// Controller :
type Controller struct {
	v  *validator.Validate
	ns application.NetworkService
	hs application.HostService
	is application.InterfaceService
	ls application.LinkService
	ss application.SynchronizationService
	ts application.TokenService
}

const (
	dashedAlphaNumericRegexString = "^[A-Za-z0-9][A-Za-z0-9_-]*$"
)

var (
	dashedAlphaNumericRegex = regexp.MustCompile(dashedAlphaNumericRegexString)
)

func dashedAlphanumValidator(fl validator.FieldLevel) bool {
	return dashedAlphaNumericRegex.MatchString(fl.Field().String())
}

// New :
func New(ns application.NetworkService,
	hs application.HostService,
	is application.InterfaceService,
	ls application.LinkService,
	ss application.SynchronizationService,
	ts application.TokenService) (*Controller, error) {


	//implement dashedAlphanumeric validator
	v := validator.New()
	v.RegisterValidation("dashedalphanum", dashedAlphanumValidator)


	return &Controller{
		v:  v,
		ns: ns,
		hs: hs,
		is: is,
		ls: ls,
		ss: ss,
		ts: ts,
	}, nil
}

