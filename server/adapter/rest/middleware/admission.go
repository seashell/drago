package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/seashell/drago/server/application"
	"github.com/seashell/drago/server/controller"
	"github.com/seashell/drago/server/domain"
)

type (
	AdmissionMiddlewareConfig struct {
		TokenContextKey string
	}
)

var (
	DefaultAdmissionMiddlewareConfig = AdmissionMiddlewareConfig{
		TokenContextKey: DefaultTokenContextKey,
	}
)

func AdmissionMiddleware(as application.AdmissionService) echo.MiddlewareFunc {
	c := DefaultAdmissionMiddlewareConfig
	return AdmissionMiddlewareWithConfig(c, as)
}

func AdmissionMiddlewareWithConfig(config AdmissionMiddlewareConfig, as application.AdmissionService) echo.MiddlewareFunc {

	if config.TokenContextKey == "" {
		config.TokenContextKey = DefaultAdmissionMiddlewareConfig.TokenContextKey
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			token := ctx.Get(config.TokenContextKey).(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)

			// Extract claims
			id := claims["sub"].(string)
			typ := claims["type"].(string)

			iflabels := claims["labels"].([]interface{})
			labels := make([]string, len(iflabels))
			for i, v := range iflabels {
				labels[i] = fmt.Sprint(v)
			}

			if typ != TokenTypeClient {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid credentials")
			}

			in := &controller.CreateHostWithIDInput{
				ID: &id,
			}

			err := validator.New().Struct(in)
			if err != nil {
				fmt.Println(err)
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid credentials")
			}

			host, err := as.GetHostOrCreate(&domain.Host{ID: in.ID, Labels: labels})
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid credentials")
			}

			if host.ID == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching/registering host")
			}

			ctx.Set(HostIDContextKey, *host.ID)

			return next(ctx)
		}
	}
}
