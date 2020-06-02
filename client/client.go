package client

import (
	"fmt"
	"time"
)

type Config struct {
	Enabled      bool
	DataDir      string
	Servers      []string
	SyncInterval int
}

type client struct {
	config Config
}

func New(c Config) (*client, error) {
	return &client{config: c}, nil
}

func (c *client) Run() {
	go func() {
		for {
			time.Sleep(time.Duration(c.config.SyncInterval) * time.Second)
			fmt.Println("Running client cycle")
		}
	}()
}
