package server

func PopulateRepository(repo Repository) error {

	repo.CreateHost(&Host{
		Name:             "bounce-server",
		AdvertiseAddress: "152.74.12.1",
		ListenPort:       "51820",
		Address:          "192.168.2.1/24",
	})

	repo.CreateHost(&Host{
		Name:    "raspberry-pi-2",
		Address: "192.168.2.2/24",
	})

	repo.CreateHost(&Host{
		Name:    "raspberry-pi-3",
		Address: "192.168.2.3/24",
	})

	repo.CreateHost(&Host{
		Name:    "raspberry-pi-4",
		Address: "192.168.2.4/24",
	})

	repo.CreateLink(&Link{
		FromID:              1,
		ToID:                2,
		AllowedIPs:          "192.168.2.2/32",
		PersistentKeepalive: 20,
	})

	repo.CreateLink(&Link{
		FromID:              2,
		ToID:                1,
		AllowedIPs:          "192.168.2.1/24",
		PersistentKeepalive: 20,
	})

	repo.CreateLink(&Link{
		FromID:              1,
		ToID:                3,
		AllowedIPs:          "192.168.2.3/32",
		PersistentKeepalive: 20,
	})

	repo.CreateLink(&Link{
		FromID:              3,
		ToID:                1,
		AllowedIPs:          "192.168.2.1/24",
		PersistentKeepalive: 20,
	})

	repo.CreateLink(&Link{
		FromID:              4,
		ToID:                1,
		AllowedIPs:          "192.168.2.1/24",
		PersistentKeepalive: 20,
	})

	repo.CreateLink(&Link{
		FromID:              1,
		ToID:                4,
		AllowedIPs:          "192.168.2.4/32",
		PersistentKeepalive: 20,
	})

	return nil
}
