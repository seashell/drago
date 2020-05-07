package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/seashell/drago/server/adapter/repository"
	"github.com/seashell/drago/server/adapter/rest"
	"github.com/seashell/drago/server/application"
	"github.com/seashell/drago/server/controller"
	"github.com/seashell/drago/server/infrastructure/storage"
	"github.com/seashell/drago/server/migrations/postgresql"
)

type HTTPServerConfig struct {
	APIBindAddr string
}

// HTTPServer
type HTTPServer struct {
	config         *HTTPServerConfig
	restAPIHandler *rest.Handler
	echo           *echo.Echo
	ch             chan struct{}
}

func NewHTTPServer(c *HTTPServerConfig) (*HTTPServer, error) {

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.AddTrailingSlash())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"HEAD", "GET"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

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
	handler := rest.NewHandler(ctrl)

	server := &HTTPServer{
		config:         c,
		echo:           e,
		restAPIHandler: handler,
		ch:             make(chan struct{}),
	}

	server.restAPIHandler.RegisterRoutes(server.echo)

	return server, nil
}

func (s *HTTPServer) Start() {
	go func() {
		defer close(s.ch)
		s.echo.Logger.Fatal(s.echo.StartServer(&http.Server{
			Addr:         ":8080",
			ReadTimeout:  2 * time.Minute,
			WriteTimeout: 2 * time.Minute,
		}))
	}()
}

func (s *HTTPServer) Shutdown() {
	if s != nil {
		s.echo.Close()
		<-s.ch
	}
}
