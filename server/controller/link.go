package controller

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
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
	FromInterfaceID     *string  `json:"fromInterfaceId" validate:"required,uuid4"`
	ToInterfaceID       *string  `json:"toInterfaceId" validate:"required,uuid4"`
	AllowedIPs          []string `json:"allowedIPs" validate:"dive,omitempty,cidr"`
	PersistentKeepalive *int     `json:"persistentKeepalive" validate:"omitempty,numeric"`
}

// UpdateLinkInput :
type UpdateLinkInput struct {
	ID                  *string  `json:"id" validate:"required,uuid4"`
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
	NetworkIDFilter         string `query:"networkId" validate:"omitempty,uuid4"`
	SourceHostIDFilter      string `query:"fromHostId" validate:"omitempty,uuid4"`
	SourceInterfaceIDFilter string `query:"fromInterfaceId" validate:"omitempty,uuid4"`
}

// GetLink :
func (c *Controller) GetLink(ctx context.Context, in *GetLinkInput) (*domain.Link, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	l, err := c.ls.GetByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return l, nil
}

// CreateLink :
func (c *Controller) CreateLink(ctx context.Context, in *CreateLinkInput) (*domain.Link, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	l := &domain.Link{}

	errs := model.Copy(l, in)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	res, err := c.ls.Create(l)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// UpdateLink :
func (c *Controller) UpdateLink(ctx context.Context, in *UpdateLinkInput) (*domain.Link, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	l := &domain.Link{}

	errs := model.Copy(l, in)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, errors.Wrap(ErrInternal, err.Error())
	}
	fmt.Println("==== controller link ====")
	spew.Dump(l)
	res, err := c.ls.Update(l)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// DeleteLink :
func (c *Controller) DeleteLink(ctx context.Context, in *DeleteLinkInput) (*domain.Link, error) {
	err := c.v.Struct(in)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidInput, err.Error())
	}

	res, err := c.ls.DeleteByID(in.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternal, err.Error())
	}

	return res, nil
}

// ListLinks :
func (c *Controller) ListLinks(ctx context.Context, in *ListLinksInput) (*pagination.Page, error) {

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

	l := []*domain.Link{}
	p := &domain.Page{}

	if in.SourceInterfaceIDFilter != "" {
		l, p, err = c.ls.FindAllBySourceInterfaceID(in.SourceInterfaceIDFilter, *pageInfo)
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
	} else if in.SourceHostIDFilter != "" {
		l, p, err = c.ls.FindAllBySourceHostID(in.SourceHostIDFilter, *pageInfo)
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
	} else if in.NetworkIDFilter != "" {
		l, p, err = c.ls.FindAllByNetworkID(in.NetworkIDFilter, *pageInfo)
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
	} else {
		l, p, err = c.ls.FindAll(*pageInfo)
		if err != nil {
			return nil, errors.Wrap(ErrInternal, err.Error())
		}
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
