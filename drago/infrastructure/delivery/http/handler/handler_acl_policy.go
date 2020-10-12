package handler

import (
	"net/http"

	application "github.com/seashell/drago/drago/application"
	structs "github.com/seashell/drago/drago/application/structs"
)

// ACLPolicyHandlerAdapter :
type ACLPolicyHandlerAdapter struct {
	policyService application.ACLPolicyService
}

// NewACLPolicyHandlerAdapter :
func NewACLPolicyHandlerAdapter(ps application.ACLPolicyService) *ACLPolicyHandlerAdapter {
	return &ACLPolicyHandlerAdapter{
		policyService: ps,
	}
}

// Handle :
func (a *ACLPolicyHandlerAdapter) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case http.MethodGet:
		return a.handleGet(rw, req)
	case http.MethodPost:
		return nil, NewError(501, ErrNotImplemented)
	case http.MethodDelete:
		return nil, NewError(501, ErrNotImplemented)
	case http.MethodPatch:
		return nil, NewError(501, ErrNotImplemented)
	default:
		return nil, NewError(405, ErrInvalidMethod)
	}
}

func (a *ACLPolicyHandlerAdapter) handleGet(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)

	if len(params) > 1 {
		return nil, NewError(404, ErrNotFound)
	}

	name := params[0]
	if name == "" {
		return a.handleList(rw, req)
	}

	out, err := a.policyService.GetByName(req.Context(), &structs.ACLPolicyGetInput{Name: name})
	if err != nil {
		return nil, NewError(404, ErrNotFound)
	}

	return out, nil
}

func (a *ACLPolicyHandlerAdapter) handleList(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	in := &structs.ACLPolicyListInput{}

	out, err := a.policyService.List(req.Context(), in)
	if err != nil {
		return nil, err
	}

	return out, nil
}
