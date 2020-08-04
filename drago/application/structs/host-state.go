package structs

import "time"

// HostStateOutput :
type SynchronizeHostOutput struct {
	Interfaces []*WgInterfaceOutput `json:"interfaces"`
	Peers      []*WgPeerOutput      `json:"peers"`
}

// WgInterfaceOutput :
type WgInterfaceOutput struct {
	Name       *string `json:"name,omitempty"`
	Address    *string `json:"address,omitempty"`
	ListenPort *string `json:"listenPort,omitempty"`
	Table      *string `json:"table,omitempty"`
	DNS        *string `json:"dns,omitempty"`
	MTU        *string `json:"mtu,omitempty"`
	PreUp      *string `json:"preUp,omitempty"`
	PostUp     *string `json:"postUp,omitempty"`
	PreDown    *string `json:"preDown,omitempty"`
	PostDown   *string `json:"postDown,omitempty"`
}

// WgPeerOutput :
type WgPeerOutput struct {
	Interface           string    `json:"interface,omitempty"`
	Address             *string   `json:"address,omitempty"`
	Port                *string   `json:"port,omitempty"`
	PublicKey           *string   `json:"publicKey,omitempty"`
	AllowedIPs          []string  `json:"allowedIps"`
	PersistentKeepalive *int      `json:"persistentKeepalive,omitempty"`
	LatencyMs           uint64    `json:"latencyMs,omitempty"`
	LastHandshake       time.Time `json:"lastHandshake,omitempty"`
}

// HostStateInput :
type SynchronizeHostInput struct {
	ID         string              `validate:"required,uuid4"`
	Interfaces []*WgInterfaceInput `json:"interfaces"`
	Peers      []*WgPeerInput      `json:"peers"`
}

// WgInterfaceInput :
type WgInterfaceInput struct {
	Name      *string `json:"name" validate:"required"`
	PublicKey *string `json:"publicKey" validate:"required"`
}

// WgPeerInput :
type WgPeerInput struct {
	Interface     string    `json:"interface,omitempty"`
	PublicKey     *string   `json:"publicKey,omitempty"`
	LatencyMs     uint64    `json:"latencyMs,omitempty"`
	LastHandshake time.Time `json:"lastHandshake,omitempty"`
}
