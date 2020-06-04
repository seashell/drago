package server

import (
	"fmt"

	"github.com/seashell/drago/server/adapter/repository"
	"github.com/seashell/drago/server/adapter/rest"
	"github.com/seashell/drago/server/adapter/spa"
	"github.com/seashell/drago/server/application"
	"github.com/seashell/drago/server/controller"
	"github.com/seashell/drago/server/infrastructure/delivery/http"
	"github.com/seashell/drago/server/infrastructure/storage"
	"github.com/seashell/drago/server/migrations/postgresql"
	"github.com/seashell/drago/ui"
)

type server struct {
	config    Config
	apiServer *http.Server
	uiServer  *http.Server
}

type Config struct {
	Enabled bool
	DataDir string
	Storage storage.Config
}

func New(c Config) (*server, error) {

	// Create storage backend

	backend, err := storage.NewBackend(&c.Storage)
	if err != nil {
		fmt.Println(err)
		panic("Error connecting to storage backend")
	}

	// Apply migrations
	backend.ApplyMigrations(postgresql.Migrations...)NewBackendNewBackend

	// Create repository adapters for each domain
	networkRepo, err := repository.NewNetworkRepositoryAdapter(backend)
	if err != nil {
		fmt.Println(err)
		panic("Error creating network repository")
	}

	hostRepo, err := repository.NewHostRepositoryAdapter(backend)
	if err != nil {
		panic("Error creating host repository")
	}

	ifaceRepo, err := repository.NewInterfaceRepositoryAdapter(backend)
	if err != nil {
		panic("Error creating interface repository")
	}

	linkRepo, err := repository.NewLinkRepositoryAdapter(backend)
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

	is, err := application.NewInterfaceService(ifaceRepo, networkRepo)
	if err != nil {
		panic("Error creating interface service")
	}

	ls, err := application.NewLinkService(linkRepo)
	if err != nil {
		panic("Error creating link service")
	}

	// Create API controller
	ctrl, err := controller.New(ns, hs, is, ls)
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

	// Create SPA adapter to handle static content
	spaHandler, err := spa.NewHandler(ui.Bundle)
	if err != nil {
		panic("Error creating SPA handler")
	}

	// Create HTTP server for static files (UI)
	uiServer, err := http.NewHTTPServer(spaHandler, &http.ServerConfig{BindAddress: ":8000"})
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
