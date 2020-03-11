//go:generate echo "==> Bundling web UI"
//go:generate go run github.com/rakyll/statik -f -src=../../ui/build
//go:generate echo "==> Done"

package agent

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rakyll/statik/fs"
	"github.com/seashell/drago/pkg/agent/storage"
	"github.com/seashell/drago/pkg/version"

	client "github.com/seashell/drago/pkg/agent/client"

	log "github.com/sirupsen/logrus"

	_ "github.com/seashell/drago/pkg/agent/statik"
)

var AgentVersion string

func init() {
	AgentVersion = version.GetVersion().VersionNumber()
}

type DragoClaims struct {
	jwt.StandardClaims
	ID    int    `json:"id"`
	Kind  string `json:"kind"`
	Label string `json:"label"`
}

type agent struct {
	store  storage.Store
	config AgentConfig
}

type AgentConfig struct {
	Ui     bool   `mapstructure:"ui"`
	Iface  string `mapstructure:"iface"`
	Server struct {
		Enabled  bool   `mapstructure:"enabled"`
		BindAddr string `mapstructure:"bind_addr"`
		Secret   string `mapstructure:"secret"`
		Network  string `mapstructure:"network"`
	} `mapstructure:"server"`
	Client client.ClientConfig `mapstructure:"client"`
}

func NewAgent(c AgentConfig) (*agent, error) {
	return &agent{
		config: c,
	}, nil
}

func (a *agent) Run() {

	// Create new config store
	store, err := storage.NewInMemoryStore()
	if err != nil {
		log.Error("Failed to initialize config store.")
	}

	a.store = store

	populateStore(a.store)

	// Run drago client
	if a.config.Client.Enabled {
		c, err := client.New(a.config.Client)
		if err != nil {
			log.Error(err)
		}
		c.Run()
	}

	// Run drago server
	if a.config.Server.Enabled {
		go a.serveManagementAPI()
		go a.serveDeviceAPI()
		if a.config.Ui {
			go a.serveUI()
		}
	}
}

func (a *agent) serveUI() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HideBanner = true
	e.HidePort = true

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(statikFS)
	e.GET("/", echo.WrapHandler(fileServer))
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", fileServer)))

	e.Logger.Fatal(e.StartServer(&http.Server{
		Addr:         ":9999",
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 2 * time.Minute,
	}))
}

func (a *agent) serveManagementAPI() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.HideBanner = true
	e.HidePort = true

	h, err := NewHandler(a.store)
	if err != nil {
		log.Fatal(err)
	}

	mgmt := e.Group("/api/v1")

	mgmt.GET("/hosts", h.GetAllHosts)
	mgmt.POST("/hosts", h.CreateHost)
	mgmt.GET("/hosts/:id", h.GetHost)
	mgmt.DELETE("/hosts/:id", h.DeleteHost)
	mgmt.GET("/version", h.GetVersion)

	e.Logger.Fatal(e.StartServer(&http.Server{
		Addr:         ":8080",
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 2 * time.Minute,
	}))
}

func (a *agent) serveDeviceAPI() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(a.config.Server.Secret),
		TokenLookup:   "header:Authorization",
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "client",
		AuthScheme:    "Bearer",
		Claims:        &DragoClaims{},
	}))

	e.HideBanner = true
	e.HidePort = true

	h, err := NewHandler(a.store)
	if err != nil {
		log.Fatal(err)
	}

	e.POST("/api/v1/host", h.SyncHost)

	e.Logger.Fatal(e.StartServer(&http.Server{
		Addr:         ":8181",
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 2 * time.Minute,
	}))

}

func populateStore(s storage.Store) {
	s.InsertHost(&storage.Host{
		Name:          "raspberry-pi-1",
		AdvertiseAddr: "127.0.0.1",
		PublicKey:     "PrlPdb7udp9rliCmQ2z5CerXVLeYWiynLgj32jlhek8=",
		Address:       "192.168.2.1/24",
		ListenPort:    "51820",
	})

	s.InsertHost(&storage.Host{
		Name:      "raspberry-pi-2",
		PublicKey: "nAjpVPq6nS8SpPw7MbajIjIckTun6CwNnSZD30M9njo=",
		Address:   "192.168.2.2/24",
	})

	s.InsertHost(&storage.Host{
		Name:      "raspberry-pi-3",
		PublicKey: "i2tdRKD7ObUrfLmJ4y6ZuBFvLLSTJe1mQvCZELYaqng=",
		Address:   "192.168.2.3/24",
	})

	s.InsertLink(&storage.Link{
		Source:              1,
		Target:              2,
		PersistentKeepalive: 20,
	})

	s.InsertLink(&storage.Link{
		Source:              1,
		Target:              3,
		PersistentKeepalive: 20,
	})

	s.InsertLink(&storage.Link{
		Source:              2,
		Target:              3,
		PersistentKeepalive: 20,
	})
}
