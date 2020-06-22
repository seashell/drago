package controller

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"gopkg.in/jeevatkm/go-model.v1"

	"github.com/seashell/drago/server/controller/pagination"
	"github.com/seashell/drago/server/domain"
)

// GetHostInput :
type GetHostInput struct {
	ID string `validate:"required,uuid4"`
}

// CreateHostInput :
type CreateHostInput struct {
	Name             *string  `json:"name" validate:"required,min=1,max=50"`
	AdvertiseAddress *string  `json:"advertiseAddress" validate:"omitempty,cidr|hostname"`
	Labels           []string `json:"labels" validate:"dive,omitempty,alphanum"`
}

// UpdateHostInput :
type UpdateHostInput struct {
	ID               *string  `json:"id" validate:"required,uuid4"`
	Name             *string  `json:"name" validate:"omitempty,min=1,max=50"`
	AdvertiseAddress *string  `json:"advertiseAddress" validate:"omitempty,cidr|hostname"`
	Labels           []string `json:"labels" validate:"dive,omitempty,alphanum"`
}

// DeleteHostInput :
type DeleteHostInput struct {
	ID string `validate:"required,uuid4"`
}

// ListHostsInput :
type ListHostsInput struct {
	pagination.Input
	NetworkIDFilter string `query:"networkId" validate:"omitempty,uuid4"`
}

// GetHost :
func (c *Controller) GetHost(ctx context.Context, in *GetHostInput) (*domain.Host, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	h, err := c.hs.GetByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return h, nil
}

// CreateHost :
func (c *Controller) CreateHost(ctx context.Context, in *CreateHostInput) (*domain.Host, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	h := &domain.Host{}

	errs := model.Copy(h, in)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	res, err := c.hs.Create(h)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// UpdateHost :
func (c *Controller) UpdateHost(ctx context.Context, in *UpdateHostInput) (*domain.Host, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	h := &domain.Host{}

	errs := model.Copy(h, in)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	res, err := c.hs.Update(h)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// DeleteHost :
func (c *Controller) DeleteHost(ctx context.Context, in *DeleteHostInput) (*domain.Host, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	res, err := c.hs.DeleteByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// ListHosts :
func (c *Controller) ListHosts(ctx context.Context, in *ListHostsInput) (*pagination.Page, error) {

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

	h := []*domain.Host{}
	p := &domain.Page{}

	if in.NetworkIDFilter != "" {
		h, p, err = c.hs.FindAllByNetworkID(in.NetworkIDFilter, *pageInfo)
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
	} else {
		h, p, err = c.hs.FindAll(*pageInfo)
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
	}

	page := &pagination.Page{
		Page:       p.Page,
		PerPage:    p.PerPage,
		PageCount:  p.PageCount,
		TotalCount: p.TotalCount,
		Items:      h,
	}

	return page, nil
}
