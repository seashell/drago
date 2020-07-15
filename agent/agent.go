package agent

import (
	"fmt"

	"github.com/seashell/drago/client"
	"github.com/seashell/drago/server"
	"github.com/seashell/drago/version"

	"github.com/seashell/drago/agent/logger"
)

var AgentVersion string

func init() {
	AgentVersion = version.GetVersion().VersionNumber()
}

type Config struct {
	UI      bool
	DataDir string
	Server  server.Config
	Client  client.Config
}

type Agent interface {
	Run()
}

type agent struct {
	config Config
}

func New(c Config) (*agent, error) {
	return &agent{
		config: c,
	}, nil
}

func (a *agent) Run() {
	
	logger, err := logger.New(logger.Configuration{
		Level: logger.Debug,
	})
	if err != nil {
		fmt.Println(err)
		panic("Error creating logger")
	}

	if a.config.Server.Enabled {
		logger.Infof("Initializing server agent\n")
		s, err := server.New(a.config.Server)
		if err != nil {
			panic(err)
		}
		s.Run()
	}

	if a.config.Client.Enabled {
		logger.Infof("Initializing client agent\n")
		c, err := client.New(a.config.Client, logger)
		if err != nil {
			panic(err)
		}
		c.Run()
	}
}
