package domain

import (
	"time"
)

// Repository :
type Repository interface {
	Get() (*Host, error)
	Save(*Host) error
}

// Host :
type Host struct {
	Interfaces []*Interface
	Peers      []*Peer
}

// Interface :
type Interface struct {
	Name       string
	PublicKey  string
	Address    string
	ListenPort int
	Table      string
	DNS        string
	MTU        string
	PreUp      string
	PostUp     string
	PreDown    string
	PostDown   string
	Peers      []*Peer
}

// Peer :
type Peer struct {
	Interface           string
	Address             string
	Port                int
	PublicKey           string
	AllowedIPs          []string
	PersistentKeepalive int
	LatencyMs           uint64
	LastHandshake       time.Time
}
