package server

type ServerConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	BindAddrAPI  string `mapstructure:"bind_addr_api"`
	BindAddrUI   string `mapstructure:"bind_addr_ui"`
	Secret       string `mapstructure:"secret"`
	Network      string `mapstructure:"network"`
	MockDataPath string `mapstructure:"mock_data"`
	UI           bool   `mapstructure:"ui"`
}

type server struct {
	config ServerConfig
}

func New(c ServerConfig) (*server, error) {
	s := &server{
		config: c,
	}
	return s, nil
}

func (srv *server) Run() {

	repo, err := NewInMemoryStore()
	if err != nil {
		panic(err)
	}

	PopulateRepositoryWithMockData(repo, srv.config.MockDataPath)

	controller, err := NewController(repo)
	serializer := NewJsonSerializer()
	gw, err := NewGateway(controller, serializer)

	sc := HttpServerConfig{
		BindAddrAPI: srv.config.BindAddrAPI,
		BindAddrUI:  srv.config.BindAddrUI,
		Secret:      []byte(srv.config.Secret),
	}

	s, err := NewHttpServer(gw, sc)
	if err != nil {
		panic(err)
	}

	s.Serve()
}
