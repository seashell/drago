package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/hashicorp/go-cleanhttp"
)

// Client provides a client to the Drago API
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient returns a new Drago API client
func NewClient(config *Config) (*Client, error) {

	config = DefaultConfig().Merge(config)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	client := &Client{
		config:     config,
		httpClient: cleanhttp.DefaultClient(),
	}

	return client, nil
}

func (c *Client) getResource(p string, id string, receiver interface{}) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += path.Join(p, id)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	c.addHeaders(req)

	return c.doRequest(req, receiver)
}

func (c *Client) createResource(p string, sender interface{}, receiver interface{}) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += p + "/"

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(sender)

	req, err := http.NewRequest("POST", u.String(), b)
	if err != nil {
		return err
	}

	c.addHeaders(req)

	return c.doRequest(req, receiver)
}

func (c *Client) updateResource(id, p string, sender interface{}, receiver interface{}) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += path.Join(p, id)

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(sender)

	req, err := http.NewRequest("PATCH", u.String(), b)
	if err != nil {
		return err
	}

	c.addHeaders(req)
	return c.doRequest(req, receiver)
}

func (c *Client) deleteResource(id, p string, receiver interface{}) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += path.Join(p, id)

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	c.addHeaders(req)

	return c.doRequest(req, receiver)
}

func (c *Client) listResources(p string, filters map[string][]string, receiver interface{}) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += p

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	c.addQuery(filters, req)
	c.addHeaders(req)

	return c.doRequest(req, receiver)
}

func (c *Client) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Drago-Token", c.config.Token)
}

func (c *Client) addQuery(filters map[string][]string, req *http.Request) {

	q := req.URL.Query()

	for k, values := range filters {
		for _, v := range values {
			q.Add(k, v)
		}
	}

	req.URL.RawQuery = q.Encode()
}

func (c *Client) doRequest(req *http.Request, receiver interface{}) error {

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if ok := res.StatusCode >= 200 && res.StatusCode < 300; !ok {

		err := CodedError{
			Code: res.StatusCode,
		}

		if err := json.NewDecoder(res.Body).Decode(&err); err != nil {
			resBody, _ := ioutil.ReadAll(res.Body)
			return fmt.Errorf("%v (%v)", res.Status, string(resBody))
		}

		return err
	}

	if receiver != nil {
		return json.NewDecoder(res.Body).Decode(receiver)
	}

	return nil
}
