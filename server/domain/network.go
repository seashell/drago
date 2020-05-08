package domain

import "time"

// Network : Network entity
type Network struct {
	ID             *string    `json:"id"`
	Name           *string    `json:"name"`
	IPAddressRange *string    `json:"ipAddressRange"`
	CreatedAt      *time.Time `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

// NetworkRepository : Network repository interface
type NetworkRepository interface {
	GetByID(id string) (*Network, error)
	Create(n *Network) (id *string, err error)
	Update(n *Network) (id *string, err error)
	DeleteByID(id string) error
	FindAll(pageInfo PageInfo) ([]*Network, *Page, error)
}
