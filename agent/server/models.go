package server

import "time"

type Host struct {
	ID               int
	Name             string
	Address          string
	AdvertiseAddress string
	ListenPort       string
	PublicKey        string
	Table            string
	DNS              string
	Mtu              string
	PreUp            string
	PostUp           string
	PreDown          string
	PostDown         string
	Links            []*Link
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

type Link struct {
	ID                  int
	FromID              int
	ToID                int
	From                *Host
	To                  *Host
	AllowedIPs          string
	PersistentKeepalive int
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

type Repository interface {
	CreateHost(n *Host) (*Host, error)
	UpdateHost(id int, n *Host) (*Host, error)
	DeleteHost(id int) error
	GetAllHosts() ([]*Host, error)
	GetHost(id int) (*Host, error)
	CreateLink(l *Link) (*Link, error)
	UpdateLink(id int, n *Link) (*Link, error)
	DeleteLink(id int) error
	GetAllLinks() ([]*Link, error)
}
