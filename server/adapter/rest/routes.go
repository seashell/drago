package rest

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterRoutes(e *echo.Echo) {

	e.Add("GET", "/hosts", h.ListHosts)
	e.Add("GET", "/hosts/:id", h.GetHost)
	e.Add("POST", "/hosts", h.CreateHost)
	e.Add("PATCH", "/hosts/:id", h.UpdateHost)
	e.Add("DELETE", "/hosts/:id", h.DeleteHost)

	e.Add("GET", "/links", h.ListLinks)
	e.Add("GET", "/links/:id", h.GetLink)
	e.Add("POST", "/links", h.CreateLink)
	e.Add("PATCH", "/links/:id", h.UpdateLink)
	e.Add("DELETE", "/links/:id", h.DeleteLink)

	e.Add("GET", "/networks", h.ListNetworks)
	e.Add("GET", "/networks/:id", h.GetNetwork)
	e.Add("POST", "/networks", h.CreateNetwork)
	e.Add("PATCH", "/networks/:id", h.UpdateNetwork)
	e.Add("DELETE", "/networks/:id", h.DeleteNetwork)

}
