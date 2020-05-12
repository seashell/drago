package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler interface {
	RegisterRoutes(e *echo.Echo)
}

type ServerConfig struct {
	BindAddress string `mapstructure:"bind_address"`
}

// HTTPServer
type Server struct {
	config  *ServerConfig
	handler Handler
	echo    *echo.Echo
	ch      chan struct{}
}

func NewHTTPServer(handler Handler, c *ServerConfig) (*Server, error) {

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"HEAD", "GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	server := &Server{
		config:  c,
		echo:    e,
		handler: handler,
		ch:      make(chan struct{}),
	}

	server.handler.RegisterRoutes(server.echo)

	return server, nil
}

func (s *Server) Start() {
	go func() {
		defer close(s.ch)
		s.echo.Logger.Fatal(s.echo.StartServer(&http.Server{
			Addr:         s.config.BindAddress,
			ReadTimeout:  2 * time.Minute,
			WriteTimeout: 2 * time.Minute,
		}))
	}()
}

func (s *Server) Shutdown() {
	if s != nil {
		s.echo.Close()
		<-s.ch
	}
}
