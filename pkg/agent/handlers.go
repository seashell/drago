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

func (h *Handler) GetNode(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	n, err := h.store.SelectNode(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, Error{"Requested resource does not exist"})
	}
	return c.JSON(http.StatusOK, n)
}

func (h *Handler) GetAllNodes(c echo.Context) error {

	n, _ := h.store.SelectAllNodes()
	ret := make([]*Node, 0)

	// Map stored object to JSON
	for _, v := range n {

		n := &Node{
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
		ret = append(ret, n)
	}

	return c.JSON(http.StatusOK, ret)
}

func (h *Handler) CreateNode(c echo.Context) error {

	n := Node{}
	if err := c.Bind(&n); err != nil {
		return c.JSON(http.StatusBadRequest, Error{"Bad request"})
	}

	nn, err := h.store.InsertNode(&storage.Node{Name: n.Name})
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{err.Error()})
	}

	jwt, err := generatwJwtToken("node", nn.ID, nn.Name)

	return c.JSON(http.StatusOK, Node{
		Entity: Entity{
			ID:        nn.ID,
			CreatedAt: nn.CreatedAt,
			UpdatedAt: nn.UpdatedAt,
		},
		Name: nn.Name,
		Jwt:  jwt,
	})
}

func (h *Handler) DeleteNode(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.store.DeleteNode(id)
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
// *********  NODE API  ***********
// ********************************

func (h *Handler) SyncNode(c echo.Context) error {

	cli := c.Get("client").(*jwt.Token)

	claims := cli.Claims.(*DragoClaims)

	if tp := claims.Kind; tp != "node" {
		return c.JSON(http.StatusBadRequest, Error{"Not a valid node"})
	}

	id := claims.ID

	n, err := h.store.SelectNode(id)

	if err != nil {
		return c.JSON(http.StatusNotFound, Error{"Could not find object matching the presented claims"})
	}

	peers, err := h.store.SelectAllPeersForNode(id)

	retPeers := make([]WireguardPeer, 0)

	for _, v := range peers {
		retPeers = append(retPeers, WireguardPeer{
			Endpoint:            v.Endpoint,
			AllowedIPs:          v.AllowedIPs,
			PublicKey:           v.PublicKey,
			PersistentKeepalive: v.PersistentKeepalive,
		})
	}

	ret := &Node{
		Entity: Entity{
			ID:        n.ID,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
		},

		Name:          n.Name,
		PublicKey:     n.PublicKey,
		AdvertiseAddr: n.AdvertiseAddr,

		Interface: WireguardInterface{
			Address:    n.Address,
			ListenPort: n.ListenPort,
		},
		Peers: retPeers,
	}

	return c.JSON(http.StatusOK, ret)
}
