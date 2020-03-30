package server

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	gomodel "gopkg.in/jeevatkm/go-model.v1"
)

type host struct {
	ID               int    `gorm:"primary_key;auto_increment:true"`
	Name             string `gorm:"unique"`
	Address          string
	AdvertiseAddress string
	ListenPort       string
	PublicKey        string
	Table            string
	DNS              string
	Mtu              string
	PreUp            string
	PostUp           string
	PreDown          string
	PostDown         string
	Links            []*link `gorm:"foreignkey:FromID"`
	LastSeen         time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

type link struct {
	ID                  int `gorm:"primary_key;auto_increment:true"`
	FromID              int `gorm:"type:integer REFERENCES hosts(id) ON DELETE CASCADE ON UPDATE CASCADE; unique_index:idx_directional_link"`
	ToID                int `gorm:"type:integer REFERENCES hosts(id) ON DELETE CASCADE ON UPDATE CASCADE; unique_index:idx_directional_link"`
	From                *host
	To                  *host
	PersistentKeepalive int
	AllowedIPs          string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

type inMemoryStore struct {
	db *gorm.DB
}

func NewInMemoryStore() (*inMemoryStore, error) {

	db, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	db.LogMode(true)

	db.AutoMigrate(&host{}, &link{})

	return &inMemoryStore{
		db: db,
	}, nil
}

func (s *inMemoryStore) CreateHost(h *Host) (*Host, error) {

	nh := &host{}
	gomodel.Copy(nh, h)

	err := s.db.Create(nh).Error
	if err != nil {
		return nil, err
	}

	res := &Host{}
	gomodel.Copy(res, nh)

	return res, err
}

func (s *inMemoryStore) UpdateHost(id int, h *Host) (*Host, error) {

	uh := &host{ID: id}
	err := s.db.First(&uh).Updates(h).Error
	if err != nil {
		return nil, err
	}

	res := &Host{}
	gomodel.Copy(res, uh)

	return res, err
}

func (s *inMemoryStore) DeleteHost(id int) error {
	err := s.db.Delete(&host{ID: id}).Error
	if err != nil {
		return err
	}

	// Force cascade behavior. TODO: improve this
	err = s.db.Where("from_id = ?", id).Or("to_id = ?", id).Delete(&link{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *inMemoryStore) GetAllHosts() ([]*Host, error) {

	hosts := make([]*Host, 0)

	err := s.db.Find(&hosts).Error
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func (s *inMemoryStore) GetHost(id int) (*Host, error) {

	h := host{}
	err := s.db.Joins("LEFT JOIN links ON links.from_id = hosts.id").Preload("Links").Preload("Links.To").Preload("Links.From").Where("hosts.id = ?", id).First(&h).Error
	if err != nil {
		return nil, err
	}

	res := &Host{
		Links: make([]*Link, 0),
	}

	for _, l := range h.Links {
		L := &Link{
			To:   &Host{},
			From: &Host{},
		}
		gomodel.Copy(L, l)
		gomodel.Copy(L.To, l.To)
		gomodel.Copy(L.From, l.From)
		res.Links = append(res.Links, L)
	}

	gomodel.Copy(res, h)

	return res, nil
}

func (s *inMemoryStore) GetHostWithLinks(id int) (*Host, error) {

	h := &host{}
	err := s.db.Where("id = ?", id).Find(&h).Error
	if err != nil {
		return nil, err
	}

	res := &Host{}
	gomodel.Copy(res, h)

	return res, nil
}

func (s *inMemoryStore) CreateLink(l *Link) (*Link, error) {

	nl := link{}
	gomodel.Copy(&nl, l)

	err := s.db.Create(&nl).Error
	if err != nil {
		return nil, err
	}

	res := &Link{}
	gomodel.Copy(res, nl)

	return l, err
}

func (s *inMemoryStore) DeleteLink(id int) error {
	err := s.db.Delete(&link{ID: id}).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *inMemoryStore) UpdateLink(id int, l *Link) (*Link, error) {

	ul := &link{}
	gomodel.Copy(ul, l)

	err := s.db.Save(&ul).Error
	if err != nil {
		return nil, err
	}

	res := &Link{}
	gomodel.Copy(res, ul)

	return res, err
}

func (s *inMemoryStore) GetAllLinks() ([]*Link, error) {

	links := make([]*Link, 0)

	err := s.db.Preload("To").Preload("From").Find(&links).Error
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (host) TableName() string {
	return "hosts"
}

func (link) TableName() string {
	return "links"
}
