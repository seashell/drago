package state

import (
	"encoding/json"
	"github.com/seashell/drago/api"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	dbFileName = "state.json"
)

type FileDB struct {
	fileName string
}

func NewFileDB(stateDir string) (StateDB, error) {
	fn := filepath.Join(stateDir, dbFileName)

	// Check to see if the DB already exists
	fi, err := os.Stat(fn)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	firstRun := fi == nil

	// If is firstRun, create empty file
	if firstRun {
		if err := ioutil.WriteFile(fn, []byte("{}"), 0644); err != nil {
			return nil, err
		}
	}

	return &FileDB{
		fileName: fn,
	}, nil

}

func (f *FileDB) Name() string {
	return "filedb"
}

func (f *FileDB) GetHostSettings() (*api.HostSettings, error) {

	fc, err := ioutil.ReadFile(f.fileName)
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

func (f *FileDB) PutHostSettings(hs *api.HostSettings) error {

	fc, err := json.MarshalIndent(hs, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.fileName, fc, 0644)
	if err != nil {
		return err
	}

	return nil
}
