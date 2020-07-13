package api

// HostSettings :
type HostSettings struct {
	NetworkInterfaces []NetworkInterface `json:"interfaces"`
	WireguardPeers    []WireguardPeer    `json:"peers"`
}

// NetworkInterface :
type NetworkInterface struct {
	Name       *string `json:"name"`
	ListenPort *string `json:"listenPort"`
	Address    *string `json:"address"`
	Table      *string `json:"table"`
	DNS        *string `json:"dns"`
	MTU        *string `json:"mtu"`
	PreUp      *string `json:"preUp"`
	PostUp     *string `json:"postUp"`
	PreDown    *string `json:"preDown"`
	PostDown   *string `json:"postDown"`
}

// WireguardPeer :
type WireguardPeer struct {
	Interface           *string  `json:"interface"`
	Address             *string  `json:"address"`
	Port                *string  `json:"port"`
	PublicKey           *string  `json:"publicKey"`
	AllowedIps          []string `json:"allowedIps"`
	PersistentKeepalive *int     `json:"persistentKeepalive"`
}

// NetworkInterfaceState :
type NetworkInterfaceState struct {
	Name        string `json:"name"`
	WgPublicKey string `json:"publicKey"`
}

// HostState :
type HostState struct {
	NetworkInterfaces []NetworkInterfaceState `json:"interfaces"`
}

// Syncronization is used to query the sync-related endpoints.
type Synchronization struct {
	client *Client
}

// Syncronization returns a handle on the Syncronization endpoints.
func (c *Client) Synchronization() *Synchronization {
	return &Synchronization{client: c}
}

// Fetch self settings
func (s *Synchronization) GetSelfSettings() (*HostSettings, error) {
	var r HostSettings
	err := s.client.Get("/hosts/self/settings", &r, nil)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Fetch self settings
func (s *Synchronization) PostSelfState(hs *HostState) error {
	err := s.client.Post("/hosts/self/state", hs, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// sync self
func (s *Synchronization) SynchronizeSelf(hs *HostState) (*HostSettings, error) {
	var r HostSettings
	err := s.client.Post("/hosts/self/sync", hs, &r, nil)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
