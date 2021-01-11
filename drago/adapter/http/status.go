package http

import (
	"net/http"

	"github.com/seashell/drago/drago/structs"
	"github.com/seashell/drago/pkg/rpc"
)

// StatusHandler is used to check on server status
type StatusHandler struct {
	rpcClient *rpc.Client
}

// NewStatusHandler :
func NewStatusHandler(rpcClient *rpc.Client) *StatusHandler {
	return &StatusHandler{
		rpcClient: rpcClient,
	}
}

// Handle :
func (h *StatusHandler) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case "GET":
		return h.handleGet(rw, req)
	default:
		return nil, NewCodedError(405, ErrMethodNotAllowed)
	}
}

func (h *StatusHandler) handleGet(resp http.ResponseWriter, req *http.Request) (interface{}, error) {

	var args structs.GenericRequest

	var out structs.GenericResponse
	if err := h.rpcClient.Call("Status.Ping", &args, &out); err != nil {
		return nil, err
	}

	return out, nil
}
