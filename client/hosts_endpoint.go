package client

import (
	"github.com/seashell/drago/api"
)

type Hosts struct {
	c *Client
}

func NewHostsEndpoint(c *Client) *Hosts {
	return &Hosts{c: c}
}

func (h *Hosts) Sync(cs *api.HostState) (*api.HostSettings, error) {
	ts, err := h.c.apiClient.Hosts().PostSelfSync(cs)
	if err != nil {
		return nil, err
	}
	return ts, nil
}
