package server

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/seashell/drago/server/adapter/repository"
	"github.com/seashell/drago/server/adapter/rest"
	"github.com/seashell/drago/server/adapter/rest/middleware"
	"github.com/seashell/drago/server/adapter/spa"
	"github.com/seashell/drago/server/application"
	"github.com/seashell/drago/server/controller"
	"github.com/seashell/drago/server/infrastructure/delivery/http"
	"github.com/seashell/drago/server/infrastructure/storage"
	"github.com/seashell/drago/ui"
)

type server struct {
	config     Config
	httpServer *http.Server
}

// Config : Drago server configuration
type Config struct {
	Enabled bool
	DataDir string
	Storage storage.Config
}

func init() {
	os.Setenv("TZ", "UTC")
}

// New : Create a new Drago server
func New(c Config) (*server, error) {

	// Create storage backend
	var backend repository.Backend
	for {
		if b, err := storage.NewBackend(&c.Storage); err == nil {
			backend = b
			break
		} else {
			fmt.Println("Error creating storage backend. Trying again in 2s...")
			fmt.Println(err)
			time.Sleep(2 * time.Second)
		}
	}

	// Create repository adapters for each domain
	networkRepo, err := repository.NewNetworkRepositoryAdapter(backend)
	if err != nil {
		fmt.Println(err)
		panic("Error creating network repository")
	}

	hostRepo, err := repository.NewHostRepositoryAdapter(backend)
	if err != nil {
		fmt.Println(err)
		panic("Error creating host repository")
	}

	ifaceRepo, err := repository.NewInterfaceRepositoryAdapter(backend)
	if err != nil {
		fmt.Println(err)
		panic("Error creating interface repository")
	}

	linkRepo, err := repository.NewLinkRepositoryAdapter(backend)
	if err != nil {
		fmt.Println(err)
		panic("Error creating link repository")
	}

	// Create application services
	ns, err := application.NewNetworkService(networkRepo)
	if err != nil {
		fmt.Println(err)
		panic("Error creating network service")
	}

	hs, err := application.NewHostService(hostRepo)
	if err != nil {
		fmt.Println(err)
		panic("Error creating host service")
	}

	is, err := application.NewInterfaceService(ifaceRepo, networkRepo)
	if err != nil {
		fmt.Println(err)
		panic("Error creating interface service")
	}

	ls, err := application.NewLinkService(linkRepo)
	if err != nil {
		fmt.Println(err)
		panic("Error creating link service")
	}

	ss, err := application.NewSynchronizationService(hostRepo, ifaceRepo, linkRepo)
	if err != nil {
		panic("Error creating link service")
	}

	ts, err := application.NewTokenService(hostRepo)
	if err != nil {
		panic("Error creating token service")
	}

	as, err := application.NewAdmissionService(hostRepo)
	if err != nil {
		panic("Error creating admission service")
	}

	// Create API controller
	ctrl, err := controller.New(ns, hs, is, ls, ss, ts)
	if err != nil {
		fmt.Println(err)
		panic("Error creating controller")
	}

	// Create REST handler
	apiHandler, err := rest.NewHandler(ctrl, rest.Middleware{
		VerifyAuth: middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey: []byte(os.Getenv("ROOT_SECRET")),
			Claims:     jwt.MapClaims{},
		}),
		AdmitHost: middleware.AdmissionMiddleware(as),
	})

	if err != nil {
		fmt.Println(err)
		panic("Error creating API handler")
	}

	// Create SPA adapter to handle static content
	spaHandler, err := spa.NewHandler(ui.Bundle)
	if err != nil {
		fmt.Println(err)
		panic("Error creating SPA handler")
	}

	// Create HTTP server for the API
	httpServer, err := http.NewHTTPServer(&http.ServerConfig{BindAddress: ":8080"})
	if err != nil {
		fmt.Println(err)
		panic("Error creating API server")
	}

	httpServer.RegisterHandler(apiHandler)
	httpServer.RegisterHandler(spaHandler)

	return &server{
		config:     c,
		httpServer: httpServer,
	}, nil
}

// Run : Start Drago server
func (s *server) Run() {
	s.httpServer.Start()
}
