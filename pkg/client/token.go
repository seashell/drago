package client

import "context"

// Token resource methods
func (d *DragoAPIClient) GetToken(ctx context.Context, t *Token) (*string, error) {
	receiver := struct {
		Token *string `json:"secret"`
	}{}

	err := d.createResource(
		tokensPath,
		t,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.Token, nil
}
