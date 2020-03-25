package server

func PopulateRepository(repo Repository) error {

	repo.CreateHost(&Host{
		Name:             "HOST_1",
		AdvertiseAddress: "192.168.128.2:51820",
		ListenPort:       "51820",
		Address:          "192.168.2.1/24",
	})

	repo.CreateHost(&Host{
		Name:    "HOST_2",
		Address: "192.168.2.2/24",
	})

	repo.CreateHost(&Host{
		Name:    "HOST_3",
		Address: "192.168.2.3/24",
	})

	repo.CreateLink(&Link{
		FromID:              1,
		ToID:                2,
		AllowedIPs:          "192.168.2.2/32",
		PersistentKeepalive: 20,
	})

	repo.CreateLink(&Link{
		FromID:              1,
		ToID:                3,
		AllowedIPs:          "192.168.2.3/32",
		PersistentKeepalive: 20,
	})

	repo.CreateLink(&Link{
		FromID:              2,
		ToID:                1,
		AllowedIPs:          "192.168.2.0/24",
		PersistentKeepalive: 20,
	})

	repo.CreateLink(&Link{
		FromID:              3,
		ToID:                1,
		AllowedIPs:          "192.168.2.0/24",
		PersistentKeepalive: 20,
	})


	return nil
}
