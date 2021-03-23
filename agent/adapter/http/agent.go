package http

import (
	"net/http"

	conn "github.com/seashell/drago/agent/conn"
	structs "github.com/seashell/drago/drago/structs"
)

type AgentAdapter interface {
	Config() map[string]interface{}
	Stats() map[string]map[string]string
}

// AgentHandler provides an API for interacting with an Agent
// in runtime, getting stats and setting configs.
type AgentHandler struct {
	agent AgentAdapter
}

// NewAgentHandler :
func NewAgentHandler(conn conn.RPCConnection, ag AgentAdapter) *AgentHandler {
	return &AgentHandler{
		agent: ag,
	}
}

// Handle :
func (h *AgentHandler) Handle(rw http.ResponseWriter, req *http.Request) (interface{}, error) {

	params := parsePathParams(req)
	if len(params) > 1 {
		return nil, NewCodedError(404, ErrNotFound)
	}

	id := params[0]

	switch req.Method {
	case "GET":
		return h.handleGet(rw, req, id)
	default:
		return nil, NewCodedError(405, ErrMethodNotAllowed)
	}
}

func (h *AgentHandler) handleGet(rw http.ResponseWriter, req *http.Request, id string) (interface{}, error) {

	if id != "self" {
		return nil, NewCodedError(404, ErrNotFound)
	}

	self := &structs.Agent{
		Config: map[string]interface{}{},
		Stats:  h.agent.Stats(),
	}

	return self, nil
}
