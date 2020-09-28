package domain

import "time"

// Interface :
type Interface struct {
	ID         *string
	Name       *string
	HostID     *string
	NetworkID  *string
	IPAddress  *string
	ListenPort *string
	Table      *string
	DNS        *string
	MTU        *string
	PreUp      *string
	PostUp     *string
	PreDown    *string
	PostDown   *string
	PublicKey  *string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
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

func (i *Interface) Interface(b *Interface) *Interface {

	result := *i

	return &result
}
