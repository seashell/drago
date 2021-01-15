package http

import (
	"net/http"

	structs "github.com/seashell/drago/drago/structs"
	rpc "github.com/seashell/drago/pkg/rpc"
)

// NetworkHandler :
type NetworkHandler struct {
	rpcClient *rpc.Client
}

// NewNetworkHandler :
func NewNetworkHandler(rpcClient *rpc.Client) *NetworkHandler {
	return &NetworkHandler{
		rpcClient: rpcClient,
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
	if err := h.rpcClient.Call("Network.GetNetwork", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return out.Network, nil
}

func (h *NetworkHandler) handleList(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	args := &structs.NetworkListRequest{
		QueryOptions: parseQueryOptions(req),
	}

	var out structs.NetworkListResponse
	if err := h.rpcClient.Call("Network.ListNetworks", &args, &out); err != nil {
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
	if err := h.rpcClient.Call("Network.UpsertNetwork", &args, &out); err != nil {
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
	if err := h.rpcClient.Call("Network.DeleteNetwork", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return nil, nil
}
