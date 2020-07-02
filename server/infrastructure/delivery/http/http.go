package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Handler :
type Handler interface {
	RegisterRoutes(e *echo.Echo)
}

// ServerConfig :
type ServerConfig struct {
	BindAddress string `mapstructure:"bind_address"`
}

// HTTPServer :
type Server struct {
	config *ServerConfig
	echo   *echo.Echo
	ch     chan struct{}
}

// NewHTTPServer :
func NewHTTPServer(c *ServerConfig) (*Server, error) {

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
		config: c,
		echo:   e,
		ch:     make(chan struct{}),
	}

	return server, nil
}

// RegisterHandler :
func (s *Server) RegisterHandler(handler Handler) {
	handler.RegisterRoutes(s.echo)
}

// Start :
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

// Shutdown :
func (s *Server) Shutdown() {
	if s != nil {
		s.echo.Close()
		<-s.ch
	}
}
