package domain

import "time"

// Interface :
type Interface struct {
	ID         *string    `json:"id"`
	Name       *string    `json:"name,omitempty"`
	HostID     *string    `json:"host,omitempty"`
	NetworkID  *string    `json:"network,omitempty"`
	IPAddress  *string    `json:"ipAddress,omitempty"`
	ListenPort *string    `json:"listenPort,omitempty"`
	Table      *string    `json:"table,omitempty"`
	DNS        *string    `json:"dns,omitempty"`
	MTU        *string    `json:"mtu,omitempty"`
	PreUp      *string    `json:"preUp,omitempty"`
	PostUp     *string    `json:"postUp,omitempty"`
	PreDown    *string    `json:"preDown,omitempty"`
	PostDown   *string    `json:"postDown,omitempty"`
	PublicKey  *string    `json:"publicKey,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

// InterfaceRepository :
type InterfaceRepository interface {
	GetByID(id string) (*Interface, error)
	Create(i *Interface) (*string, error)
	Update(i *Interface) (*string, error)
	DeleteByID(id string) (*string, error)
	FindAll(pageInfo PageInfo) ([]*Interface, *Page, error)
	FindAllByNetworkID(id string, pageInfo PageInfo) ([]*Interface, *Page, error)
	FindAllByHostID(id string, pageInfo PageInfo) ([]*Interface, *Page, error)
}
