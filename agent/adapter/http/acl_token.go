package http

import (
	"net/http"

	conn "github.com/seashell/drago/agent/conn"
	structs "github.com/seashell/drago/drago/structs"
)

// ACLTokenHandler :
type ACLTokenHandler struct {
	rpcConn conn.RPCConnection
}

// NewACLTokenHandler :
func NewACLTokenHandler(conn conn.RPCConnection) *ACLTokenHandler {
	return &ACLTokenHandler{
		rpcConn: conn,
	}
}

// Handle :
func (h *ACLTokenHandler) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)
	if len(params) > 1 {
		return nil, NewCodedError(404, ErrNotFound)
	}

	tokenID := params[0]

	switch req.Method {
	case "GET":
		return h.handleGet(rw, req, tokenID)
	case "POST":
		return h.handlePost(rw, req, tokenID)
	case "DELETE":
		return h.handleDelete(rw, req, tokenID)
	default:
		return nil, NewCodedError(405, ErrMethodNotAllowed)
	}
}

func (h *ACLTokenHandler) handleGet(rw http.ResponseWriter, req *http.Request, tokenID string) (interface{}, error) {

	if tokenID == "" {
		return h.handleList(rw, req)
	}

	if tokenID == "self" {
		return h.handleGetSelf(rw, req)
	}

	args := structs.ACLTokenSpecificRequest{
		QueryOptions: parseQueryOptions(req),
		ACLTokenID:   tokenID,
	}

	var out structs.SingleACLTokenResponse
	if err := h.rpcConn.Call("ACL.GetToken", &args, &out); err != nil {
		return nil, parseError(err)
	}

	if out.ACLToken == nil {
		return nil, NewCodedError(404, "ACL token not found")
	}

	return out.ACLToken, nil
}

func (h *ACLTokenHandler) handleGetSelf(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	args := structs.ResolveACLTokenRequest{
		QueryOptions: parseQueryOptions(req),
		Secret:       parseAuthToken(req),
	}

	var out structs.SingleACLTokenResponse
	if err := h.rpcConn.Call("ACL.ResolveToken", &args, &out); err != nil {
		return nil, parseError(err)
	}

	if out.ACLToken == nil {
		return nil, NewCodedError(404, "ACL token not found")
	}

	return out.ACLToken, nil
}

func (h *ACLTokenHandler) handleList(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	args := &structs.ACLTokenListRequest{
		QueryOptions: parseQueryOptions(req),
	}

	var out structs.ACLTokenListResponse
	if err := h.rpcConn.Call("ACL.ListTokens", &args, &out); err != nil {
		return nil, parseError(err)
	}

	if out.Items == nil {
		out.Items = make([]*structs.ACLTokenListStub, 0)
	}

	return out.Items, nil
}

func (h *ACLTokenHandler) handlePost(rw http.ResponseWriter, req *http.Request, tokenID string) (interface{}, error) {

	var token structs.ACLToken
	err := parseBody(req.Body, &token)
	if err != nil {
		return nil, NewCodedError(400, ErrBadRequest, err)
	}

	// Make sure the token ID matches
	if token.ID != tokenID {
		return nil, NewCodedError(400, "ACL token ID does not match request path")
	}

	args := &structs.ACLTokenUpsertRequest{
		ACLToken:     &token,
		WriteRequest: parseWriteRequestOptions(req),
	}

	var out structs.ACLTokenUpsertResponse
	if err := h.rpcConn.Call("ACL.UpsertToken", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return out.ACLToken, nil
}

func (h *ACLTokenHandler) handleDelete(rw http.ResponseWriter, req *http.Request, tokenID string) (interface{}, error) {

	args := structs.ACLTokenDeleteRequest{
		WriteRequest: parseWriteRequestOptions(req),
		ACLTokenIDs:  []string{tokenID},
	}

	var out structs.GenericResponse
	if err := h.rpcConn.Call("ACL.DeleteToken", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return nil, nil
}
