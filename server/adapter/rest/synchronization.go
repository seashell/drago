package rest

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/seashell/drago/server/controller"
)

// SynchronizeSelf :
func (h *Handler) SynchronizeSelf(c echo.Context) error {

	token := c.Get(TokenContextKey).(*jwt.Token)
	claims := token.Claims.(jwt.StandardClaims)

	in := &controller.SynchronizeHostInput{
		ID: claims.Subject,
	}

	ctx := c.Request().Context()

	res, err := h.controller.SynchronizeHost(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// GetSelfSettings :
func (h *Handler) GetSelfSettings(c echo.Context) error {

	token := c.Get(TokenContextKey).(*jwt.Token)
	claims := token.Claims.(jwt.StandardClaims)

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

// UpdateSelfState :
func (h *Handler) UpdateSelfState(c echo.Context) error {

	token := c.Get(TokenContextKey).(*jwt.Token)
	claims := token.Claims.(jwt.StandardClaims)

	in := &controller.UpdateHostStateInput{
		ID: claims.Subject,
	}

	ctx := c.Request().Context()

	state, err := h.controller.UpdateHostState(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, state)
}

// GetSettingsByHostID :
func (h *Handler) GetHostSettings(c echo.Context) error {
	in := &controller.GetHostSettingsInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.GetHostSettings(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateStateByHostID :
func (h *Handler) UpdateHostState(c echo.Context) error {
	in := &controller.UpdateHostStateInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.UpdateHostState(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}
