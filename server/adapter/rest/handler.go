package rest

import (
	"github.com/seashell/drago/server/controller"
)

const (
	PaginationPageQueryKey    = "page"
	PaginationPerPageQueryKey = "per_page"
)

// Handler :
type Handler struct {
	controller *controller.Controller
}

// NewHandler : Create a new REST API handler
func NewHandler(c *controller.Controller) (*Handler, error) {
	return &Handler{
		controller: c,
	}, nil
}
