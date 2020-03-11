package storage

import "time"

type Peer struct {
	Endpoint            string
	AllowedIPs          string
	PublicKey           string
	PersistentKeepalive int
}

type Host struct {
	ID   int
	Name string

	Address    string
	ListenPort string
	PrivateKey string
	Table      string
	DNS        string
	MTU        string
	PreUp      string
	PostUp     string
	PreDown    string
	PostDown   string

	// Drago control fields
	PublicKey     string
	AdvertiseAddr string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Link struct {
	ID int

	Source int
	Target int

	AllowedIPs          string
	PersistentKeepalive int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Store interface {
	SelectAllHosts() (map[int]*Host, error)
	SelectHost(int) (*Host, error)
	DeleteHost(int) error
	InsertHost(*Host) (*Host, error)
	UpdateHost(int, *Host) (*Host, error)

	SelectAllPeersForHost(id int) ([]*Peer, error)

	SelectAllLinks() (map[int]*Link, error)
	SelectLink(int) (*Link, error)
	DeleteLink(int) error
	InsertLink(*Link) (*Link, error)
	UpdateLink(int, *Link) (*Link, error)
}
