package domain

// Link :
type Link struct {
	BaseModel
	NetworkID           *string  `json:"network"`
	FromHostID          *string  `json:"fromHost"`
	ToHostID            *string  `json:"toHost"`
	AllowedIPs          []string `json:"allowedIps"`
	PersistentKeepalive *int     `json:"persistentKeepalive"`
}

// LinkRepository :
type LinkRepository interface {
	GetByID(id string) (*Link, error)
	Create(l *Link) (id *string, err error)
	Update(l *Link) (id *string, err error)
	DeleteByID(id string) error
	FindAllByNetworkID(id string, pageInfo PageInfo) ([]*Link, *Page, error)
	FindAllByHostID(id string, pageInfo PageInfo) ([]*Link, *Page, error)
}
