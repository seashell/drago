package structs

import (
	"fmt"
	"time"
)

const (
	NodeStatusPreregistered = "preregistered"
	NodeStatusInit          = "initializing"
	NodeStatusReady         = "ready"
	NodeStatusDown          = "down"
)

// Node :
type Node struct {
	ID        string
	SecretID  string
	Name      string
	Address   string
	Status    string
	Meta      map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate validates a structs.Node object
func (n *Node) Validate() error {

	if !IsValidNodeStatus(n.Status) {
		return fmt.Errorf("invalid node status")
	}

	return nil
}

// IsValidNodeStatus returns true if the status passed as argument
// corresponds to a valid node status. Otherwise returns false.
func IsValidNodeStatus(s string) bool {

	valid := map[string]interface{}{
		NodeStatusInit:  nil,
		NodeStatusReady: nil,
		NodeStatusDown:  nil,
	}

	if _, ok := valid[s]; !ok {
		return false
	}

	return true
}

// Merge :
func (n *Node) Merge(in *Node) *Node {

	result := *n

	if in.ID != "" {
		result.ID = in.ID
	}
	if in.SecretID != "" {
		result.SecretID = in.SecretID
	}
	if in.Name != "" {
		result.Name = in.Name
	}
	if in.Status != "" {
		result.Status = in.Status
	}

	return &result
}

// Stub :
func (n *Node) Stub() *NodeListStub {
	return &NodeListStub{
		ID:        n.ID,
		Name:      n.Name,
		Status:    n.Status,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

// NodeListStub :
type NodeListStub struct {
	ID        string
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NodeSpecificRequest :
type NodeSpecificRequest struct {
	ID string

	QueryOptions
}

// SingleNodeResponse :
type SingleNodeResponse struct {
	Node *Node

	Response
}

// NodePreregisterRequest :
type NodePreregisterRequest struct {
	Node *Node

	WriteRequest
}

// NodePreregisterResponse :
type NodePreregisterResponse struct {
	Node *Node

	Response
}

// NodeRegisterRequest :
type NodeRegisterRequest struct {
	Node *Node

	WriteRequest
}

// Validate validates a structs.NodeRegisterRequest
func (r *NodeRegisterRequest) Validate() error {

	if r.Node == nil {
		return fmt.Errorf("missing node")
	}
	if r.Node.ID == "" {
		return fmt.Errorf("missing node ID")
	}
	if r.Node.Name == "" {
		return fmt.Errorf("missing node name")
	}
	if r.Node.SecretID == "" {
		return fmt.Errorf("missing node secret ID")
	}

	return nil
}

// NodeUpdateStatusRequest :
type NodeUpdateStatusRequest struct {
	ID     string
	Status string

	WriteRequest
}

// NodeUpdateResponse is used to update nodes
type NodeUpdateResponse struct {
	Servers []string

	Response
}

// NodeListRequest :
type NodeListRequest struct {
	QueryOptions
}

// NodeListResponse :
type NodeListResponse struct {
	Items []*NodeListStub

	Response
}

// NodeInterfacesResponse :
type NodeInterfacesResponse struct {
	Items []*Interface

	Response
}

// NodePeersResponse :
type NodePeersResponse struct {
	// Items []*Peer

	Response
}