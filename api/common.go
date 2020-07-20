package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) getResource(path string, id string, receiver interface{}) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += path
	if len(id) > 0 {
		u.Path += "/" + id
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	c.addHeaders(req)

	return c.doRequest(req, receiver)
}

func (c *Client) createResource(path string, sender interface{}, receiver interface{}) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += path

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(sender)

	req, err := http.NewRequest("POST", u.String(), b)
	if err != nil {
		return err
	}

	c.addHeaders(req)

	return c.doRequest(req, receiver)
}

func (c *Client) updateResource(id, path string, sender interface{}) error {

	base, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	base.Path += path
	base.Path += "/" + id

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(sender)

	req, err := http.NewRequest("PATCH", base.String(), b)
	if err != nil {
		return err
	}

	c.addHeaders(req)
	return c.doRequest(req, nil)
}

func (c *Client) deleteResource(id, path string) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += path
	u.Path += "/" + id

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	c.addHeaders(req)
	return c.doRequest(req, nil)
}

func (c *Client) listResources(path string, page *Page, filter map[string]string, receiver interface{}) error {
	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += path

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	c.addQuery(filter, page, req)
	c.addHeaders(req)

	return c.doRequest(req, receiver)
}

func (c *Client) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Drago-Token", c.config.Token)
}

func (c *Client) addQuery(query map[string]string, page *Page, req *http.Request) {
	q := req.URL.Query()

	if len(query) > 0 {
		for k, v := range query {
			q.Add(k, v)
		}
	}

	if page != nil {
		if page.Page != nil {
			q.Add(pageQueryParamKey, strconv.Itoa(*page.Page))
		}
		if page.PerPage != nil {
			q.Add(perPageQueryParamKey, strconv.Itoa(*page.PerPage))
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
		resBody, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("%v: %v", res.Status, string(resBody))
	}

	if receiver != nil {
		return json.NewDecoder(res.Body).Decode(receiver)
	}

	return nil
}
