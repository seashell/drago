package http

import (
	"net/http"

	structs "github.com/seashell/drago/drago/structs"
	rpc "github.com/seashell/drago/pkg/rpc"
)

// InterfaceHandler :
type InterfaceHandler struct {
	rpcClient *rpc.Client
}

// NewInterfaceHandler :
func NewInterfaceHandler(rpcClient *rpc.Client) *InterfaceHandler {
	return &InterfaceHandler{
		rpcClient: rpcClient,
	}
}

// Handle :
func (h *InterfaceHandler) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)
	if len(params) > 1 {
		return nil, NewCodedError(404, ErrNotFound)
	}

	ifaceID := params[0]

	switch req.Method {
	case "GET":
		return h.handleGet(rw, req, ifaceID)
	case "PUT", "POST":
		return h.handlePost(rw, req, ifaceID)
	case "DELETE":
		return h.handleDelete(rw, req, ifaceID)
	default:
		return nil, NewCodedError(405, ErrMethodNotAllowed)
	}
}

func (h *InterfaceHandler) handleGet(rw http.ResponseWriter, req *http.Request, ifaceID string) (interface{}, error) {

	if ifaceID == "" {
		return h.handleList(rw, req)
	}

	args := structs.NetworkSpecificRequest{
		QueryOptions: parseQueryOptions(req),
		ID:           ifaceID,
	}

	var out structs.SingleInterfaceResponse
	if err := h.rpcClient.Call("Interface.GetInterface", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return out.Interface, nil
}

func (h *InterfaceHandler) handleList(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	args := &structs.InterfaceListRequest{
		QueryOptions: parseQueryOptions(req),
		NodeID:       req.URL.Query().Get("node"),
		NetworkID:    req.URL.Query().Get("network"),
	}

	var out structs.InterfaceListResponse
	if err := h.rpcClient.Call("Interface.ListInterfaces", &args, &out); err != nil {
		return nil, parseError(err)
	}

	if out.Items == nil {
		out.Items = make([]*structs.InterfaceListStub, 0)
	}

	return out.Items, nil
}

func (h *InterfaceHandler) handlePost(rw http.ResponseWriter, req *http.Request, ifaceID string) (interface{}, error) {

	var iface structs.Interface
	err := parseBody(req.Body, &iface)
	if err != nil {
		return nil, NewCodedError(500, ErrInternal, err)
	}

	// Make sure the interface ID matches
	if iface.ID != ifaceID {
		return nil, NewCodedError(400, "Interface ID does not match request path")
	}

	args := &structs.InterfaceUpsertRequest{
		Interface:    &iface,
		WriteRequest: parseWriteRequestOptions(req),
	}

	var out structs.GenericResponse
	if err := h.rpcClient.Call("Interface.UpsertInterface", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return nil, nil
}

func (h *InterfaceHandler) handleDelete(rw http.ResponseWriter, req *http.Request, ifaceID string) (interface{}, error) {

	args := structs.InterfaceDeleteRequest{
		WriteRequest: parseWriteRequestOptions(req),
		IDs:          []string{ifaceID},
	}

	var out structs.GenericResponse
	if err := h.rpcClient.Call("Interface.DeleteInterface", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return nil, nil
}
