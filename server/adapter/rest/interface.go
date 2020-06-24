package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/seashell/drago/server/controller"
)

// GetInterface :
func (h *Handler) GetInterface(c echo.Context) error {
	in := &controller.GetInterfaceInput{
		ID: c.Param("id"),
	}

	ctx := c.Request().Context()

	res, err := h.controller.GetInterface(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// CreateInterface :
func (h *Handler) CreateInterface(c echo.Context) error {
	in := &controller.CreateInterfaceInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.CreateInterface(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateInterface :
func (h *Handler) UpdateInterface(c echo.Context) error {
	in := &controller.UpdateInterfaceInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.UpdateInterface(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteInterface :
func (h *Handler) DeleteInterface(c echo.Context) error {
	in := &controller.DeleteInterfaceInput{
		ID: c.Param("id"),
	}

	ctx := c.Request().Context()

	res, err := h.controller.DeleteInterface(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)

}

// ListInterfaces :
func (h *Handler) ListInterfaces(c echo.Context) error {
	in := &controller.ListInterfacesInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.ListInterfaces(ctx, in)

	e := WrapControllerError(err)
	if e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}
