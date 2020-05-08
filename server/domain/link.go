package domain

import "time"

// Link :
type Link struct {
	ID                  *string    `json:"id"`
	NetworkID           *string    `json:"network"`
	FromHostID          *string    `json:"fromHost"`
	ToHostID            *string    `json:"toHost"`
	AllowedIPs          []string   `json:"allowedIps"`
	PersistentKeepalive *int       `json:"persistentKeepalive"`
	CreatedAt           *time.Time `json:"createdAt"`
	UpdatedAt           *time.Time `json:"updatedAt"`
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
