package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HttpServer interface {
	Serve()
}

type HttpServerConfig struct {
	BindAddr string
	Secret   []byte
}

type httpServer struct {
	gateway *Gateway
	config  HttpServerConfig
	e       *echo.Echo
}

func NewHttpServer(gw *Gateway, c HttpServerConfig) (*httpServer, error) {

	s := &httpServer{
		e:       echo.New(),
		gateway: gw,
		config:  c,
	}

	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())

	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	s.e.HideBanner = true
	s.e.HidePort = true

	v1 := s.e.Group("/api/v1")

	// Hosts API
	v1.Add("GET", "/hosts", EchoHandlerFunc(s.gateway.HandleGetAllHosts))
	v1.Add("GET", "/hosts/:id", EchoHandlerFunc(s.gateway.HandleGetHost))
	v1.Add("POST", "/hosts", EchoHandlerFunc(s.gateway.HandleCreateHost))
	v1.Add("PUT", "/hosts/:id", EchoHandlerFunc(s.gateway.HandleUpdateHost))
	v1.Add("DELETE", "/hosts/:id", EchoHandlerFunc(s.gateway.HandleDeleteHost))

	v1.Add("POST", "/hosts/self/settings", JwtProtected(s.config.Secret)(EchoHandlerFunc(s.gateway.HandleSyncHost)))

	// Links API
	v1.Add("GET", "/links", EchoHandlerFunc(s.gateway.HandleGetAllLinks))
	v1.Add("POST", "/links", EchoHandlerFunc(s.gateway.HandleCreateLink))
	v1.Add("PUT", "/links/:id", EchoHandlerFunc(s.gateway.HandleUpdateLink))
	v1.Add("DELETE", "/links/:id", EchoHandlerFunc(s.gateway.HandleDeleteLink))

	return s, nil
}

func EchoHandlerFunc(f func(Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return f(c)
	}
}

func (s *httpServer) Serve() {
	s.e.Logger.Fatal(s.e.StartServer(&http.Server{
		Addr:         s.config.BindAddr,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 2 * time.Minute,
	}))
}
