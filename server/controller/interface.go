package controller

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"gopkg.in/jeevatkm/go-model.v1"

	"github.com/seashell/drago/server/controller/pagination"
	"github.com/seashell/drago/server/domain"
)

// GetInterfaceInput :
type GetInterfaceInput struct {
	ID string `validate:"required,uuid4"`
}

// CreateInterfaceInput :
type CreateInterfaceInput struct {
	Name       *string `json:"name" validate:"required,min=1,max=50"`
	HostID     *string `json:"hostId" validate:"required,uuid4"`
	NetworkID  *string `json:"networkId" validate:"omitempty,uuid4"`
	IPAddress  *string `json:"ipAddress" validate:"omitempty,cidr"`
	ListenPort *string `json:"listenPort" validate:"omitempty,numeric,min=1,max=5"`
	PublicKey  *string `json:"publicKey" validate:""`
	Table      *string `json:"table" validate:""`
	DNS        *string `json:"dns" validate:""`
	MTU        *string `json:"mtu" validate:""`
	PreUp      *string `json:"preUp" validate:""`
	PostUp     *string `json:"postUp" validate:""`
	PreDown    *string `json:"preDown" validate:""`
	PostDown   *string `json:"postDown" validate:""`
}

// UpdateInterfaceInput :
type UpdateInterfaceInput struct {
	ID         *string `json:"id" validate:"required,uuid4"`
	NetworkID  *string `json:"networkId" validate:"omitempty,uuid4"`
	Name       *string `json:"name" validate:"min=1,max=50"`
	IPAddress  *string `json:"ipAddress" validate:"omitempty,cidr"`
	ListenPort *string `json:"listenPort" validate:"omitempty,numeric,min=1,max=5"`
	PublicKey  *string `json:"publicKey" validate:""`
	Table      *string `json:"table" validate:""`
	DNS        *string `json:"dns" validate:""`
	MTU        *string `json:"mtu" validate:""`
	PreUp      *string `json:"preUp" validate:""`
	PostUp     *string `json:"postUp" validate:""`
	PreDown    *string `json:"preDown" validate:""`
	PostDown   *string `json:"postDown" validate:""`
}

// DeleteInterfaceInput :
type DeleteInterfaceInput struct {
	ID string `validate:"required,uuid4"`
}

// ListInterfacesInput :
type ListInterfacesInput struct {
	pagination.Input
	HostIDFilter    string `query:"hostId" validate:"omitempty,uuid4"`
	NetworkIDFilter string `query:"networkId" validate:"omitempty,uuid4"`
}

// GetInterface :
func (c *Controller) GetInterface(ctx context.Context, in *GetInterfaceInput) (*domain.Interface, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	h, err := c.is.GetByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return h, nil
}

// CreateInterface :
func (c *Controller) CreateInterface(ctx context.Context, in *CreateInterfaceInput) (*domain.Interface, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	iface := &domain.Interface{}

	errs := model.Copy(iface, in)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	res, err := c.is.Create(iface)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// UpdateInterface :
func (c *Controller) UpdateInterface(ctx context.Context, in *UpdateInterfaceInput) (*domain.Interface, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	h := &domain.Interface{}

	errs := model.Copy(h, in)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	res, err := c.is.Update(h)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// DeleteInterface :
func (c *Controller) DeleteInterface(ctx context.Context, in *DeleteInterfaceInput) (*domain.Interface, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	res, err := c.is.DeleteByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// ListInterfaces :
func (c *Controller) ListInterfaces(ctx context.Context, in *ListInterfacesInput) (*pagination.Page, error) {

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

	i := []*domain.Interface{}
	p := &domain.Page{}

	if in.HostIDFilter != "" {
		i, p, err = c.is.FindAllByHostID(in.HostIDFilter, *pageInfo)
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
	} else if in.NetworkIDFilter != "" {
		i, p, err = c.is.FindAllByNetworkID(in.NetworkIDFilter, *pageInfo)
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
	} else {
		i, p, err = c.is.FindAll(*pageInfo)
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
	}

	page := &pagination.Page{
		Page:       p.Page,
		PerPage:    p.PerPage,
		PageCount:  p.PageCount,
		TotalCount: p.TotalCount,
		Items:      i,
	}

	return page, nil
}
