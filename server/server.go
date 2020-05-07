//go:generate go generate github.com/seashell/drago/server/ui

package server

import (
	"github.com/seashell/drago/server/adapter/repository"
	"github.com/seashell/drago/server/adapter/rest"
	"github.com/seashell/drago/server/application"
	"github.com/seashell/drago/server/controller"
	"github.com/seashell/drago/server/infrastructure/delivery/http"
	"github.com/seashell/drago/server/infrastructure/storage"
	"github.com/seashell/drago/server/migrations/postgresql"
)

type server struct {
	config     ServerConfig
	httpServer *http.HTTPServer
}

type ServerConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

func New(c ServerConfig) (*server, error) {

	// Create storage backend
	backend, err := storage.NewBackend(&storage.Config{
		Type:               "postgresql",
		PostgreSQLAddress:  "127.0.0.1",
		PostgreSQLPort:     5432,
		PostgreSQLDatabase: "seashell",
		PostgreSQLUsername: "admin",
		PostgreSQLPassword: "password",
		PostgreSQLSSLMode:  "disable",
	})
	if err != nil {
		panic("Error creating backend")
	}

	// Apply migrations
	backend.ApplyMigrations(postgresql.Migrations...)

	// Create repository adapters for each domain
	networkRepo, err := repository.NewPostgreSQLNetworkRepositoryAdapter(backend)
	if err != nil {
		panic("Error creating network repository")
	}

	hostRepo, err := repository.NewPostgreSQLHostRepositoryAdapter(backend)
	if err != nil {
		panic("Error creating host repository")
	}

	linkRepo, err := repository.NewPostgreSQLLinkRepositoryAdapter(backend)
	if err != nil {
		panic("Error creating link repository")
	}

	// Create application services
	ns, err := application.NewNetworkService(networkRepo)
	if err != nil {
		panic("Error creating network service")
	}

	hs, err := application.NewHostService(hostRepo)
	if err != nil {
		panic("Error creating host service")
	}

	ls, err := application.NewLinkService(linkRepo)
	if err != nil {
		panic("Error creating link service")
	}

	// Create API controller
	ctrl, err := controller.New(ns, hs, ls)
	if err != nil {
		panic("Error creating controller")
	}

	// Create REST handler
	handler, err := rest.NewHandler(ctrl)
	if err != nil {
		panic("Error creating API handler")
	}

	// Create HTTP server
	httpServer, err := http.NewHTTPServer(handler, &http.HTTPServerConfig{})
	if err != nil {
		panic("Error creating HTTP server")
	}

	return &server{
		config:     c,
		httpServer: httpServer,
	}, nil

}

func (s *server) Run() {
	s.httpServer.Start()
}
