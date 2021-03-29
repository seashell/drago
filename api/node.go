package api

import (
	"path"

	"github.com/seashell/drago/drago/structs"
)

const (
	nodesPath = "/api/nodes"
)

// Nodes is a handle to the nodes API
type Nodes struct {
	client *Client
}

// Nodes returns a handle on the nodes endpoints.
func (c *Client) Nodes() *Nodes {
	return &Nodes{client: c}
}

// Get :
func (t *Nodes) Get(id string) (*structs.Node, error) {

	var node *structs.Node
	err := t.client.getResource(nodesPath, id, &node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// List :
func (t *Nodes) List(filters map[string][]string) ([]*structs.NodeListStub, error) {

	var items []*structs.NodeListStub
	err := t.client.listResources(path.Join(nodesPath, "/"), filters, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
