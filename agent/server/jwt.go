package server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type DragoClaims struct {
	jwt.StandardClaims
	ID    int    `json:"id"`
	Kind  string `json:"kind"`
	Label string `json:"label"`
}

func JwtProtected(secret []byte) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(secret),
		TokenLookup:   "header:Authorization",
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "client",
		AuthScheme:    "Bearer",
		Claims:        &DragoClaims{},
	})
}
