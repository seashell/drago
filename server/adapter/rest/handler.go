package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/seashell/drago/server/controller"
)

const (
	PaginationPageQueryKey    = "page"
	PaginationPerPageQueryKey = "per_page"
)

// Handler :
type Handler struct {
	controller *controller.Controller
	middleware Middleware
}

type Middleware struct {
	VerifyAuth echo.MiddlewareFunc
	AdmitHost  echo.MiddlewareFunc
}

// NewHandler : Create a new REST API handler
func NewHandler(c *controller.Controller, middleware Middleware) (*Handler, error) {
	return &Handler{
		controller: c,
		middleware: middleware,
	}, nil
}
