package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/seashell/drago/server/controller"
)

// GetNetwork :
func (h *Handler) GetNetwork(c echo.Context) error {
	in := &controller.GetNetworkInput{
		ID: c.Param("id"),
	}

	ctx := c.Request().Context()

	res, err := h.controller.GetNetwork(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// CreateNetwork :
func (h *Handler) CreateNetwork(c echo.Context) error {
	in := &controller.CreateNetworkInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.CreateNetwork(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateNetwork :
func (h *Handler) UpdateNetwork(c echo.Context) error {
	in := &controller.UpdateNetworkInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.UpdateNetwork(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteNetwork :
func (h *Handler) DeleteNetwork(c echo.Context) error {
	in := &controller.DeleteNetworkInput{
		ID: c.Param("id"),
	}

	ctx := c.Request().Context()

	res, err := h.controller.DeleteNetwork(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// ListNetworks :
func (h *Handler) ListNetworks(c echo.Context) error {
	in := &controller.ListNetworksInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.ListNetworks(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}
