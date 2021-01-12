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
	ID string

	SecretID string

	Name string

	Status string

	CreatedAt time.Time

	UpdatedAt time.Time

	Meta map[string]string
}

// Validate validates a structs.Node object
func (n *Node) Validate() error {

	valid := map[string]interface{}{
		NodeStatusInit:  nil,
		NodeStatusReady: nil,
		NodeStatusDown:  nil,
	}

	if _, ok := valid[n.Status]; !ok {
		return fmt.Errorf("invalid node status")
	}

	return nil
}

// NodeSpecificRequest :
type NodeSpecificRequest struct {
	QueryOptions
	ID string
}

// SingleNodeResponse :
type SingleNodeResponse struct {
	Response
	Node *Node
}

// NodeRegisterRequest :
type NodeRegisterRequest struct {
	WriteRequest
	Node *Node
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

// NodeUpdateRequest :
type NodeUpdateRequest struct {
	ID     string
	Status string

	Response
}

// NodeUpdateResponse is used to update nodes
type NodeUpdateResponse struct {
	HeartbeatTTL time.Duration
	InterfaceIDs []string
	PeerIDs      []string
	Response
}

// NodeInterfacesResponse :
type NodeInterfacesResponse struct {
	Interfaces []*Interface
	Response
}

// NodePeersResponse :
type NodePeersResponse struct {
	Peers []*Peer
	Response
}
