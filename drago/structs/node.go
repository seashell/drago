package structs

import (
	"fmt"
	"time"
)

const (
	NodeStatusInit  = "initializing"
	NodeStatusReady = "ready"
	NodeStatusDown  = "down"
)

// Node :
type Node struct {
	ID               string
	SecretID         string
	Name             string
	AdvertiseAddress string
	Status           string
	Interfaces       []string
	Connections      []string
	ModifyIndex      uint64
	Meta             map[string]string
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Underlying struct for efficiently adding/removing interfaces and connections.
	// Always use the lazyInterfacesMap() and lazyConnectionsMap() methods for accessing them.
	interfacesMap  map[string]struct{}
	connectionsMap map[string]struct{}
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
	if in.AdvertiseAddress != "" {
		result.AdvertiseAddress = in.AdvertiseAddress
	}
	if in.Interfaces != nil {
		result.Interfaces = in.Interfaces
	}
	if in.Status != "" {
		result.Status = in.Status
	}
	if in.ModifyIndex != 0 {
		result.ModifyIndex = in.ModifyIndex
	}

	return &result
}

// If the node's interfacesMap was already initialized, return it.
// Otherwise initialize and synchronize it with the node interfaces slice.
func (n *Node) lazyInterfacesMap() map[string]struct{} {

	if n.interfacesMap != nil {
		return n.interfacesMap
	}

	n.interfacesMap = map[string]struct{}{}
	for _, iface := range n.Interfaces {
		n.interfacesMap[iface] = struct{}{}
	}
	return n.interfacesMap
}

// UpsertInterface :
func (n *Node) UpsertInterface(id string) {
	n.lazyInterfacesMap()[id] = struct{}{}
	tmp := n.Interfaces[:0]
	for k := range n.interfacesMap {
		tmp = append(tmp, k)
	}
	n.Interfaces = tmp
}

// RemoveInterface :
func (n *Node) RemoveInterface(id string) {

	delete(n.lazyInterfacesMap(), id)
	tmp := n.Interfaces[:0]
	for k := range n.interfacesMap {
		tmp = append(tmp, k)
	}
	n.Interfaces = tmp
}

// If the node's connectionsMap was already initialized, return it.
// Otherwise initialize and synchronize it with the node connections slice.
func (n *Node) lazyConnectionsMap() map[string]struct{} {

	if n.connectionsMap != nil {
		return n.connectionsMap
	}

	n.connectionsMap = map[string]struct{}{}
	for _, conn := range n.Connections {
		n.connectionsMap[conn] = struct{}{}
	}
	return n.connectionsMap
}

// UpsertConnection :
func (n *Node) UpsertConnection(id string) {
	n.lazyConnectionsMap()[id] = struct{}{}
	tmp := n.Connections[:0]
	for k := range n.connectionsMap {
		tmp = append(tmp, k)
	}
	n.Connections = tmp
}

// RemoveConnection :
func (n *Node) RemoveConnection(id string) {

	delete(n.lazyConnectionsMap(), id)
	tmp := n.Connections[:0]
	for k := range n.connectionsMap {
		tmp = append(tmp, k)
	}
	n.Connections = tmp
}

// Stub :
func (n *Node) Stub() *NodeListStub {
	return &NodeListStub{
		ID:               n.ID,
		Name:             n.Name,
		AdvertiseAddress: n.AdvertiseAddress,
		Status:           n.Status,
		InterfacesCount:  len(n.Interfaces),
		ConnectionsCount: len(n.Connections),
		ModifyIndex:      n.ModifyIndex,
		Meta:             n.Meta,
		CreatedAt:        n.CreatedAt,
		UpdatedAt:        n.UpdatedAt,
	}
}

// NodeListStub :
type NodeListStub struct {
	ID               string
	Name             string
	AdvertiseAddress string
	Status           string
	InterfacesCount  int
	ConnectionsCount int
	ModifyIndex      uint64
	Meta             map[string]string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NodeSpecificRequest :
type NodeSpecificRequest struct {
	NodeID   string
	SecretID string

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
	NodeID           string
	Status           string
	AdvertiseAddress string
	Meta             map[string]string
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

// NodeInterfaceUpdateRequest :
type NodeInterfaceUpdateRequest struct {
	NodeID     string
	Interfaces []*Interface

	WriteRequest
}

// NodeJoinNetworkRequest :
type NodeJoinNetworkRequest struct {
	NodeID    string
	NetworkID string

	WriteRequest
}

// NodeLeaveNetworkRequest :
type NodeLeaveNetworkRequest struct {
	NodeID    string
	NetworkID string

	WriteRequest
}
