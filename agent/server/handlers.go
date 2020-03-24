package server

type Gateway struct {
	controller Controller
	serializer Serializer
}

func NewGateway(c Controller, s Serializer) (*Gateway, error) {
	return &Gateway{
		controller: c,
		serializer: s,
	}, nil
}

func (gw *Gateway) HandleSyncHost(c Context) error {

	in, err := gw.serializer.ParseSyncHostInput(c)

	h, err := gw.controller.SyncHost(in)

	if err != nil {
		return gw.serializer.SerializeError(c, "Error message")
	}

	return gw.serializer.SerializeHostSettings(c, h)
}

func (gw *Gateway) HandleGetAllHosts(c Context) error {
	in, err := gw.serializer.ParseGetAllHostsInput(c)
	if err != nil {
		panic(err)
	}
	hs, err := gw.controller.GetAllHosts(in)
	if err != nil {
		return gw.serializer.SerializeError(c, "Error message")
	}
	return gw.serializer.SerializeHostList(c, hs)
}

func (gw *Gateway) HandleGetHost(c Context) error {
	in, err := gw.serializer.ParseGetHostInput(c)
	if err != nil {
		panic(err)
	}
	h, err := gw.controller.GetHost(in)

	if err != nil {
		return gw.serializer.SerializeError(c, "Error message")
	}
	return gw.serializer.SerializeHostDetails(c, h)
}

func (gw *Gateway) HandleCreateHost(c Context) error {
	in, err := gw.serializer.ParseCreateHostInput(c)
	if err != nil {
		panic(err)
	}

	h, err := gw.controller.CreateHost(in)
	if err != nil {
		return gw.serializer.SerializeError(c, "Error message")
	}

	return gw.serializer.SerializeHostDetails(c, h)
}

func (gw *Gateway) HandleUpdateHost(c Context) error {
	in, err := gw.serializer.ParseUpdateHostInput(c)
	if err != nil {
		panic(err)
	}
	h, err := gw.controller.UpdateHost(in)
	if err != nil {
		return gw.serializer.SerializeError(c, "Error message")
	}
	return gw.serializer.SerializeHostDetails(c, h)
}

func (gw *Gateway) HandleDeleteHost(c Context) error {
	in, err := gw.serializer.ParseDeleteHostInput(c)
	if err != nil {
		panic(err)
	}
	err = gw.controller.DeleteHost(in)
	if err != nil {
		return gw.serializer.SerializeError(c, "Error message")
	}
	return gw.serializer.SerializeHostDetails(c, nil)
}

func (gw *Gateway) HandleGetAllLinks(c Context) error {
	in, err := gw.serializer.ParseGetAllLinksInput(c)
	if err != nil {
		panic(err)
	}
	ls, err := gw.controller.GetAllLinks(in)
	if err != nil {
		return gw.serializer.SerializeError(c, "Error message")
	}
	return gw.serializer.SerializeLinkList(c, ls)
}

func (gw *Gateway) HandleCreateLink(c Context) error {
	in, err := gw.serializer.ParseCreateLinkInput(c)
	if err != nil {
		panic(err)
	}
	l, err := gw.controller.CreateLink(in)
	if err != nil {
		return gw.serializer.SerializeError(c, "Error message")
	}
	return gw.serializer.SerializeLinkDetails(c, l)
}

func (gw *Gateway) HandleUpdateLink(c Context) error {
	in, err := gw.serializer.ParseUpdateLinkInput(c)
	if err != nil {
		panic(err)
	}
	l, err := gw.controller.UpdateLink(in)
	if err != nil {
		return gw.serializer.SerializeError(c, "Error message")
	}

	return gw.serializer.SerializeLinkDetails(c, l)
}

func (gw *Gateway) HandleDeleteLink(c Context) error {
	in, err := gw.serializer.ParseDeleteLinkInput(c)
	if err != nil {
		panic(err)
	}
	err = gw.controller.DeleteLink(in)
	if err != nil {
		return gw.serializer.SerializeError(c, "Error message")
	}
	return gw.serializer.SerializeLinkDetails(c, nil)
}
