package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/seashell/drago/server/controller"
)

// GetLink :
func (h *Handler) GetLink(c echo.Context) error {
	in := &controller.GetLinkInput{
		ID: c.Param("id"),
	}

	ctx := c.Request().Context()

	res, err := h.controller.GetLink(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// CreateLink :
func (h *Handler) CreateLink(c echo.Context) error {
	in := &controller.CreateLinkInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.CreateLink(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateLink :
func (h *Handler) UpdateLink(c echo.Context) error {
	in := &controller.UpdateLinkInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.UpdateLink(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteLink :
func (h *Handler) DeleteLink(c echo.Context) error {
	in := &controller.DeleteLinkInput{
		ID: c.Param("id"),
	}

	ctx := c.Request().Context()

	err := h.controller.DeleteLink(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusNoContent, nil)
}

// ListLinks :
func (h *Handler) ListLinks(c echo.Context) error {
	in := &controller.ListLinksInput{}

	err := c.Bind(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	res, err := h.controller.ListLinks(ctx, in)
	if e := WrapControllerError(err); e != nil {
		return c.JSON(e.Code, e)
	}

	return c.JSON(http.StatusOK, res)
}
