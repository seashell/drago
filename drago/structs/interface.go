package structs

import "time"

type Interface struct {
	ID         string
	NodeID     string
	NetworkID  string
	Name       string
	Address    string
	ListenPort int
	Peers      []*Peer
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
	if in.ListenPort != 0 {
		result.ListenPort = in.ListenPort
	}

	return &result
}

// Stub :
func (i *Interface) Stub() *InterfaceListStub {
	return &InterfaceListStub{
		ID:         i.ID,
		Name:       i.Name,
		Address:    i.Address,
		ListenPort: i.ListenPort,
		NodeID:     i.NodeID,
		NetworkID:  i.NetworkID,
		CreatedAt:  i.CreatedAt,
		UpdatedAt:  i.UpdatedAt,
	}
}

// InterfaceListStub :
type InterfaceListStub struct {
	ID         string
	NodeID     string
	NetworkID  string
	Name       string
	Address    string
	ListenPort int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// InterfaceSpecificRequest :
type InterfaceSpecificRequest struct {
	ID string

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
	IDs []string

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
	ID                  string
	PublicKey           string
	Address             string
	Port                int
	AllowedIPs          []string
	PersistentKeepalive int
}

type Link struct{}
