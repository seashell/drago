package domain

import "time"

// Link :
type Link struct {
	ID                  *string    `json:"id"`
	NetworkID           *string    `json:"network,omitempty"`
	FromHostID          *string    `json:"fromHost,omitempty"`
	ToHostID            *string    `json:"toHost,omitempty"`
	AllowedIPs          []string   `json:"allowedIps,omitempty"`
	PersistentKeepalive *int       `json:"persistentKeepalive,omitempty"`
	CreatedAt           *time.Time `json:"createdAt,omitempty"`
	UpdatedAt           *time.Time `json:"updatedAt,omitempty"`
}

// LinkRepository :
type LinkRepository interface {
	GetByID(id string) (*Link, error)
	Create(l *Link) (*string, error)
	Update(l *Link) (*string, error)
	DeleteByID(id string) (*string, error)
	FindAllByNetworkID(id string, pageInfo PageInfo) ([]*Link, *Page, error)
	FindAllByHostID(id string, pageInfo PageInfo) ([]*Link, *Page, error)
}
