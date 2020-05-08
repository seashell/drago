package domain

import "time"

// Host :
type Host struct {
	ID               *string    `json:"id"`
	NetworkID        *string    `json:"network"`
	LinkIDs          []string   `json:"links"`
	Name             *string    `json:"name"`
	IPAddress        *string    `json:"ipAddress"`
	AdvertiseAddress *string    `json:"advertiseAddress"`
	ListenPort       *string    `json:"listenPort"`
	PublicKey        *string    `json:"publicKey"`
	Table            *string    `json:"table"`
	DNS              *string    `json:"dns"`
	MTU              *string    `json:"mtu"`
	PreUp            *string    `json:"preUp"`
	PostUp           *string    `json:"postUp"`
	PreDown          *string    `json:"preDown"`
	PostDown         *string    `json:"postDown"`
	CreatedAt        *time.Time `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

// HostRepository :
type HostRepository interface {
	GetByID(id string) (*Host, error)
	Create(h *Host) (id *string, err error)
	Update(h *Host) (id *string, err error)
	DeleteByID(id string) error
	FindAllByNetworkID(id string, pageInfo PageInfo) ([]*Host, *Page, error)
}
