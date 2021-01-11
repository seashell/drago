package http

import (
	"net/http"

	structs "github.com/seashell/drago/drago/structs"
	rpc "github.com/seashell/drago/pkg/rpc"
)

// ACLPolicyHandler :
type ACLPolicyHandler struct {
	rpcClient *rpc.Client
}

// NewACLPolicyHandler :
func NewACLPolicyHandler(rpcClient *rpc.Client) *ACLPolicyHandler {
	return &ACLPolicyHandler{
		rpcClient: rpcClient,
	}
}

// Handle :
func (h *ACLPolicyHandler) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)
	if len(params) > 1 {
		return nil, NewCodedError(404, ErrNotFound)
	}

	policyName := params[0]

	switch req.Method {
	case "GET":
		return h.handleGet(rw, req, policyName)
	case "POST":
		return h.handlePost(rw, req, policyName)
	case "DELETE":
		return h.handleDelete(rw, req, policyName)
	default:
		return nil, NewCodedError(405, ErrMethodNotAllowed)
	}
}

func (h *ACLPolicyHandler) handleGet(rw http.ResponseWriter, req *http.Request, policyName string) (interface{}, error) {

	if policyName == "" {
		return h.handleList(rw, req)
	}

	args := structs.ACLPolicySpecificRequest{
		QueryOptions: parseQueryOptions(req),
		Name:         policyName,
	}

	var out structs.SingleACLPolicyResponse
	if err := h.rpcClient.Call("ACL.GetPolicy", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return out.ACLPolicy, nil
}

func (h *ACLPolicyHandler) handleList(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	args := &structs.ACLPolicyListRequest{
		QueryOptions: parseQueryOptions(req),
	}

	var out structs.ACLPolicyListResponse
	if err := h.rpcClient.Call("ACL.ListPolicies", &args, &out); err != nil {
		return nil, parseError(err)
	}

	if out.Items == nil {
		out.Items = make([]*structs.ACLPolicyListStub, 0)
	}

	return out.Items, nil
}

func (h *ACLPolicyHandler) handlePost(rw http.ResponseWriter, req *http.Request, policyName string) (interface{}, error) {

	var policy structs.ACLPolicy
	err := parseBody(req.Body, &policy)
	if err != nil {
		return nil, NewCodedError(500, ErrInternal, err)
	}

	// Make sure the policy name matches
	if policy.Name != policyName {
		return nil, NewCodedError(400, "ACL policy name does not match request path")
	}

	args := &structs.ACLPolicyUpsertRequest{
		ACLPolicy:    &policy,
		WriteRequest: parseWriteRequestOptions(req),
	}

	var out structs.GenericResponse
	if err := h.rpcClient.Call("ACL.UpsertACLPolicy", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return nil, nil
}

func (h *ACLPolicyHandler) handleDelete(rw http.ResponseWriter, req *http.Request, policyName string) (interface{}, error) {

	args := structs.ACLPolicyDeleteRequest{
		WriteRequest: parseWriteRequestOptions(req),
		Names:        []string{policyName},
	}

	var out structs.GenericResponse
	if err := h.rpcClient.Call("ACL.DeleteACLPolicy", &args, &out); err != nil {
		return nil, parseError(err)
	}

	return nil, nil
}
