package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Common methods
func (d *DragoAPIClient) getResource(
	id, path string,
	receiver interface{},
) error {
	req, err := d.getGetRequestPreamble(id, path)
	if err != nil {
		return err
	}

	d.addHeaders(req)
	return d.doRequest(receiver, req)
}

func (d *DragoAPIClient) createResource(
	path string,
	sender interface{},
	receiver interface{},
) error {
	req, err := d.getPostRequestPreamble(path, sender)
	if err != nil {
		return err
	}

	d.addHeaders(req)
	return d.doRequest(receiver, req)
}

func (d *DragoAPIClient) updateResource(
	id, path string,
	sender interface{},
	receiver interface{},
) error {
	req, err := d.getPatchRequestPreamble(id, path, sender)
	if err != nil {
		return err
	}

	d.addHeaders(req)
	return d.doRequest(receiver, req)
}

func (d *DragoAPIClient) deleteResource(
	id, path string,
) error {
	req, err := d.getDeleteRequestPreamble(id, path)
	if err != nil {
		return err
	}

	d.addHeaders(req)
	return d.doRequest(nil, req)
}

func (d *DragoAPIClient) listResources(
	path string,
	page *Page,
	filter map[string]string,
	receiver interface{},
) error {
	req, err := d.getGetRequestPreamble("", path)
	if err != nil {
		return err
	}

	d.addQueries(filter, page, req)
	d.addHeaders(req)
	return d.doRequest(receiver, req)
}

func (d *DragoAPIClient) getGetRequestPreamble(id, path string) (*http.Request, error) {
	base, err := url.Parse(d.ServerURL)
	if err != nil {
		return nil, err
	}

	base.Path += path
	if len(id) > 0 {
		base.Path += "/" + id
	}

	req, err := http.NewRequest("GET", base.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, err
}

func (d *DragoAPIClient) getPostRequestPreamble(path string, sender interface{}) (*http.Request, error) {
	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(sender)

	base, err := url.Parse(d.ServerURL)
	if err != nil {
		return nil, err
	}

	base.Path += path

	req, err := http.NewRequest("POST", base.String(), b)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (d *DragoAPIClient) getPatchRequestPreamble(id, path string, sender interface{}) (*http.Request, error) {
	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(sender)

	base, err := url.Parse(d.ServerURL)
	if err != nil {
		return nil, err
	}

	base.Path += path
	base.Path += "/" + id

	req, err := http.NewRequest("PATCH", base.String(), b)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (d *DragoAPIClient) getDeleteRequestPreamble(id, path string) (*http.Request, error) {
	base, err := url.Parse(d.ServerURL)
	if err != nil {
		return nil, err
	}

	base.Path += path
	base.Path += "/" + id

	req, err := http.NewRequest("DELETE", base.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (d *DragoAPIClient) addQueries(filter map[string]string, page *Page, req *http.Request) {
	q := req.URL.Query()

	if len(filter) > 0 {
		for k, v := range filter {
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

func (d *DragoAPIClient) addHeaders(req *http.Request) {
	req.Header.Add(contentTypeHeader, contentTypeApplicationJson)

	if d.Token != "" {
		req.Header.Add(defaultTokenHeader, d.Token)
	}
}

func (d *DragoAPIClient) doRequest(receiver interface{}, req *http.Request) error {
	client := &http.Client{Timeout: d.Timeout}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if ok := res.StatusCode >= 200 && res.StatusCode < 300; !ok {
		resBody, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return fmt.Errorf("%v: %v", res.Status, string(resBody))
	}

	if receiver != nil {
		return json.NewDecoder(res.Body).Decode(receiver)
	}

	return nil
}
