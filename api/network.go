package api

import (
	"path"

	"github.com/seashell/drago/drago/structs"
)

const (
	networksPath = "/api/networks"
)

// Networks is a handle to the nodes API
type Networks struct {
	client *Client
}

// Networks returns a handle on the networks endpoints.
func (c *Client) Networks() *Networks {
	return &Networks{client: c}
}

// Get :
func (n *Networks) Get(id string, opts *structs.QueryOptions) (*structs.Network, error) {

	var network *structs.Network
	err := n.client.getResource(networksPath, id, &network)
	if err != nil {
		return nil, err
	}

	return network, nil
}

// List :
func (n *Networks) List(opts *structs.QueryOptions) ([]*structs.NetworkListStub, error) {

	var items []*structs.NetworkListStub
	err := n.client.listResources(path.Join(networksPath, "/"), nil, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
