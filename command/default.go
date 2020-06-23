package command

var DevConfig = DragoConfig{
	// TODO: define sane defaults
}

var DefaultConfig = DragoConfig{
	UI: false,
	DataDir: "/tmp/drago",
	BindAddr: "0.0.0.0",
	Server: &ServerStanza{
		Enabled: false,
		DataDir: "",
		Storage: &StorageStanza{
			Type:               "",
			Path:               "",
			PostgreSQLAddress:  "",
			PostgreSQLPort:     0,
			PostgreSQLDatabase: "",
			PostgreSQLUsername: "",
			PostgreSQLPassword: "",
			PostgreSQLSSLMode:  "",
		},
	},
	Client: &ClientStanza{
		Enabled: false,
		Servers: nil,
		DataDir: "",
		Token: "",
		SyncInterval: "5s",
	},
}