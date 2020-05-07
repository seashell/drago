package domain

// Network : Network entity
type Network struct {
	BaseModel
	Name           *string `json:"name"`
	IPAddressRange *string `json:"ipAddressRange"`
}

// NetworkRepository : Network repository interface
type NetworkRepository interface {
	GetByID(id string) (*Network, error)
	Create(n *Network) (id *string, err error)
	Update(n *Network) (id *string, err error)
	DeleteByID(id string) error
	FindAll(pageInfo PageInfo) ([]*Network, *Page, error)
}
