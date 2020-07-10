package client

import (
	"github.com/seashell/drago/api"
)

type Synchronization struct {
	c *Client
}

func NewSynchronizationEndpoint(c *Client) *Synchronization {
	return &Synchronization{c: c}
}

func (s *Synchronization) SynchronizeSelf(cs *api.HostState) (*api.HostSettings, error) {
	ts, err := s.c.apiClient.Synchronization().SynchronizeSelf(cs)
	if err != nil {
		return nil, err
	}
	return ts, nil
}
