package api

type HostSettings struct {
	NetworkInterfaces []NetworkInterface 	`json:"interfaces"`
	WireguardPeers 	  []WireguardPeer 		`json:"peers"`
}

type NetworkInterface struct {
	Name 		*string	`json:"name"`
	ListenPort 	*string `json:"listenPort"`
	Address 	*string `json:"address"`
	Table      	*string `json:"table"`
	DNS        	*string `json:"dns"`
	MTU        	*string `json:"mtu"`
	PreUp      	*string `json:"preUp"`
	PostUp     	*string `json:"postUp"`
	PreDown    	*string `json:"preDown"`
	PostDown   	*string `json:"postDown"`
}

type WireguardPeer struct {
	Interface 			*string 	`json:"interface"`
	Address 			*string		`json:"address"` 
	Port 				*string 	`json:"port"`
	PublicKey 			*string 	`json:"publicKey"`
	AllowedIps 			[]string 	`json:"allowedIps"`
	PersistentKeepalive *int		`json:"persistentKeepalive"`
}

type NetworkInterfaceState struct {
	Name 	  	string `json:"name"`
	WgPublicKey string `json:"publicKey"`
} 

type HostState struct {
	NetworkInterfaces []NetworkInterfaceState `json:"interfaces"`
}

// Hosts is used to query the host-related endpoints.
type Hosts struct {
	client *Client
}

// Settings returns a handle on the settings endpoints.
func (c *Client) Hosts() *Hosts {
	return &Hosts{client: c}
}

// Fetch self settings
func (h Hosts) GetSelfSettings() (*HostSettings, error) {
	var r HostSettings
	err := h.client.Get("/hosts/self/settings", &r)
	if err != nil {
		return nil,err
	}
	return &r,nil
}

// Fetch self settings
func (h Hosts) PostSelfState(hs *HostState) (error) {
		err := h.client.Post("/hosts/self/state", hs, nil)
	if err != nil {
		return err
	}
	return nil
}

// sync self
func (h Hosts) PostSelfSync(hs *HostState) (*HostSettings,error) {
	var r HostSettings
		err := h.client.Post("/hosts/self/sync", hs, &r)
	if err != nil {
		return nil,err
	}
	return &r,nil
}

//cli.Add("GET", "hosts/self/settings", h.GetSelfSettings)
//cli.Add("POST", "hosts/self/state", h.UpdateSelfState)
//cli.Add("POST", "hosts/self/sync", h.SynchronizeSelf)