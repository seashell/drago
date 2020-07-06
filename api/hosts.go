package api

// HostInput :
type HostInput struct {
	Name			string 		`json:"name"`
	Labels			[]string 	`json:"labels,omitempty"`
	AdvertiseAddress	string 	`json:"advertiseAddress,omitempty"`
}

// Host :
type Host struct {	
	ID					string 		`json:"id,omitempty"`
	Name				string 		`json:"name,omitempty"`
	Labels				[]string 	`json:"labels,omitempty"`
	AdvertiseAddress	string 		`json:"advertiseAddress,omitempty"`
	//...
}

// HostsList :
type HostsList struct {
	Items []*Host `json:"items"`
}

// Hosts is used to query the network-related endpoints.
type Hosts struct {
	client *Client
}

// Hosts returns a handle on the networks endpoints.
func (c *Client) Hosts() *Hosts {
	return &Hosts{client: c}
}

// ListHosts :
func (h *Hosts) ListHosts() (*HostsList,error) {
	var r HostsList
	err := h.client.Get("/hosts", &r, nil)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateHost :
func (h *Hosts) CreateHost(nh *HostInput) (*Host,error) {
	var r Host
	err := h.client.Post("/hosts", nh, &r, nil)
	if err != nil {
		return nil, err
	}
	return &r, nil
}