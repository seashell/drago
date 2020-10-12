package handler

import (
	"net/http"

	application "github.com/seashell/drago/drago/application"
	structs "github.com/seashell/drago/drago/application/structs"
)

// ACLTokenHandlerAdapter :
type ACLTokenHandlerAdapter struct {
	tokenService application.ACLTokenService
}

// NewACLTokenHandlerAdapter :
func NewACLTokenHandlerAdapter(ts application.ACLTokenService) *ACLTokenHandlerAdapter {
	return &ACLTokenHandlerAdapter{
		tokenService: ts,
	}
}

// Handle :
func (a *ACLTokenHandlerAdapter) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
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

func (a *ACLTokenHandlerAdapter) handleGet(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)

	if len(params) > 1 {
		return nil, NewError(404, ErrNotFound)
	}

	id := params[0]
	if id == "" {
		return a.handleList(rw, req)
	}

	out, err := a.tokenService.GetByID(req.Context(), &structs.ACLTokenGetInput{
		BaseInput: baseInputFromReq(req),
		ID:        id,
	})
	if err != nil {
		return nil, NewError(404, ErrNotFound)
	}

	return out, nil
}

func (a *ACLTokenHandlerAdapter) handlePost(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	in := &structs.ACLTokenCreateInput{}
	err := parseBody(req.Body, in)
	if err != nil {
		return nil, NewError(400, ErrBadRequest, err)
	}
	in.BaseInput = baseInputFromReq(req)

	out, err := a.tokenService.Create(req.Context(), in)
	if err != nil {
		return nil, NewError(500, ErrInternal, err)
	}

	return out, nil
}

func (a *ACLTokenHandlerAdapter) handleDelete(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)
	if len(params) > 1 {
		return nil, NewError(404, ErrNotFound)
	}

	id := params[0]
	if id == "" {
		return nil, NewError(400, ErrBadRequest)
	}

	_, err := a.tokenService.Delete(req.Context(), &structs.ACLTokenDeleteInput{
		BaseInput: baseInputFromReq(req),
		ID:        id,
	})
	if err != nil {
		return nil, NewError(404, ErrNotFound)
	}

	return nil, nil
}

func (a *ACLTokenHandlerAdapter) handlePatch(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return nil, NewError(501, ErrNotImplemented)
}

func (a *ACLTokenHandlerAdapter) handleList(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	in := &structs.ACLTokenListInput{
		BaseInput: baseInputFromReq(req),
	}

	out, err := a.tokenService.List(req.Context(), in)
	if err != nil {
		return nil, err
	}

	return out, nil
}
