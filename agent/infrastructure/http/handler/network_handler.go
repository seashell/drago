package handler

import (
	application "github.com/seashell/drago/agent/application"
	http "github.com/seashell/drago/agent/infrastructure/http"
	"github.com/seashell/drago/drago/application/structs"
	log "github.com/seashell/drago/pkg/log"
)

type NetworkHandlerAdapter struct {
	http.BaseHandlerAdapter
	service application.NetworkService
	logger  log.Logger
}

func NewNetworkHandlerAdapter(service application.NetworkService, logger log.Logger) *NetworkHandlerAdapter {

	a := &NetworkHandlerAdapter{
		service: service,
	}

	a.logger = logger

	a.RegisterHandlerFunc("GET", "/:id", a.getNetwork)
	a.RegisterHandlerFunc("POST", "/", a.createNetwork)
	a.RegisterHandlerFunc("PATCH", "/:id", a.updateNetwork)
	a.RegisterHandlerFunc("DELETE", "/:id", a.deleteNetwork)
	a.RegisterHandlerFunc("GET", "/", a.listNetworks)

	return a
}

func (a *NetworkHandlerAdapter) getNetwork(resp http.Response, req *http.Request) (interface{}, error) {

	id := req.Params["id"]

	in := &structs.GetNetworkInput{
		ID: &id,
	}

	out, err := a.service.Get(in)
	if err != nil {
		return nil, http.NewError(500, err.Error()) // TODO: get more accurate status code from err
	}

	return out, nil
}

func (a *NetworkHandlerAdapter) createNetwork(resp http.Response, req *http.Request) (interface{}, error) {
	return nil, http.NewError(501, http.ErrNotImplemented)
}

func (a *NetworkHandlerAdapter) updateNetwork(resp http.Response, req *http.Request) (interface{}, error) {
	return nil, http.NewError(501, http.ErrNotImplemented)
}

func (a *NetworkHandlerAdapter) deleteNetwork(resp http.Response, req *http.Request) (interface{}, error) {
	return nil, http.NewError(501, http.ErrNotImplemented)
}

func (a *NetworkHandlerAdapter) listNetworks(resp http.Response, req *http.Request) (interface{}, error) {
	return nil, http.NewError(501, http.ErrNotImplemented)
}
