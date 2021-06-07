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
func (n *Connections) Get(id string) (*structs.Connection, error) {

	var conn *structs.Connection
	err := n.client.getResource(connectionsPath, id, &conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// List :
func (n *Connections) List() ([]*structs.ConnectionListStub, error) {

	var items []*structs.ConnectionListStub
	err := n.client.listResources(path.Join(connectionsPath, "/"), nil, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (n *Connections) Create(connection *structs.Connection) (*structs.Connection, error) {

	out := &structs.Connection{}

	err := n.client.createResource(networksPath, connection, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}


func (n *Connections) Delete(id string) error {

	err := n.client.deleteResource(id, connectionsPath, nil)
	if err != nil {
		return err
	}

	return nil
}
