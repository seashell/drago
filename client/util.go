package client

import (
	"io/ioutil"
	"os"
	"reflect"

	"github.com/seashell/drago/drago/structs"
)

type Diff struct {
	created   []string
	deleted   []string
	updated   []string
	unchanged []string
}

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

func interfacesDiff(old, new map[string]*structs.Interface) Diff {

	diff := Diff{}

	for _, vold := range old {
		if _, ok := new[vold.ID]; !ok {
			diff.deleted = append(diff.deleted, vold.ID)
		}
	}

	for _, vnew := range new {
		if vold, ok := old[vnew.ID]; !ok {
			diff.created = append(diff.created, vnew.ID)
		} else {
			vnew = vold.Merge(vnew)
			if !reflect.DeepEqual(vold, vnew) {
				diff.updated = append(diff.updated, vnew.ID)
			} else {
				diff.unchanged = append(diff.unchanged, vold.ID)
			}
		}
	}

	return diff
}
