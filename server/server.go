package server

import (
	"github.com/seashell/drago/server/adapter/repository"
	"github.com/seashell/drago/server/adapter/rest"
	"github.com/seashell/drago/server/adapter/static"
	"github.com/seashell/drago/server/application"
	"github.com/seashell/drago/server/controller"
	"github.com/seashell/drago/server/infrastructure/delivery/http"
	"github.com/seashell/drago/server/infrastructure/storage"
	"github.com/seashell/drago/server/migrations/postgresql"
	"github.com/seashell/drago/ui"
)

type server struct {
	config    ServerConfig
	apiServer *http.Server
	uiServer  *http.Server
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
		panic("Error connecting to storage backend")
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
	restHandler, err := rest.NewHandler(ctrl)
	if err != nil {
		panic("Error creating API handler")
	}

	// Create HTTP server for the API
	apiServer, err := http.NewHTTPServer(restHandler, &http.ServerConfig{BindAddress: ":8080"})
	if err != nil {
		panic("Error creating API server")
	}

	// Create static content handler
	staticHandler, err := static.NewHandler(ui.Bundle)
	if err != nil {
		panic("Error creating API handler")
	}

	// Create HTTP server for static files (UI)
	uiServer, err := http.NewHTTPServer(staticHandler, &http.ServerConfig{BindAddress: ":8000"})
	if err != nil {
		panic("Error creating UI server")
	}

	return &server{
		config:    c,
		apiServer: apiServer,
		uiServer:  uiServer,
	}, nil

}

func (s *server) Run() {
	s.apiServer.Start()
	s.uiServer.Start()
}
