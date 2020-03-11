package agent

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/seashell/drago/pkg/agent/storage"
)

type Handler struct {
	store storage.Store
}

func NewHandler(s storage.Store) (*Handler, error) {
	return &Handler{
		store: s,
	}, nil
}

// ********************************
// *********  MGMT API  ***********
// ********************************

func (h *Handler) GetHost(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	ht, err := h.store.SelectHost(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, Error{"Requested resource does not exist"})
	}
	return c.JSON(http.StatusOK, ht)
}

func (h *Handler) GetAllHosts(c echo.Context) error {

	hs, _ := h.store.SelectAllHosts()
	ret := make([]*Host, 0)

	// Map stored object to JSON
	for _, v := range hs {

		ht := &Host{
			Entity: Entity{
				ID:        v.ID,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			},
			Name:          v.Name,
			PublicKey:     v.PublicKey,
			AdvertiseAddr: v.AdvertiseAddr,
			Interface: WireguardInterface{
				Address:    v.Address,
				ListenPort: v.ListenPort,
			},
		}
		ret = append(ret, ht)
	}

	return c.JSON(http.StatusOK, ret)
}

func (h *Handler) CreateHost(c echo.Context) error {

	ht := Host{}
	if err := c.Bind(&ht); err != nil {
		return c.JSON(http.StatusBadRequest, Error{"Bad request"})
	}

	nht, err := h.store.InsertHost(&storage.Host{Name: ht.Name})
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{err.Error()})
	}

	jwt, err := generatwJwtToken("host", nht.ID, nht.Name)

	return c.JSON(http.StatusOK, Host{
		Entity: Entity{
			ID:        nht.ID,
			CreatedAt: nht.CreatedAt,
			UpdatedAt: nht.UpdatedAt,
		},
		Name: nht.Name,
		Jwt:  jwt,
	})
}

func (h *Handler) DeleteHost(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.store.DeleteHost(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, Error{"Requested resource does not exist"})
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, &Version{
		Version: AgentVersion,
	})
}

func generatwJwtToken(k string, i int, l string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &DragoClaims{
		ID:    i,
		Kind:  k,
		Label: l,
	})

	t, err := token.SignedString([]byte("seashell-test-secret-123"))

	if err != nil {
		return "", err
	}

	return t, nil
}

// ********************************
// *********  HOST API  ***********
// ********************************

func (h *Handler) SyncHost(c echo.Context) error {

	cli := c.Get("client").(*jwt.Token)

	claims := cli.Claims.(*DragoClaims)

	if tp := claims.Kind; tp != "host" {
		return c.JSON(http.StatusBadRequest, Error{"Not a valid host"})
	}

	id := claims.ID

	ht, err := h.store.SelectHost(id)

	if err != nil {
		return c.JSON(http.StatusNotFound, Error{"Could not find object matching the presented claims"})
	}

	peers, err := h.store.SelectAllPeersForHost(id)

	retPeers := make([]WireguardPeer, 0)

	for _, v := range peers {
		retPeers = append(retPeers, WireguardPeer{
			Endpoint:            v.Endpoint,
			AllowedIPs:          v.AllowedIPs,
			PublicKey:           v.PublicKey,
			PersistentKeepalive: v.PersistentKeepalive,
		})
	}

	ret := &Host{
		Entity: Entity{
			ID:        ht.ID,
			CreatedAt: ht.CreatedAt,
			UpdatedAt: ht.UpdatedAt,
		},

		Name:          ht.Name,
		PublicKey:     ht.PublicKey,
		AdvertiseAddr: ht.AdvertiseAddr,

		Interface: WireguardInterface{
			Address:    ht.Address,
			ListenPort: ht.ListenPort,
		},
		Peers: retPeers,
	}

	return c.JSON(http.StatusOK, ret)
}
