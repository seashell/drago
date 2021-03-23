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

// Register :
func (t *Nodes) Register(req *structs.NodeRegisterRequest) (*structs.NodeUpdateResponse, error) {

	var resp structs.NodeUpdateResponse
	err := t.client.createResource(path.Join(nodesPath, "register"), req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Get :
func (t *Nodes) Get(id string, opts *structs.QueryOptions) (*structs.Node, error) {

	var node *structs.Node
	err := t.client.getResource(nodesPath, id, &node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// List :
func (t *Nodes) List(opts *structs.QueryOptions) ([]*structs.NodeListStub, error) {

	var items []*structs.NodeListStub
	err := t.client.listResources(path.Join(nodesPath, "/"), nil, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Update :
func (t *Nodes) UpdateStatus(req *structs.NodeUpdateStatusRequest) (*structs.NodeUpdateResponse, error) {

	var resp structs.NodeUpdateResponse
	err := t.client.createResource(path.Join(nodesPath, req.NodeID, "status"), req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetNodeInterfaces :
func (t *Nodes) GetNodeInterfaces(req *structs.NodeSpecificRequest) (*structs.NodeInterfacesResponse, error) {

	var resp structs.NodeInterfacesResponse
	err := t.client.listResources(path.Join(nodesPath, req.NodeID, "interfaces"), nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
