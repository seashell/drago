package http

import (
	"net/http"

	"github.com/seashell/drago/agent/conn"
	structs "github.com/seashell/drago/drago/structs"
)

// ConnectionHandler :
type ConnectionHandler struct {
	rpcConn conn.RPCConnection
}

// NewConnectionHandler :
func NewConnectionHandler(conn conn.RPCConnection) *ConnectionHandler {
	return &ConnectionHandler{
		rpcConn: conn,
	}
}

// Handle :
func (h *ConnectionHandler) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)
	if len(params) > 1 {
		return nil, NewCodedError(404, ErrNotFound)
	}

	connID := params[0]

	switch req.Method {
	case "GET":
		return h.handleGet(rw, req, connID)
	case "PUT", "POST":
		return h.handlePost(rw, req, connID)
	case "DELETE":
		return h.handleDelete(rw, req, connID)
	default:
		return nil, NewCodedError(405, ErrMethodNotAllowed)
	}
}

func (h *ConnectionHandler) handleGet(rw http.ResponseWriter, req *http.Request, connID string) (interface{}, error) {

	if connID == "" {
		return h.handleList(rw, req)
	}

	args := structs.ConnectionSpecificRequest{
		QueryOptions: parseQueryOptions(req),
		ConnectionID: connID,
	}

	var out structs.SingleConnectionResponse
	if err := h.rpcConn.Call("Connection.GetConnection", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return out.Connection, nil
}

func (h *ConnectionHandler) handleList(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	args := &structs.ConnectionListRequest{
		QueryOptions: parseQueryOptions(req),
		InterfaceID:  req.URL.Query().Get("interface"),
		NodeID:       req.URL.Query().Get("node"),
		NetworkID:    req.URL.Query().Get("network"),
	}

	var out structs.ConnectionListResponse
	if err := h.rpcConn.Call("Connection.ListConnections", &args, &out); err != nil {
		return nil, parseError(err)
	}

	if out.Items == nil {
		out.Items = make([]*structs.ConnectionListStub, 0)
	}

	return out.Items, nil
}

func (h *ConnectionHandler) handlePost(rw http.ResponseWriter, req *http.Request, ifaceID string) (interface{}, error) {

	var conn structs.Connection
	err := parseBody(req.Body, &conn)
	if err != nil {
		return nil, NewCodedError(500, ErrInternal, err)
	}

	// Make sure the interface ID matches
	if conn.ID != ifaceID {
		return nil, NewCodedError(400, "Connection ID does not match request path")
	}

	args := &structs.ConnectionUpsertRequest{
		Connection:   &conn,
		WriteRequest: parseWriteRequestOptions(req),
	}

	var out structs.GenericResponse
	if err := h.rpcConn.Call("Connection.UpsertConnection", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return nil, nil
}

func (h *ConnectionHandler) handleDelete(rw http.ResponseWriter, req *http.Request, connID string) (interface{}, error) {

	args := structs.ConnectionDeleteRequest{
		WriteRequest:  parseWriteRequestOptions(req),
		ConnectionIDs: []string{connID},
	}

	var out structs.GenericResponse
	if err := h.rpcConn.Call("Connection.DeleteConnection", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return nil, nil
}
