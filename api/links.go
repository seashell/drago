package api

// LinkInput :
type LinkInput struct {
	FromInterfaceID     string  	`json:"fromInterfaceId"`
	ToInterfaceID       string  	`json:"toInterfaceId"`
	AllowedIPs          []string	`json:"allowedIPs"`
	PersistentKeepalive int 		`json:"persistentKeepalive"`
	//...
}

// Link :
type Link struct {	
	ID					string 		`json:"id,omitempty"`
	FromInterfaceID     string    	`json:"fromInterfaceId,omitempty"`
	ToInterfaceID       string    	`json:"toInterfaceId,omitempty"`
	AllowedIPs          []string 	`json:"allowedIps,omitempty"`
	PersistentKeepalive	int       	`json:"persistentKeepalive,omitempty"`
	//...
}


// ListLinksInput :
type ListLinksInput struct {
	NetworkIDFilter         string `url:"networkId"`
	SourceHostIDFilter      string `url:"fromHostId"`
	SourceInterfaceIDFilter string `url:"fromInterfaceId"`
}

// LinksList :
type LinksList struct {
	Items []*Link `json:"items"`
}

// Links is used to query the link-related endpoints.
type Links struct {
	client *Client
}

// Links returns a handle on the networks endpoints.
func (c *Client) Links() *Links {
	return &Links{client: c}
}

// ListLinks :
func (i *Links) ListLinks(q ListLinksInput) (*LinksList,error) {
	var r LinksList
	err := i.client.Get("/links", &r, &q)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateLink :
func (i *Links) CreateLink(ni *LinkInput) (*Link,error) {
	var r Link
	err := i.client.Post("/links", ni, &r, nil)
	if err != nil {
		return nil, err
	}
	return &r, nil
}