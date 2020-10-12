package handler

import (
	"net/http"

	application "github.com/seashell/drago/drago/application"
	structs "github.com/seashell/drago/drago/application/structs"
)

// NetworkHandlerAdapter :
type NetworkHandlerAdapter struct {
	networkService application.NetworkService
}

// NewNetworkHandlerAdapter :
func NewNetworkHandlerAdapter(ts application.NetworkService) *NetworkHandlerAdapter {
	return &NetworkHandlerAdapter{
		networkService: ts,
	}
}

// Handle :
func (a *NetworkHandlerAdapter) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case http.MethodGet:
		return a.handleGet(rw, req)
	case http.MethodPost:
		return a.handlePost(rw, req)
	case http.MethodDelete:
		return a.handleDelete(rw, req)
	case http.MethodPatch:
		return a.handlePatch(rw, req)
	default:
		return nil, NewError(405, ErrInvalidMethod)
	}
}

func (a *NetworkHandlerAdapter) handleGet(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)
	if len(params) > 1 {
		return nil, NewError(404, ErrNotFound)
	}

	id := params[0]
	if id == "" {
		return a.handleList(rw, req)
	}

	out, err := a.networkService.GetByID(req.Context(), &structs.NetworkGetInput{
		BaseInput: baseInputFromReq(req),
		ID:        id,
	})
	if err != nil {
		return nil, NewError(404, ErrNotFound)
	}

	return out, nil
}

func (a *NetworkHandlerAdapter) handlePost(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	in := &structs.NetworkCreateInput{}
	err := parseBody(req.Body, in)
	if err != nil {
		return nil, NewError(400, ErrBadRequest, err)
	}
	in.BaseInput = baseInputFromReq(req)

	out, err := a.networkService.Create(req.Context(), in)
	if err != nil {
		return nil, NewError(500, ErrInternal, err)
	}

	return out, nil
}

func (a *NetworkHandlerAdapter) handleDelete(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)
	if len(params) > 1 {
		return nil, NewError(404, ErrNotFound)
	}

	id := params[0]
	if id == "" {
		return nil, NewError(400, ErrBadRequest)
	}

	_, err := a.networkService.Delete(req.Context(), &structs.NetworkDeleteInput{
		BaseInput: baseInputFromReq(req),
		ID:        id,
	})
	if err != nil {
		return nil, NewError(404, ErrNotFound)
	}

	return nil, nil
}

func (a *NetworkHandlerAdapter) handlePatch(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return nil, NewError(501, ErrNotImplemented)
}

func (a *NetworkHandlerAdapter) handleList(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	in := &structs.NetworkListInput{
		BaseInput: baseInputFromReq(req),
	}

	out, err := a.networkService.List(req.Context(), in)
	if err != nil {
		return nil, err
	}

	return out, nil
}
