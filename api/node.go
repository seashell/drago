package api

import (
	"context"
	"path"

	"github.com/seashell/drago/drago/structs"
)

const (
	nodesPath = "/api/node"
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
func (t *Nodes) Register(ctx context.Context, req *structs.NodeRegisterRequest) (*structs.NodeUpdateResponse, error) {

	var resp structs.NodeUpdateResponse
	err := t.client.createResource(path.Join(nodesPath, "register"), req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Update :
func (t *Nodes) UpdateStatus(ctx context.Context, req *structs.NodeUpdateStatusRequest) (*structs.NodeUpdateResponse, error) {

	var resp structs.NodeUpdateResponse
	err := t.client.createResource(path.Join(nodesPath, req.ID, "status"), req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetNodeInterfaces :
func (t *Nodes) GetNodeInterfaces(ctx context.Context, req *structs.NodeSpecificRequest) (*structs.NodeInterfacesResponse, error) {

	var resp structs.NodeInterfacesResponse
	err := t.client.listResources(path.Join(nodesPath, req.ID, "interfaces"), nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetNodePeers :
func (t *Nodes) GetNodePeers(ctx context.Context, req *structs.NodeSpecificRequest) (*structs.NodePeersResponse, error) {

	var resp structs.NodePeersResponse
	err := t.client.listResources(path.Join(nodesPath, req.ID, "peers"), nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
