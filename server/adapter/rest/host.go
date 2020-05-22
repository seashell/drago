package rest

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/seashell/drago/server/controller"
)

// GetHost :
func (h *Handler) GetHost(c echo.Context) error {
	in := &controller.GetHostInput{
		ID: c.Param("id"),
	}

	ctx := c.Request().Context()

	res, err := h.controller.GetHost(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// CreateHost :
func (h *Handler) CreateHost(c echo.Context) error {
	in := &controller.CreateHostInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.CreateHost(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateHost :
func (h *Handler) UpdateHost(c echo.Context) error {
	in := &controller.UpdateHostInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.UpdateHost(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteHost :
func (h *Handler) DeleteHost(c echo.Context) error {
	in := &controller.DeleteHostInput{
		ID: c.Param("id"),
	}

	ctx := c.Request().Context()

	res, err := h.controller.DeleteHost(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)

}

// ListHosts :
func (h *Handler) ListHosts(c echo.Context) error {
	in := &controller.ListHostsInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.ListHosts(ctx, in)

	e := WrapControllerError(err)
	if e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// GetHostSettings :
func (h *Handler) GetSelfSettings(c echo.Context) error {

	token := c.Get(TokenContextKey).(*jwt.Token)
	claims := token.Claims.(DragoClaims)

	in := &controller.GetHostSettingsInput{
		ID: claims.Subject,
	}

	ctx := c.Request().Context()

	res, err := h.controller.GetHostSettings(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}
