package server

type ServerConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	BindAddr string `mapstructure:"bind_addr"`
	Secret   string `mapstructure:"secret"`
	Network  string `mapstructure:"network"`
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

	PopulateRepository(repo)

	controller, err := NewController(repo)
	serializer := NewJsonSerializer()
	gw, err := NewGateway(controller, serializer)

	sc := HttpServerConfig{
		BindAddr: srv.config.BindAddr,
		Secret:   []byte(srv.config.Secret),
	}

	s, err := NewHttpServer(gw, sc)
	if err != nil {
		panic(err)
	}

	s.Serve()
}
