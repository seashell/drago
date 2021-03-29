package http

import (
	"net/http"

	"github.com/seashell/drago/agent/conn"
	structs "github.com/seashell/drago/drago/structs"
)

// NetworkHandler :
type NetworkHandler struct {
	rpcConn conn.RPCConnection
}

// NewNetworkHandler :
func NewNetworkHandler(conn conn.RPCConnection) *NetworkHandler {
	return &NetworkHandler{
		rpcConn: conn,
	}
}

// Handle :
func (h *NetworkHandler) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)
	if len(params) > 1 {
		return nil, NewCodedError(404, ErrNotFound)
	}

	networkID := params[0]

	switch req.Method {
	case "GET":
		return h.handleGet(rw, req, networkID)
	case "PUT", "POST":
		return h.handlePost(rw, req, networkID)
	case "DELETE":
		return h.handleDelete(rw, req, networkID)
	default:
		return nil, NewCodedError(405, ErrMethodNotAllowed)
	}
}

func (h *NetworkHandler) handleGet(rw http.ResponseWriter, req *http.Request, networkID string) (interface{}, error) {

	if networkID == "" {
		return h.handleList(rw, req)
	}

	args := structs.NetworkSpecificRequest{
		QueryOptions: parseQueryOptions(req),
		NetworkID:    networkID,
	}

	var out structs.SingleNetworkResponse
	if err := h.rpcConn.Call("Network.GetNetwork", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return out.Network, nil
}

func (h *NetworkHandler) handleList(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	args := &structs.NetworkListRequest{
		QueryOptions: parseQueryOptions(req),
	}

	var out structs.NetworkListResponse
	if err := h.rpcConn.Call("Network.ListNetworks", &args, &out); err != nil {
		return nil, parseError(err)
	}

	if out.Items == nil {
		out.Items = make([]*structs.NetworkListStub, 0)
	}

	return out.Items, nil
}

func (h *NetworkHandler) handlePost(rw http.ResponseWriter, req *http.Request, networkID string) (interface{}, error) {

	var network structs.Network
	err := parseBody(req.Body, &network)
	if err != nil {
		return nil, NewCodedError(500, ErrInternal, err)
	}

	// Make sure the network ID matches
	if network.ID != networkID {
		return nil, NewCodedError(400, "Network ID does not match request path")
	}

	args := &structs.NetworkUpsertRequest{
		Network:      &network,
		WriteRequest: parseWriteRequestOptions(req),
	}

	var out structs.GenericResponse
	if err := h.rpcConn.Call("Network.UpsertNetwork", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return nil, nil
}

func (h *NetworkHandler) handleDelete(rw http.ResponseWriter, req *http.Request, networkID string) (interface{}, error) {

	args := structs.NetworkDeleteRequest{
		WriteRequest: parseWriteRequestOptions(req),
		NetworkIDs:   []string{networkID},
	}

	var out structs.GenericResponse
	if err := h.rpcConn.Call("Network.DeleteNetwork", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return nil, nil
}
