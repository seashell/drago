package rest

import (
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

func JWTProtected(secret []byte) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    secret,
		TokenLookup:   "header:" + tokenHeader,
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    TokenContextKey,
		AuthScheme:    "",
		Claims:        &jwt.StandardClaims{},
	})
}
