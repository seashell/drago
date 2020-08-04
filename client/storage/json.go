package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/seashell/drago/api"
)

type store struct {
	filename string
}

// NewJsonStorage :
func NewJsonStorage(path string) (StateRpository, error) {

	// Check to see if the DB already exists
	info, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if info == nil {
		if err := ioutil.WriteFile(path, []byte("{}"), 0644); err != nil {
			return nil, err
		}
	}

	return &store{
		filename: path,
	}, nil

}

// Name :
func (f *store) Name() string {
	return "filedb"
}

// GetHostSettings :
func (s *store) GetHostSettings() (*api.HostSettings, error) {

	fc, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}
	hs := api.HostSettings{}
	err = json.Unmarshal([]byte(fc), &hs)
	if err != nil {
		return nil, err
	}

	return &hs, nil
}

// PutHostSettings :
func (s *store) PutHostSettings(hs *api.HostSettings) error {

	fc, err := json.MarshalIndent(hs, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.filename, fc, 0644)
	if err != nil {
		return err
	}

	return nil
}
