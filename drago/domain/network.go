package domain

import (
	"errors"
	"net"
	"time"
)

// Network : Network entity
type Network struct {
	ID             *string
	Name           *string
	IPAddressRange *string
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

// NetworkRepository : Network repository interface
type NetworkRepository interface {
	GetByID(string) (*Network, error)
	Create(n *Network) (*string, error)
	Update(n *Network) (*string, error)
	DeleteByID(string) (*string, error)
	FindAll(pageInfo PageInfo) ([]*Network, *Page, error)
}

func (n *Network) Merge(b *Network) *Network {

	if b == nil {
		return n
	}

	result := *n

	if b.Name != nil {
		result.Name = b.Name
	}
	if b.IPAddressRange != nil {
		result.IPAddressRange = b.IPAddressRange
	}

	return &result
}

// IsAddressInRange : Check whether an IP address in CIDR notation
// is within the allowed range of the network.
func (n *Network) IsAddressInRange(ip string) error {
	_, subnet, _ := net.ParseCIDR(*n.IPAddressRange)
	addr, _, _ := net.ParseCIDR(ip)
	if subnet.Contains(addr) {
		return nil
	}
	return errors.New("ip address not within network's allowed range")
}
