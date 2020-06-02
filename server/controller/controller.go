package controller

import (
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
}

// New :
func New(ns application.NetworkService, hs application.HostService, is application.InterfaceService, ls application.LinkService) (*Controller, error) {
	return &Controller{
		v:  validator.New(),
		ns: ns,
		hs: hs,
		is: is,
		ls: ls,
	}, nil
}
