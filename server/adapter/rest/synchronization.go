package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/seashell/drago/server/adapter/rest/middleware"
	"github.com/seashell/drago/server/controller"
)

// SynchronizeSelf :
func (h *Handler) SynchronizeSelf(c echo.Context) error {

	hostID := c.Get(middleware.HostIDContextKey).(string)

	in := &controller.SynchronizeHostInput{
		ID: hostID,
	}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
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

	hostID := c.Get(middleware.HostIDContextKey).(string)

	in := &controller.GetHostSettingsInput{
		ID: hostID,
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

	hostID := c.Get(middleware.HostIDContextKey).(string)

	in := &controller.UpdateHostStateInput{
		ID: hostID,
	}
	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	state, err := h.controller.UpdateHostState(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, state)
}
