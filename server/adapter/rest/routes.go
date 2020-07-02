package rest

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes :
func (h *Handler) RegisterRoutes(e *echo.Echo) {

	e.Add("GET", "/healthcheck", h.Healthcheck)

	api := e.Group("/api/")

	mgmt := api.Group("")

	mgmt.Add("GET", "hosts", h.ListHosts)
	mgmt.Add("GET", "hosts/:id", h.GetHost)
	mgmt.Add("POST", "hosts", h.CreateHost)
	mgmt.Add("PATCH", "hosts/:id", h.UpdateHost)
	mgmt.Add("DELETE", "hosts/:id", h.DeleteHost)

	mgmt.Add("GET", "interfaces", h.ListInterfaces)
	mgmt.Add("GET", "interfaces/:id", h.GetInterface)
	mgmt.Add("POST", "interfaces", h.CreateInterface)
	mgmt.Add("PATCH", "interfaces/:id", h.UpdateInterface)
	mgmt.Add("DELETE", "interfaces/:id", h.DeleteInterface)

	mgmt.Add("GET", "links", h.ListLinks)
	mgmt.Add("GET", "links/:id", h.GetLink)
	mgmt.Add("POST", "links", h.CreateLink)
	mgmt.Add("PATCH", "links/:id", h.UpdateLink)
	mgmt.Add("DELETE", "links/:id", h.DeleteLink)

	mgmt.Add("GET", "networks", h.ListNetworks)
	mgmt.Add("GET", "networks/:id", h.GetNetwork)
	mgmt.Add("POST", "networks", h.CreateNetwork)
	mgmt.Add("PATCH", "networks/:id", h.UpdateNetwork)
	mgmt.Add("DELETE", "networks/:id", h.DeleteNetwork)

	mgmt.Add("POST", "tokens", h.CreateToken)
	mgmt.Add("GET", "tokens/self", h.GetSelfToken)

	cli := api.Group("", h.middleware.VerifyAuth, h.middleware.AdmitHost)

	cli.Add("GET", "hosts/self/settings", h.GetSelfSettings)
	cli.Add("POST", "hosts/self/state", h.UpdateSelfState)
	cli.Add("POST", "hosts/self/sync", h.SynchronizeSelf)
}
