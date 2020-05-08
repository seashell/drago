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
	Name             *string `json:"name" validate:"required,stringlength(1|50)"`
	IPAddress        *string `json:"ipAddress" validate:"required,cidr"`
	AdvertiseAddress *string `json:"advertiseAddress" validate:"required,cidr"`
	ListenPort       *string `json:"listenPort" validate:"required,numeric,stringlength(1|5)"`
	PublicKey        *string `json:"publicKey" validate:""`
	Table            *string `json:"table" validate:""`
	DNS              *string `json:"dns" validate:""`
	MTU              *string `json:"mtu" validate:""`
	PreUp            *string `json:"preUp" validate:""`
	PostUp           *string `json:"postUp" validate:""`
	PreDown          *string `json:"preDown" validate:""`
	PostDown         *string `json:"postDown" validate:""`
	NetworkID        *string `json:"networkId" validate:"required,uuid4"`
}

// UpdateHostInput :
type UpdateHostInput struct {
	ID               *string `json:"id" validate:"required,uuid4"`
	Name             *string `json:"name" validate:"stringlength(1|50)"`
	IPAddress        *string `json:"ipAddress" validate:"cidr"`
	AdvertiseAddress *string `json:"advertiseAddress" validate:"cidr"`
	ListenPort       *string `json:"listenPort" validate:"numeric,stringlength(1|5)"`
	PublicKey        *string `json:"publicKey" validate:""`
	Table            *string `json:"table" validate:""`
	DNS              *string `json:"dns" validate:""`
	MTU              *string `json:"mtu" validate:""`
	PreUp            *string `json:"preUp" validate:""`
	PostUp           *string `json:"postUp" validate:""`
	PreDown          *string `json:"preDown" validate:""`
	PostDown         *string `json:"postDown" validate:""`
}

// DeleteHostInput :
type DeleteHostInput struct {
	ID string `validate:"required,uuid4"`
}

// ListHostsInput :
type ListHostsInput struct {
	pagination.Input
	NetworkIDFilter string `query:"networkId" validate:"required,uuid4"`
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
func (c *Controller) CreateHost(ctx context.Context, in *CreateHostInput) (*string, error) {
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

	id, err := c.hs.Create(h)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return id, nil
}

// UpdateHost :
func (c *Controller) UpdateHost(ctx context.Context, in *UpdateHostInput) (*string, error) {
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

	id, err := c.hs.Update(h)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return id, nil
}

// DeleteHost :
func (c *Controller) DeleteHost(ctx context.Context, in *DeleteHostInput) error {
	err := c.v.Struct(in)
	if err != nil {
		return errors.Wrap(ErrInvalidInput, err.Error())
	}

	err = c.hs.DeleteByID(in.ID)
	if err != nil {
		return errors.Wrap(ErrInternal, err.Error())
	}

	return nil
}

// ListHosts :
func (c *Controller) ListHosts(ctx context.Context, in *ListHostsInput) (*pagination.Page, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	pageInfo := &domain.PageInfo{
		Page:    in.Page,
		PerPage: in.PerPage,
	}

	h, p, err := c.hs.FindAllByNetworkID(in.NetworkIDFilter, *pageInfo)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
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
