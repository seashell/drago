package middleware

import "github.com/labstack/echo/v4/middleware"

const (
	HostIDContextKey = "hostID"

	TokenTypeManagement = "management"
	TokenTypeClient     = "client"

	DefaultTokenContextKey = "client"
	DefaultTokenHeader     = "X-Drago-Token"
	DefaultSigningMethod   = middleware.AlgorithmHS256
	DefaultTokenLookup     = "header:" + DefaultTokenHeader
	DefaultAuthScheme      = ""
)
