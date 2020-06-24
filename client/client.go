package client

import (
	"time"
	"github.com/seashell/drago/client/host"
	"github.com/seashell/drago/api"
)

type client struct {
	config     Config
	hostClient	*host.Client
}

type Config struct {
	Enabled bool
	DataDir string
	Servers []string
	Token string
	SyncInterval time.Duration
}

func New(c Config) (*client, error) {
	h,err := host.New(&host.Config{
		DataDir: c.DataDir,
		SyncInterval: c.SyncInterval,
	})
	if err != nil {
		return nil,err
	}

	a,err := api.NewClient(&api.Config{
		Address: c.Servers[0], //TODO: add support for multiple API addresses
		Token: c.Token,
	})
	if err != nil {
		return nil,err
	}
	h.SetAPIClient(a)

	return &client{
		config:     c,
		hostClient: h,
	}, nil
}

func (c *client) Run() {
	c.hostClient.Start()
}

