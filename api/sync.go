package api

import (
	"context"

	structs "github.com/seashell/drago/drago/application/structs"
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

// SynchronizeSelf :
func (a *Agent) SynchronizeSelf(ctx context.Context, in *structs.HostSynchronizeInput) (*structs.HostSynchronizeOutput, error) {

	out := struct {
		*structs.HostSynchronizeOutput
	}{}

	err := a.client.updateResource(agentPath, "self", in, out)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
