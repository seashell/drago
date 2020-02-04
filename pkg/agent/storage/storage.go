package storage

import "time"

type Peer struct {
	Endpoint            string
	AllowedIPs          string
	PublicKey           string
	PersistentKeepalive int
}

type Node struct {
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

type Edge struct {
	ID int

	Source int
	Target int

	AllowedIPs          string
	PersistentKeepalive int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Store interface {
	SelectAllNodes() (map[int]*Node, error)
	SelectNode(int) (*Node, error)
	DeleteNode(int) error
	InsertNode(*Node) (*Node, error)
	UpdateNode(int, *Node) (*Node, error)

	SelectAllPeersForNode(id int) ([]*Peer, error)

	SelectAllEdges() (map[int]*Edge, error)
	SelectEdge(int) (*Edge, error)
	DeleteEdge(int) error
	InsertEdge(*Edge) (*Edge, error)
	UpdateEdge(int, *Edge) (*Edge, error)
}
