package server

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	gomodel "gopkg.in/jeevatkm/go-model.v1"
)

type Serializer interface {
	SerializeError(Context, string) error
	SerializeHostList(Context, []*Host) error
	SerializeHostDetails(Context, *Host) error
	SerializeLinkList(Context, []*Link) error
	SerializeLinkDetails(Context, *Link) error
	ParseGetHostInput(c Context) (*GetHostInput, error)
	ParseGetAllHostsInput(c Context) (*GetAllHostsInput, error)
	ParseCreateHostInput(c Context) (*CreateHostInput, error)
	ParseUpdateHostInput(c Context) (*UpdateHostInput, error)
	ParseDeleteHostInput(c Context) (*DeleteHostInput, error)
	ParseGetAllLinksInput(c Context) (*GetAllLinksInput, error)
	ParseCreateLinkInput(c Context) (*CreateLinkInput, error)
	ParseUpdateLinkInput(c Context) (*UpdateLinkInput, error)
	ParseDeleteLinkInput(c Context) (*DeleteLinkInput, error)
	ParseSyncHostInput(Context) (*SyncHostInput, error)
	SerializeHostSettings(Context, *Host) error
}

type jsonSerializer struct{}

func NewJsonSerializer() Serializer {
	return &jsonSerializer{}
}

func (p *jsonSerializer) SerializeError(c Context, msg string) error {
	ctx := c.(echo.Context)
	return ctx.JSON(http.StatusInternalServerError, &ApiError{msg})
}

func (p *jsonSerializer) SerializeHostList(c Context, hs []*Host) error {
	ctx := c.(echo.Context)
	if hs != nil {
		res := make([]*HostSummary, 0)
		for _, h := range hs {
			hi := HostSummary{}
			gomodel.Copy(&hi, h)
			res = append(res, &hi)
		}
		return ctx.JSON(http.StatusOK, &HostList{
			Count: len(hs),
			Items: res,
		})
	}
	return ctx.JSON(http.StatusOK, nil)
}

func (p *jsonSerializer) SerializeHostDetails(c Context, h *Host) error {
	ctx := c.(echo.Context)

	if h != nil {
		hd := HostDetails{
			Links: LinkList{
				Count: len(h.Links),
				Items: make([]*LinkDetails, 0),
			},
		}
		gomodel.Copy(&hd, h)

		for _, l := range h.Links {
			ld := &LinkDetails{
				To: &HostSummary{},
			}
			gomodel.Copy(ld, l)
			gomodel.Copy(ld.To, l.To)
			hd.Links.Items = append(hd.Links.Items, ld)
		}
		return ctx.JSON(http.StatusOK, hd)
	}
	return ctx.JSON(http.StatusOK, nil)
}

func (p *jsonSerializer) SerializeHostSettings(c Context, h *Host) error {
	ctx := c.(echo.Context)

	if h != nil {
		hs := HostSettings{
			ID:        h.ID,
			Name:      h.Name,
			CreatedAt: h.UpdatedAt,
			UpdatedAt: h.UpdatedAt,
			Interface: WireguardInterface{},
			Peers:     make([]*WireguardPeer, 0),
		}
		gomodel.Copy(&hs.Interface, h)

		for _, l := range h.Links {
			peer := &WireguardPeer{
				Name:                l.To.Name,
				Endpoint:            l.To.AdvertiseAddress,
				PublicKey:           l.To.PublicKey,
				AllowedIPs:          l.AllowedIPs,
				PersistentKeepalive: l.PersistentKeepalive,
			}
			hs.Peers = append(hs.Peers, peer)
		}
		return ctx.JSON(http.StatusOK, hs)
	}
	return ctx.JSON(http.StatusOK, nil)
}

func (p *jsonSerializer) ParseSyncHostInput(c Context) (*SyncHostInput, error) {
	ctx := c.(echo.Context)

	cli := ctx.Get("client").(*jwt.Token)
	claims := cli.Claims.(*DragoClaims)

	in := &SyncHostInput{}
	if err := ctx.Bind(in); err != nil {
		return nil, err
	}

	in.ID = claims.ID
	in.HostType = claims.Kind

	return in, nil
}

func (p *jsonSerializer) ParseCreateHostInput(c Context) (*CreateHostInput, error) {
	ctx := c.(echo.Context)

	in := &CreateHostInput{}
	if err := ctx.Bind(in); err != nil {
		return nil, err
	}

	return in, nil
}

func (p *jsonSerializer) ParseGetHostInput(c Context) (*GetHostInput, error) {
	ctx := c.(echo.Context)

	in := &GetHostInput{}
	if err := ctx.Bind(in); err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	in.ID = id

	return in, nil
}

func (p *jsonSerializer) ParseGetAllHostsInput(c Context) (*GetAllHostsInput, error) {
	ctx := c.(echo.Context)

	in := &GetAllHostsInput{}
	if err := ctx.Bind(in); err != nil {
		return nil, err
	}

	return in, nil
}

func (p *jsonSerializer) ParseUpdateHostInput(c Context) (*UpdateHostInput, error) {
	ctx := c.(echo.Context)

	in := &UpdateHostInput{}
	if err := ctx.Bind(in); err != nil {
		return nil, err
	}

	return in, nil
}

func (p *jsonSerializer) ParseDeleteHostInput(c Context) (*DeleteHostInput, error) {
	ctx := c.(echo.Context)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, err
	}

	return &DeleteHostInput{
		ID: id,
	}, nil
}

func (p *jsonSerializer) ParseCreateLinkInput(c Context) (*CreateLinkInput, error) {
	ctx := c.(echo.Context)

	in := &CreateLinkInput{}
	if err := ctx.Bind(in); err != nil {
		return nil, err
	}

	return in, nil
}

func (p *jsonSerializer) ParseGetAllLinksInput(c Context) (*GetAllLinksInput, error) {
	ctx := c.(echo.Context)

	in := &GetAllLinksInput{}
	if err := ctx.Bind(in); err != nil {
		return nil, err
	}

	return in, nil
}

func (p *jsonSerializer) SerializeLinkList(c Context, ls []*Link) error {
	ctx := c.(echo.Context)

	if ls != nil {

		res := make([]*LinkDetails, 0)

		for _, l := range ls {
			ld := &LinkDetails{
				From: &HostSummary{},
				To:   &HostSummary{},
			}
			gomodel.Copy(ld, l)
			gomodel.Copy(ld.To, l.To)
			gomodel.Copy(ld.From, l.From)
			res = append(res, ld)
		}

		return ctx.JSON(http.StatusOK, &LinkList{
			Count: len(res),
			Items: res,
		})
	}
	return ctx.JSON(http.StatusOK, nil)
}

func (p *jsonSerializer) SerializeLinkDetails(c Context, l *Link) error {
	ctx := c.(echo.Context)
	if l != nil {
		ld := &LinkDetails{}
		gomodel.Copy(&ld, l)
		return ctx.JSON(http.StatusOK, ld)
	}
	return ctx.JSON(http.StatusOK, nil)
}

func (p *jsonSerializer) ParseUpdateLinkInput(c Context) (*UpdateLinkInput, error) {
	ctx := c.(echo.Context)

	in := &UpdateLinkInput{}
	if err := ctx.Bind(in); err != nil {
		return nil, err
	}

	return in, nil
}

func (p *jsonSerializer) ParseDeleteLinkInput(c Context) (*DeleteLinkInput, error) {
	ctx := c.(echo.Context)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, err
	}

	return &DeleteLinkInput{
		ID: id,
	}, nil
}
