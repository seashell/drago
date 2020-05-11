package domain

import "time"

// Host :
type Host struct {
	ID               *string    `json:"id"`
	NetworkID        *string    `json:"network,omitempty"`
	LinkIDs          []string   `json:"links,omitempty"`
	Name             *string    `json:"name,omitempty"`
	IPAddress        *string    `json:"ipAddress,omitempty"`
	AdvertiseAddress *string    `json:"advertiseAddress,omitempty"`
	ListenPort       *string    `json:"listenPort,omitempty"`
	PublicKey        *string    `json:"publicKey,omitempty"`
	Table            *string    `json:"table,omitempty"`
	DNS              *string    `json:"dns,omitempty"`
	MTU              *string    `json:"mtu,omitempty"`
	PreUp            *string    `json:"preUp,omitempty"`
	PostUp           *string    `json:"postUp,omitempty"`
	PreDown          *string    `json:"preDown,omitempty"`
	PostDown         *string    `json:"postDown,omitempty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
}

// HostRepository :
type HostRepository interface {
	GetByID(id string) (*Host, error)
	Create(h *Host) (*string, error)
	Update(h *Host) (*string, error)
	DeleteByID(id string) (*string, error)
	FindAllByNetworkID(id string, pageInfo PageInfo) ([]*Host, *Page, error)
}
