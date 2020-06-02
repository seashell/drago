package domain

import (
	"errors"
	"net"
	"time"
)

// Network : Network entity
type Network struct {
	ID             *string    `json:"id"`
	Name           *string    `json:"name,omitempty"`
	IPAddressRange *string    `json:"ipAddressRange,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
}

// NetworkRepository : Network repository interface
type NetworkRepository interface {
	GetByID(string) (*Network, error)
	Create(n *Network) (*string, error)
	Update(n *Network) (*string, error)
	DeleteByID(string) (*string, error)
	FindAll(pageInfo PageInfo) ([]*Network, *Page, error)
}

// CheckAddressInRange : Check whether an IP address in CIDR notation
// is within the allowed range of the network.
func (n *Network) CheckAddressInRange(ip string) error {
	_, subnet, _ := net.ParseCIDR(*n.IPAddressRange)
	addr, _, _ := net.ParseCIDR(ip)
	if subnet.Contains(addr) {
		return nil
	}
	return errors.New("ip address not within network's allowed range")
}
