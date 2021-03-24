package api

import (
	"path"

	"github.com/seashell/drago/drago/structs"
)

const (
	connectionsPath = "/api/connections"
)

// Connections is a handle to the connection API
type Connections struct {
	client *Client
}

// Connections returns a handle on the connections endpoints.
func (c *Client) Connections() *Connections {
	return &Connections{client: c}
}

// Get :
func (n *Connections) Get(id string, opts *structs.QueryOptions) (*structs.Connection, error) {

	var conn *structs.Connection
	err := n.client.getResource(connectionsPath, id, &conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// List :
func (n *Connections) List(opts *structs.QueryOptions) ([]*structs.ConnectionListStub, error) {

	var items []*structs.ConnectionListStub
	err := n.client.listResources(path.Join(connectionsPath, "/"), nil, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
