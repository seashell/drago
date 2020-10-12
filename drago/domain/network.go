package domain

import (
	"context"
	"errors"
	"net"
	"time"
)

// Network ...
type Network struct {
	ID           string
	Name         string
	AddressRange string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NetworkRepository : Network repository interface
type NetworkRepository interface {
	GetByID(ctx context.Context, id string) (*Network, error)
	Create(ctx context.Context, n *Network) (*string, error)
	Update(ctx context.Context, n *Network) (*string, error)
	DeleteByID(ctx context.Context, id string) (*string, error)
	FindAll(ctx context.Context) ([]*Network, error)
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
