package rest

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	TokenContextKey     = "client"
	TokenTypeManagement = "management"
	TokenTypeClient     = "client"
	tokenHeader         = "X-Drago-Token"
)

func (h *Handler) RegisterRoutes(e *echo.Echo) {

	api := e.Group("/api/")

	secret := os.Getenv("ROOT_SECRET")

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

	cli := api.Group("")
	cli.Use(JWTProtected([]byte(secret)))

	cli.Add("GET", "hosts/self/settings", h.GetSelfSettings)
	cli.Add("POST", "hosts/self/state", h.UpdateSelfState)
	cli.Add("POST", "hosts/self/sync", h.SynchronizeSelf)
}

func JWTProtected(secret []byte) echo.MiddlewareFunc {
	return JWTWithConfig(JWTConfig{
		SigningKey:    secret,
		TokenLookup:   "header:X-Drago-Token",
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    TokenContextKey,
		AuthScheme:    "",
		Claims:        jwt.MapClaims{},
		BeforeFunc: func(c echo.Context) {
			fmt.Println(c.Request().Header.Get("X-Drago-Token"))
		},
	})
}
