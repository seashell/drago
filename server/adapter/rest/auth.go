package rest

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const TokenContextKey = "client"
const TokenTypeManagement = "management"
const TokenTypeClient = "client"

type DragoClaims struct {
	jwt.StandardClaims
	Type     string   `json:"type"`
	Policies []string `json:"policies"`
	Label    string   `json:"label"`
}

func JWTProtected(secret []byte) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte("secret"),
		TokenLookup:   "header:X-Drago-Token",
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    TokenContextKey,
		AuthScheme:    "",
		Claims:        &DragoClaims{},
	})
}
