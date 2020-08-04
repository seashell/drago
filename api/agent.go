package api

import (
	"context"
	"time"
)

const (
	agentPath = "/api/agent"
)

// HostSettings :
type HostSettings struct {
	Interfaces []*WgInterfaceSettings `json:"interfaces,omitempty"`
	Peers      []*WgPeerSettings      `json:"peers,omitempty"`
}

// HostState :
type HostState struct {
	Interfaces []*WgInterfaceState `json:"interfaces,omitempty"`
	Peers      []*WgPeerState      `json:"peers,omitempty"`
}

// WgInterfaceState :
type WgInterfaceState struct {
	Name      *string `json:"name,omitempty"`
	PublicKey *string `json:"publicKey,omitempty"`
}

// WgPeerState :
type WgPeerState struct {
	LatencyMs     uint64    `json:"latencyMs,omitempty"`
	LastHandshake time.Time `json:"lastHandshake,omitempty"`
}

// WgInterfaceSettings :
type WgInterfaceSettings struct {
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

// WgPeerSettings :
type WgPeerSettings struct {
	Interface           string   `json:"interface,omitempty"`
	Address             *string  `json:"address,omitempty"`
	Port                *string  `json:"port,omitempty"`
	PublicKey           *string  `json:"publicKey,omitempty"`
	AllowedIPs          []string `json:"allowedIps"`
	PersistentKeepalive *int     `json:"persistentKeepalive,omitempty"`
}

// Hosts is a handle to the agent API
type Agent struct {
	client *Client
}

// Hosts returns a handle on the agent endpoints.
func (c *Client) Agent() *Hosts {
	return &Hosts{client: c}
}

// GetSelfSettings :
func (h *Hosts) GetSelfSettings(ctx context.Context) (*HostSettings, error) {

	receiver := struct {
		*HostSettings
	}{}

	err := h.client.getResource(agentPath, "self", &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.HostSettings, nil
}

// UpdateSelfState :
func (h *Hosts) UpdateSelfState(ctx context.Context) error {
	err := h.client.updateResource(agentPath, "self", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// SynchronizeSelf :
func (h *Hosts) SynchronizeSelf(ctx context.Context, state *HostState) (*HostSettings, error) {
	return nil, nil
}
