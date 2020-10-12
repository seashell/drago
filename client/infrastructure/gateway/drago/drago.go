package drago

import (
	api "github.com/seashell/drago/api"
	application "github.com/seashell/drago/client/application"
)

// DragoGatewayAdapter implements the application.DragoGateway interface.
type DragoGatewayAdapter struct {
	client *api.Client
}

func NewDragoGatewayAdapter(addr, token string) (application.DragoGateway, error) {

	cli, err := api.NewClient(&api.Config{
		Address: addr,
		Token:   token,
	})
	if err != nil {
		return nil, err
	}

	return &DragoGatewayAdapter{
		client: cli,
	}, nil
}

// Agent returns a handler to the Agent API.
func (g *DragoGatewayAdapter) Agent() application.DragoAgentGateway {
	return g.client.Agent()
}

// Tokens returns a handler to the Tokens API.
func (g *DragoGatewayAdapter) Tokens() application.DragoTokenGateway {
	return g.client.Tokens()
}
