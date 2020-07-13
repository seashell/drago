package api

import (
	"fmt"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/google/go-querystring/query"
)

const (
	// Default protocol for communicating with the the Drago server
	DefaultScheme string = "http"
	// Prefix applied to all paths
	DefaultPreprendPath string = "/api"
)

// Config : API Client configuration
type Config struct {
	Address string
	Token   string
}

// Client provides a client to the Drago API
type Client struct {
	config     Config
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


// newRequest :
func (c Client) newRequest(method, path string, body io.Reader, queries interface{}) (*http.Request, error) {
	u := url.URL{
		Scheme: DefaultScheme,
		Host:   c.config.Address,
		Path:   DefaultPreprendPath + path,
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Drago-Token", c.config.Token)
	
	//set queries
	if queries != nil {
		v, err := query.Values(queries)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = v.Encode()
	}
	
	return req, nil
}

// doRequest :
func (c Client) doRequest(req *http.Request, out interface{}) (error){
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if ok := res.StatusCode >= 200 && res.StatusCode < 300; !ok {
		resBody, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return fmt.Errorf("%v: %v", res.Status, string(resBody))
	}

	if err := decodeBody(res, out); err != nil {
		return err
	}

	return nil
}


// encodeBody is used to encode a JSON body
func encodeBody(obj interface{}) (io.Reader, error) {
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
func decodeBody(resp *http.Response, out interface{}) error {
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

// Get :
func (c Client) Get(endpoint string, out interface{}, queries interface{}) error {
	req, err := c.newRequest("GET", endpoint, nil, queries)
	if err != nil {
		return err
	}

	return c.doRequest(req, out)
}

// Post :
func (c Client) Post(endpoint string, in, out interface{}, queries interface{}) error {

	body, err := encodeBody(in)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", endpoint, body, queries)
	if err != nil {
		return err
	}

	return c.doRequest(req, out)
}
