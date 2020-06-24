package api

import (
	"bytes"
	"errors"
	"io"
	"encoding/json"
	"net/http"
	"net/url"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

const (
	DefaultScheme      	string = "http"
	DefaultPreprendPath	string = "/api"
)

type Config struct {
	Address string
	Token string
}

// Client provides a client to the Drago API
type Client struct {
	config Config
	httpClient *http.Client
}

// NewClient returns a new client
func NewClient(c *Config) (*Client, error) {
	h := cleanhttp.DefaultClient()

	client := &Client{
		config:     *c,
		httpClient: h,
	}
	return client, nil
}

func (c Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	u := url.URL{
		Scheme: DefaultScheme,
		Host: c.config.Address,
		Path: DefaultPreprendPath+path,
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil,err
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Drago-Token", c.config.Token)
	return req, nil
}

func (c Client) Get(endpoint string, out interface{}) (error) {
	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := decodeBody(resp, out); err != nil {
		return err
	}

	return nil
}

func (c Client) Post(endpoint string, in, out interface{}) (error) {

	body, err := encodeBody(in)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", endpoint, body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := decodeBody(resp, out); err != nil {
		return err
	}

	return nil
}

// encodeBody is used to encode a JSON body
func encodeBody (obj interface{}) (io.Reader, error) {
		if reader, ok := obj.(io.Reader); ok {
			return reader, nil
		}
	
		buf := bytes.NewBuffer(nil)
		enc := json.NewEncoder(buf)
		if err := enc.Encode(obj); err != nil {
			return nil, err
		}
		return buf, nil
	}


// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface{}) (error) {
	switch resp.ContentLength {
	case 0:
		if out == nil {
			return nil
		}
		return errors.New("Got 0 byte response with non-nil decode object")
	default:
		dec := json.NewDecoder(resp.Body)
		return dec.Decode(out)
	}
}