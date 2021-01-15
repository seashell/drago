package structs

import (
	"errors"
	"net"
	"time"
)

// Network :
type Network struct {
	ID           string
	Name         string
	AddressRange string
	ModifyIndex  uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (n *Network) Validate() error {
	return nil
}

// CheckAddressInRange : Check whether an IP address in CIDR notation
// is within the allowed range of the network.
func (n *Network) CheckAddressInRange(ip string) error {
	_, subnet, _ := net.ParseCIDR(n.AddressRange)
	addr, _, _ := net.ParseCIDR(ip)
	if subnet.Contains(addr) {
		return nil
	}
	return errors.New("ip address not within network's allowed range")
}

// Merge :
func (n *Network) Merge(in *Network) *Network {

	result := *n

	if in.ID != "" {
		result.ID = in.ID
	}
	if in.Name != "" {
		result.Name = in.Name
	}
	if in.ModifyIndex != 0 {
		result.ModifyIndex = in.ModifyIndex
	}
	if in.AddressRange != "" {
		result.AddressRange = in.AddressRange
	}

	return &result
}

// Stub :
func (n *Network) Stub() *NetworkListStub {
	return &NetworkListStub{
		ID:           n.ID,
		Name:         n.Name,
		AddressRange: n.AddressRange,
		ModifyIndex:  n.ModifyIndex,
		CreatedAt:    n.CreatedAt,
		UpdatedAt:    n.UpdatedAt,
	}
}

// NetworkListStub :
type NetworkListStub struct {
	ID           string
	Name         string
	AddressRange string
	ModifyIndex  uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NetworkSpecificRequest :
type NetworkSpecificRequest struct {
	NetworkID string

	QueryOptions
}

// SingleNetworkResponse :
type SingleNetworkResponse struct {
	Network *Network

	Response
}

// NetworkUpsertRequest :
type NetworkUpsertRequest struct {
	Network *Network

	WriteRequest
}

// NetworkDeleteRequest :
type NetworkDeleteRequest struct {
	NetworkIDs []string

	WriteRequest
}

// NetworkListRequest :
type NetworkListRequest struct {
	QueryOptions
}

// NetworkListResponse :
type NetworkListResponse struct {
	Items []*NetworkListStub

	Response
}
