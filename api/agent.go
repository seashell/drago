package api

import (
	"path"

	"github.com/seashell/drago/drago/structs"
)

const (
	agentPath = "/api/agent"
)

// Agent is a handle to the agent API
type Agent struct {
	client *Client
}

// Agent returns a handle on the agent endpoints.
func (c *Client) Agent() *Agent {
	return &Agent{client: c}
}

// Self :
func (t *Agent) Self() (*structs.Agent, error) {

	var agent *structs.Agent
	err := t.client.getResource(path.Join(agentPath, "self"), "", &agent)
	if err != nil {
		return nil, err
	}

	return agent, nil
}
