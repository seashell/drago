package rest

import "github.com/labstack/echo/v4"

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	hosts := e.Group("/hosts")
	hosts.Add("GET", "/", h.ListHosts)
	hosts.Add("GET", "/:id", h.GetHost)
	hosts.Add("POST", "/", h.CreateHost)
}
