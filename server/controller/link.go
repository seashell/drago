package controller

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/controller/pagination"
	"github.com/seashell/drago/server/domain"
	"gopkg.in/jeevatkm/go-model.v1"
)

// GetLinkInput :
type GetLinkInput struct {
	ID string `validate:"required,uuid4"`
}

// CreateLinkInput :
type CreateLinkInput struct {
	NetworkID           *string `json:"networkId" validate:"required,uuid4"`
	ToHostID            *string `json:"toHostId" validate:"required,uuid4"`
	FromHostID          *string `json:"fromhostId" validate:"required,uuid4"`
	AllowedIPs          *string `json:"allowedIPs" validate:"dive,omitempty,cidr"`
	PersistentKeepalive *int    `json:"persistentKeepalive" validate:"numeric"`
}

// UpdateLinkInput :
type UpdateLinkInput struct {
	AllowedIPs          []string `json:"allowedIPs" validate:"dive,omitempty,cidr"`
	PersistentKeepalive *int     `json:"persistentKeepalive"`
}

// DeleteLinkInput :
type DeleteLinkInput struct {
	ID string `validate:"required,uuid4"`
}

// ListLinksInput :
type ListLinksInput struct {
	pagination.Input
	NetworkIDFilter string `query:"networkId" validate:"uuid4"`
	HostIDFilter    string `query:"hostId" validate:"uuid4"`
}

// GetLink :
func (c *Controller) GetLink(ctx context.Context, in *GetLinkInput) (*domain.Link, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidInput.Error())
	}

	l, err := c.ls.GetByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(err, ErrInternal.Error())
	}

	return l, nil
}

// CreateLink :
func (c *Controller) CreateLink(ctx context.Context, in *CreateLinkInput) (*string, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidInput.Error())
	}

	l := &domain.Link{}

	errs := model.Copy(l, in)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, errors.Wrap(err, ErrInternal.Error())
	}

	id, err := c.ls.Create(l)
	if err != nil {
		return nil, errors.Wrap(err, ErrInternal.Error())
	}

	return id, nil
}

// UpdateLink :
func (c *Controller) UpdateLink(ctx context.Context, in *UpdateLinkInput) (*string, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidInput.Error())
	}

	l := &domain.Link{}

	errs := model.Copy(l, in)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, errors.Wrap(err, ErrInternal.Error())
	}

	id, err := c.ls.Update(l)
	if err != nil {
		return nil, errors.Wrap(err, ErrInternal.Error())
	}

	return id, nil
}

// DeleteLink :
func (c *Controller) DeleteLink(ctx context.Context, in *DeleteLinkInput) error {
	err := c.v.Struct(in)
	if err != nil {
		return errors.Wrap(ErrInvalidInput, err.Error())
	}

	err = c.ls.DeleteByID(in.ID)
	if err != nil {
		return errors.Wrap(ErrInternal, err.Error())
	}

	return nil
}

// ListLinks :
func (c *Controller) ListLinks(ctx context.Context, in *ListLinksInput) (*pagination.Page, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidInput.Error())
	}

	pageInfo := &domain.PageInfo{
		Page:    in.Page,
		PerPage: in.PerPage,
	}

	l := []*domain.Link{}
	p := &domain.Page{}

	if in.NetworkIDFilter != "" {
		l, p, err = c.ls.FindAllByNetworkID(in.NetworkIDFilter, *pageInfo)
		if err != nil {
			return nil, errors.Wrap(err, ErrInternal.Error())
		}
	} else if in.HostIDFilter != "" {
		l, p, err = c.ls.FindAllByHostID(in.NetworkIDFilter, *pageInfo)
		if err != nil {
			return nil, errors.Wrap(err, ErrInternal.Error())
		}
	} else {
		err = errors.New("No filter provided")
		return nil, errors.Wrap(err, ErrInvalidInput.Error())
	}

	page := &pagination.Page{
		Page:       p.Page,
		PerPage:    p.PerPage,
		PageCount:  p.PageCount,
		TotalCount: p.TotalCount,
		Items:      l,
	}

	return page, nil
}
