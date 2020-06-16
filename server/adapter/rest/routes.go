package rest

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterRoutes(e *echo.Echo) {

	api := e.Group("/api/")

	// If ACL has been boostrapped, protect all routes
	// api.Use(JWTProtected([]byte{}))

	api.Add("GET", "hosts", h.ListHosts)
	api.Add("GET", "hosts/:id", h.GetHost)
	api.Add("POST", "hosts", h.CreateHost)
	api.Add("PATCH", "hosts/:id", h.UpdateHost)
	api.Add("DELETE", "hosts/:id", h.DeleteHost)
	//api.Add("GET", "hosts/self/settings", JWTProtected([]byte{})(h.GetSelfSettings))
	api.Add("GET", "hosts/self/settings", h.GetSelfSettings)

	api.Add("GET", "interfaces", h.ListInterfaces)
	api.Add("GET", "interfaces/:id", h.GetInterface)
	api.Add("POST", "interfaces", h.CreateInterface)
	api.Add("PATCH", "interfaces/:id", h.UpdateInterface)
	api.Add("DELETE", "interfaces/:id", h.DeleteInterface)

	api.Add("GET", "links", h.ListLinks)
	api.Add("GET", "links/:id", h.GetLink)
	api.Add("POST", "links", h.CreateLink)
	api.Add("PATCH", "links/:id", h.UpdateLink)
	api.Add("DELETE", "links/:id", h.DeleteLink)

	api.Add("GET", "networks", h.ListNetworks)
	api.Add("GET", "networks/:id", h.GetNetwork)
	api.Add("POST", "networks", h.CreateNetwork)
	api.Add("PATCH", "networks/:id", h.UpdateNetwork)
	api.Add("DELETE", "networks/:id", h.DeleteNetwork)

	//api.Add("POST", "acl/bootstrap", h.BoostrapAcl)
	//api.Add("GET", "acl/tokens", h.ListTokens)
	//api.Add("GET", "acl/tokens/:id", h.GetAclToken)
	//api.Add("POST", "acl/tokens", h.CreateAclToken)
	//api.Add("PATCH", "acl/tokens/:id", h.UpdateAclToken)
	//api.Add("DELETE", "acl/tokens/:id", h.DeleteAclToken)
	//api.Add("GET", "acl/tokens/self", h.GetSelfToken)
}
