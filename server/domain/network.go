package domain

import "time"

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
