package http

import (
	"net/http"

	structs "github.com/seashell/drago/drago/structs"
	rpc "github.com/seashell/drago/pkg/rpc"
)

// ACLHandler :
type ACLHandler struct {
	rpcClient *rpc.Client
}

// NewACLHandler :
func NewACLHandler(rpcClient *rpc.Client) *ACLHandler {
	return &ACLHandler{
		rpcClient: rpcClient,
	}
}

// Handle :
func (h *ACLHandler) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)
	if len(params) > 1 {
		return nil, NewCodedError(404, ErrNotFound)
	}

	switch params[0] {
	case "bootstrap":
		return h.handleBootstrap(rw, req)
	default:
		return nil, NewCodedError(404, "Not found")
	}

}

func (h *ACLHandler) handleBootstrap(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	if req.Method != "POST" {
		return nil, NewCodedError(405, ErrMethodNotAllowed)
	}

	args := structs.ACLBootstrapRequest{
		WriteRequest: parseWriteRequestOptions(req),
	}

	var out structs.ACLTokenUpsertResponse
	if err := h.rpcClient.Call("ACL.BootstrapACL", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return out.ACLToken, nil
}
