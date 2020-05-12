package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/seashell/drago/server/adapter/rest"
	"github.com/sirupsen/logrus"
)

type HTTPServerConfig struct{}

// HTTPServer
type HTTPServer struct {
	config  *HTTPServerConfig
	handler *rest.Handler
	echo    *echo.Echo
	ch      chan struct{}
}

func NewHTTPServer(handler *rest.Handler, c *HTTPServerConfig) (*HTTPServer, error) {

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	Logger = logrus.New()
	e.Logger = GetEchoLogger()
	e.Use(LoggerMiddleware())

	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"HEAD", "GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	server := &HTTPServer{
		config:  c,
		echo:    e,
		handler: handler,
		ch:      make(chan struct{}),
	}

	server.handler.RegisterRoutes(server.echo)

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
