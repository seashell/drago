package structs

import (
	"errors"
	"fmt"
	"net"
	"time"
)

// Network :
type Network struct {
	ID           string
	Name         string
	AddressRange string
	Interfaces   []string
	Connections  []string
	ModifyIndex  uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Underlying structs for efficiently adding/removing interfaces and connections.
	// Always use the lazyInterfacesMap() and lazyConnectionsMap() methods for accessing them.
	interfacesMap  map[string]struct{}
	connectionsMap map[string]struct{}
}

// Validate :
func (n *Network) Validate() error {
	if n.Name == "" {
		return fmt.Errorf("Name is empty")
	}
	if n.AddressRange == "" {
		return fmt.Errorf("Address range is empty")
	}
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

// If the networks's interfacesMap was already initialized, return it.
// Otherwise initialize and synchronize it with the network interfaces slice.
func (n *Network) lazyInterfacesMap() map[string]struct{} {

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
func (n *Network) UpsertInterface(id string) {
	n.lazyInterfacesMap()[id] = struct{}{}
	tmp := n.Interfaces[:0]
	for k := range n.interfacesMap {
		tmp = append(tmp, k)
	}
	n.Interfaces = tmp
}

// RemoveInterface :
func (n *Network) RemoveInterface(id string) {

	delete(n.lazyInterfacesMap(), id)
	tmp := n.Interfaces[:0]
	for k := range n.interfacesMap {
		tmp = append(tmp, k)
	}
	n.Interfaces = tmp
}

// If the networks's connectionsMap was already initialized, return it.
// Otherwise initialize and synchronize it with the network connections slice.
func (n *Network) lazyConnectionsMap() map[string]struct{} {

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
func (n *Network) UpsertConnection(id string) {
	n.lazyConnectionsMap()[id] = struct{}{}
	tmp := n.Connections[:0]
	for k := range n.connectionsMap {
		tmp = append(tmp, k)
	}
	n.Connections = tmp
}

// RemoveConnection :
func (n *Network) RemoveConnection(id string) {

	delete(n.lazyConnectionsMap(), id)
	tmp := n.Connections[:0]
	for k := range n.connectionsMap {
		tmp = append(tmp, k)
	}
	n.Connections = tmp
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
	if in.Interfaces != nil {
		result.Interfaces = in.Interfaces
	}
	if in.Connections != nil {
		result.Connections = in.Connections
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
		ID:               n.ID,
		Name:             n.Name,
		AddressRange:     n.AddressRange,
		InterfacesCount:  len(n.Interfaces),
		ConnectionsCount: len(n.Connections),
		ModifyIndex:      n.ModifyIndex,
		CreatedAt:        n.CreatedAt,
		UpdatedAt:        n.UpdatedAt,
	}
}

// NetworkListStub :
type NetworkListStub struct {
	ID               string
	Name             string
	AddressRange     string
	InterfacesCount  int
	ConnectionsCount int
	ModifyIndex      uint64
	CreatedAt        time.Time
	UpdatedAt        time.Time
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
