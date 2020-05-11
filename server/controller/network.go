package controller

import (
	"context"

	"github.com/pkg/errors"
	"github.com/seashell/drago/server/controller/pagination"
	"github.com/seashell/drago/server/domain"
)

// GetNetworkInput :
type GetNetworkInput struct {
	ID string `validate:"required,uuid4"`
}

// CreateNetworkInput :
type CreateNetworkInput struct {
	Name           *string `json:"name" validate:"required,min=1,max=50"`
	IPAddressRange *string `json:"ipAddressRange" validate:"required,cidr"`
}

// UpdateNetworkInput :
type UpdateNetworkInput struct {
	ID             *string `json:"id" validate:"required,uuid4"`
	Name           *string `json:"name" validate:"omitempty,min=1,max=50"`
	IPAddressRange *string `json:"ipAddressRange" validate:"omitempty,cidr"`
}

// DeleteNetworkInput :
type DeleteNetworkInput struct {
	ID string `validate:"required,uuid4"`
}

// ListNetworksInput :
type ListNetworksInput struct {
	pagination.Input
}

// GetNetwork :
func (c *Controller) GetNetwork(ctx context.Context, in *GetNetworkInput) (*domain.Network, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	n, err := c.ns.GetByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return n, nil
}

// CreateNetwork :
func (c *Controller) CreateNetwork(ctx context.Context, in *CreateNetworkInput) (*domain.Network, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	n := &domain.Network{
		Name:           in.Name,
		IPAddressRange: in.IPAddressRange,
	}

	res, err := c.ns.Create(n)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// UpdateNetwork :
func (c *Controller) UpdateNetwork(ctx context.Context, in *UpdateNetworkInput) (*domain.Network, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	n := &domain.Network{
		ID:             in.ID,
		Name:           in.Name,
		IPAddressRange: in.IPAddressRange,
	}

	res, err := c.ns.Update(n)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// DeleteNetwork :
func (c *Controller) DeleteNetwork(ctx context.Context, in *DeleteNetworkInput) (*domain.Network, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	res, err := c.ns.DeleteByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// ListNetworks :
func (c *Controller) ListNetworks(ctx context.Context, in *ListNetworksInput) (*pagination.Page, error) {

	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	if in.Page == 0 {
		in.Page = pagination.PAGINATION_DEFAULT_PAGE
		in.PerPage = pagination.PAGINATION_DEFAULT_PER_PAGE
	}

	pageInfo := &domain.PageInfo{
		Page:    in.Page,
		PerPage: in.PerPage,
	}

	n, p, err := c.ns.FindAll(*pageInfo)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	page := &pagination.Page{
		Page:       p.Page,
		PerPage:    p.PerPage,
		PageCount:  p.PageCount,
		TotalCount: p.TotalCount,
		Items:      n,
	}

	return page, nil
}
