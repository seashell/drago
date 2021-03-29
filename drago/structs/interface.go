package structs

import (
	"time"
)

type Interface struct {
	ID          string
	NodeID      string
	NetworkID   string
	Name        *string
	Address     *string
	ListenPort  *int
	PublicKey   *string
	ModifyIndex uint64
	Peers       []*Peer
	Connections []string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Underlying struct for efficiently adding/removing connections.
	// Always use the lazyConnectionsMap() method for accessing it.
	connectionsMap map[string]struct{}
}

// Merge :
func (i *Interface) Merge(in *Interface) *Interface {

	result := *i

	if in.ID != "" {
		result.ID = in.ID
	}
	if in.NodeID != "" {
		result.NodeID = in.NodeID
	}
	if in.NetworkID != "" {
		result.NetworkID = in.NetworkID
	}
	if in.Name != nil {
		result.Name = in.Name
	}
	if in.Address != nil {
		result.Address = in.Address
	}
	if in.PublicKey != nil {
		result.PublicKey = in.PublicKey
	}
	if in.ListenPort != nil {
		result.ListenPort = in.ListenPort
	}
	if in.Peers != nil {
		result.Peers = in.Peers
	}
	if in.ModifyIndex != 0 {
		result.ModifyIndex = in.ModifyIndex
	}

	return &result
}

// Validate : validate interface fields
func (i *Interface) Validate() error {
	return nil
}

// If the interfaces's connectionsMap was already initialized, return it.
// Otherwise initialize and synchronize it with the interface connections slice.
func (i *Interface) lazyConnectionsMap() map[string]struct{} {

	if i.connectionsMap != nil {
		return i.connectionsMap
	}

	i.connectionsMap = map[string]struct{}{}
	for _, conn := range i.Connections {
		i.connectionsMap[conn] = struct{}{}
	}
	return i.connectionsMap
}

// UpsertConnection :
func (i *Interface) UpsertConnection(id string) {
	i.lazyConnectionsMap()[id] = struct{}{}
	tmp := i.Connections[:0]
	for k := range i.connectionsMap {
		tmp = append(tmp, k)
	}
	i.Connections = tmp
}

// RemoveConnection :
func (i *Interface) RemoveConnection(id string) {

	delete(i.lazyConnectionsMap(), id)
	tmp := i.Connections[:0]
	for k := range i.connectionsMap {
		tmp = append(tmp, k)
	}
	i.Connections = tmp
}

// Stub :
func (i *Interface) Stub() *InterfaceListStub {
	return &InterfaceListStub{
		ID:               i.ID,
		Name:             i.Name,
		Address:          i.Address,
		ListenPort:       i.ListenPort,
		NodeID:           i.NodeID,
		NetworkID:        i.NetworkID,
		ModifyIndex:      i.ModifyIndex,
		ConnectionsCount: len(i.Connections),
		PublicKey:        i.PublicKey,
		HasPublicKey:     i.PublicKey != nil,
		CreatedAt:        i.CreatedAt,
		UpdatedAt:        i.UpdatedAt,
	}
}

// InterfaceListStub :
type InterfaceListStub struct {
	ID               string
	NodeID           string
	NetworkID        string
	Name             *string
	Address          *string
	ListenPort       *int
	ConnectionsCount int
	ModifyIndex      uint64
	PublicKey        *string
	HasPublicKey     bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// InterfaceSpecificRequest :
type InterfaceSpecificRequest struct {
	InterfaceID string

	QueryOptions
}

// SingleInterfaceResponse :
type SingleInterfaceResponse struct {
	Interface *Interface

	Response
}

// InterfaceUpsertRequest :
type InterfaceUpsertRequest struct {
	Interface *Interface

	WriteRequest
}

// InterfaceDeleteRequest :
type InterfaceDeleteRequest struct {
	InterfaceIDs []string

	WriteRequest
}

// InterfaceListRequest :
type InterfaceListRequest struct {
	NodeID    string
	NetworkID string

	QueryOptions
}

// InterfaceListResponse :
type InterfaceListResponse struct {
	Items []*InterfaceListStub

	Response
}

type Peer struct {
	PublicKey           *string
	Address             *string
	Port                *int
	AllowedIPs          []string
	PersistentKeepalive *int
}
