package handler

import (
	"net/http"

	"github.com/seashell/drago/drago/application/structs"

	application "github.com/seashell/drago/drago/application"
)

const (
	aclBootstrapPath = "bootstrap"
)

// ACLHandlerAdapter :
type ACLHandlerAdapter struct {
	aclService application.ACLService
}

// NewACLHandlerAdapter :
func NewACLHandlerAdapter(acl application.ACLService) *ACLHandlerAdapter {
	return &ACLHandlerAdapter{
		aclService: acl,
	}
}

// Handle :
func (a *ACLHandlerAdapter) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)

	switch params[0] {
	case aclBootstrapPath:
		return a.handleBootstrap(rw, req)
	default:
		return nil, NewError(404, ErrNotFound)
	}

}

func (a *ACLHandlerAdapter) handleBootstrap(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)

	if len(params) > 1 {
		return nil, NewError(404, ErrNotFound)
	}

	if req.Method != http.MethodPost {
		return nil, NewError(405, ErrMethodNotAllowed)
	}

	out, err := a.aclService.Bootstrap(req.Context(), &structs.ACLBootstrapInput{})
	if err != nil {
		return nil, NewError(500, ErrInternal, err)
	}

	return out, nil
}
