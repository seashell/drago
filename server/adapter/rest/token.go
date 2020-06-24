package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/seashell/drago/server/controller"
)

// CreateToken :
func (h *Handler) CreateToken(c echo.Context) error {
	in := &controller.CreateTokenInput{}
	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.CreateToken(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// GetSelfToken :
func (h *Handler) GetSelfToken(c echo.Context) error {

	raw := c.Request().Header.Get("X-Drago-Token")

	in := &controller.GetSelfTokenInput{
		Raw: &raw,
	}

	ctx := c.Request().Context()

	res, err := h.controller.GetSelfToken(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}
