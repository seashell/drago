package client

import (
	"io/ioutil"
	"os"
)

// Read file contents if they exist, or persist and return a default value otherwise.
func (c *Client) readFileLazy(path string, s string) (string, error) {

	var out string

	buf, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if len(buf) != 0 {
		out = string(buf)
	} else {
		out = s
		if err := ioutil.WriteFile(path, []byte(s), 0700); err != nil {
			return "", err
		}
	}

	return out, nil
}
