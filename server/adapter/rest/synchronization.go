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
	claims := token.Claims.(DragoClaims)

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
	claims := token.Claims.(DragoClaims)

	in := &controller.GetHostSettingsInput{
		ID: claims.Subject,
	}

	ctx := c.Request().Context()

	res, err := h.controller.GetHostSettingsByID(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateSelfState :
func (h *Handler) UpdateSelfState(c echo.Context) error {

	token := c.Get(TokenContextKey).(*jwt.Token)
	claims := token.Claims.(DragoClaims)

	in := &controller.UpdateHostStateInput{
		ID: claims.Subject,
	}

	ctx := c.Request().Context()

	err := h.controller.UpdateHostStateByID(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, nil)
}
