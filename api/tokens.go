package api

// TokenInput :
type TokenInput struct {
	Type    string  `json:"type"`
	Subject string  `json:"subject"`
	Labels  []string `json:"labels"`
	//...
}

// Token :
type Token struct {	
	ID		string 		`json:"id,omitempty"`
	Type    string  	`json:"type,omitempty"`
	Subject	string  	`json:"subject,omitempty"`
	Labels 	[]string	`json:"labels,omitempty"`
	Raw     string   	`json:"secret,omitempty"`
	//...
}


// ListTokensInput :
type ListTokensInput struct {
	NetworkIDFilter	string `url:"networkId"`
}

// TokensList :
type TokensList struct {
	Items []*Token `json:"items"`
}

// Tokens is used to query the link-related endpoints.
type Tokens struct {
	client *Client
}

// Tokens returns a handle on the networks endpoints.
func (c *Client) Tokens() *Tokens {
	return &Tokens{client: c}
}

// ListTokens :
func (i *Tokens) ListTokens(q ListTokensInput) (*TokensList,error) {
	var r TokensList
	err := i.client.Get("/tokens", &r, &q)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateToken :
func (i *Tokens) CreateToken(ni *TokenInput) (*Token,error) {
	var r Token
	err := i.client.Post("/tokens", ni, &r, nil)
	if err != nil {
		return nil, err
	}
	return &r, nil
}