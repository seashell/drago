package rest

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterRoutes(e *echo.Echo) {

	// If ACL has been boostrapped, protect all routes
	// e.Use(JWTProtected([]byte{}))

	e.Add("GET", "/hosts", h.ListHosts)
	e.Add("GET", "/hosts/:id", h.GetHost)
	e.Add("POST", "/hosts", h.CreateHost)
	e.Add("PATCH", "/hosts/:id", h.UpdateHost)
	e.Add("DELETE", "/hosts/:id", h.DeleteHost)
	//e.Add("GET", "/hosts/self/settings", JWTProtected([]byte{})(h.GetSelfSettings))
	e.Add("GET", "/hosts/self/settings", h.GetSelfSettings)

	e.Add("GET", "/interfaces", h.ListInterfaces)
	e.Add("GET", "/interfaces/:id", h.GetInterface)
	e.Add("POST", "/interfaces", h.CreateInterface)
	e.Add("PATCH", "/interfaces/:id", h.UpdateInterface)
	e.Add("DELETE", "/interfaces/:id", h.DeleteInterface)

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

	//e.Add("POST", "/acl/bootstrap", h.BoostrapAcl)
	//e.Add("GET", "/acl/tokens", h.ListTokens)
	//e.Add("GET", "/acl/tokens/:id", h.GetAclToken)
	//e.Add("POST", "/acl/tokens", h.CreateAclToken)
	//e.Add("PATCH", "/acl/tokens/:id", h.UpdateAclToken)
	//e.Add("DELETE", "/acl/tokens/:id", h.DeleteAclToken)
	//e.Add("GET", "/acl/tokens/self", h.GetSelfToken)
}
