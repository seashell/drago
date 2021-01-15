package structs

import (
	"fmt"
	"time"
)

type Interface struct {
	ID          string
	NodeID      string
	NetworkID   string
	Name        string
	Address     string
	ListenPort  uint16
	PublicKey   string
	ModifyIndex uint64
	Peers       []*Peer
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
	if in.Name != "" {
		result.Name = in.Name
	}
	if in.Address != "" {
		result.Address = in.Address
	}
	if in.PublicKey != "" {
		result.PublicKey = in.PublicKey
	}
	if in.ListenPort != 0 {
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

func (i *Interface) Validate() error {
	if i.NodeID == "" {
		return fmt.Errorf("missing node id")
	}
	if i.NetworkID == "" {
		return fmt.Errorf("missing network id")
	}
	return nil
}

// Stub :
func (i *Interface) Stub() *InterfaceListStub {
	return &InterfaceListStub{
		ID:           i.ID,
		Name:         i.Name,
		Address:      i.Address,
		ListenPort:   i.ListenPort,
		NodeID:       i.NodeID,
		NetworkID:    i.NetworkID,
		ModifyIndex:  i.ModifyIndex,
		HasPublicKey: i.PublicKey != "",
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
}

// InterfaceListStub :
type InterfaceListStub struct {
	ID           string
	NodeID       string
	NetworkID    string
	Name         string
	Address      string
	ListenPort   uint16
	ModifyIndex  uint64
	HasPublicKey bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

type Link struct {
	ID                  string
	FromInterfaceID     string
	ToInterfaceID       string
	AllowedIPs          []string
	PersistentKeepalive int
}

type Peer struct {
	PublicKey           string
	Address             string
	Port                int
	AllowedIPs          []string
	PersistentKeepalive int
}
