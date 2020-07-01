package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Healthcheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, nil)
}
