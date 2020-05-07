//go:generate go generate github.com/seashell/drago/server/ui

package server

import "github.com/seashell/drago/server/infrastructure/delivery/http"

type server struct {
	config     ServerConfig
	httpServer *http.HTTPServer
}

type ServerConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

func New(c ServerConfig) (*server, error) {
	return &server{
		config: c,
	}, nil

}

func (s *server) Run() {

}
