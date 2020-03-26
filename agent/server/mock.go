package server

import (
	"encoding/json"
	"io/ioutil"

	gomodel "gopkg.in/jeevatkm/go-model.v1"
)

type MockData struct {
	Hosts []HostSummary `json:"hosts"`
	Links []LinkSummary `json:"links"`
}

func PopulateRepositoryWithMockData(repo Repository, path string) error {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	data := MockData{}

	_ = json.Unmarshal([]byte(file), &data)

	for i := 0; i < len(data.Hosts); i++ {
		h := &Host{}
		gomodel.Copy(h, data.Hosts[i])
		repo.CreateHost(h)
	}

	for i := 0; i < len(data.Links); i++ {
		l := &Link{}
		gomodel.Copy(l, data.Links[i])
		repo.CreateLink(l)
	}

	return nil
}
