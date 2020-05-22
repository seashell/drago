package domain

import "time"

// Host :
type Host struct {
	ID               *string    `json:"id"`
	Name             *string    `json:"name,omitempty"`
	AdvertiseAddress *string    `json:"advertiseAddress,omitempty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
}

// HostSettings :
type HostSettings struct {
	ID         *string  `json:"id"`
	Interfaces []string `json:"interfaces,omitempty"`
	Peers      []string `json:"peers,omitempty"`
}

// HostRepository :
type HostRepository interface {
	GetByID(id string) (*Host, error)
	Create(h *Host) (*string, error)
	Update(h *Host) (*string, error)
	DeleteByID(id string) (*string, error)
	FindAll(pageInfo PageInfo) ([]*Host, *Page, error)
	FindAllByNetworkID(id string, pageInfo PageInfo) ([]*Host, *Page, error)
}
