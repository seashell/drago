package agent

import (
	"fmt"

	"github.com/seashell/drago/client"
	"github.com/seashell/drago/server"
	"github.com/seashell/drago/version"
)

var AgentVersion string

func init() {
	AgentVersion = version.GetVersion().VersionNumber()
}

type Agent interface {
	Run()
}

type agent struct {
	config Config
}

type Config struct {
	UI      bool
	DataDir string
	Server  server.Config
	Client  client.Config
}

func New(c Config) (*agent, error) {
	return &agent{
		config: c,
	}, nil
}

func (a *agent) Run() {

	if a.config.Client.Enabled {
		fmt.Println("Initializing agent (client)")
		client, err := client.New(a.config.Client)
		if err != nil {
			panic(err)
		}
		client.Run()
	}

	if a.config.Server.Enabled {
		fmt.Println("Initializing agent (server)")
		s, err := server.New(a.config.Server)
		if err != nil {
			panic(err)
		}
		s.Run()
	}
}
