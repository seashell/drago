package repository

import (
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/adapter/repository/sql"
	"github.com/seashell/drago/server/domain"
	"gopkg.in/jeevatkm/go-model.v1"
)

type postgresqlLinkRepositoryAdapter struct {
	db *sqlx.DB
}

// NewPostgreSQLLinkRepositoryAdapter :
func NewPostgreSQLLinkRepositoryAdapter(backend Backend) (domain.LinkRepository, error) {
	if db, ok := backend.DB().(*sqlx.DB); ok {
		return &postgresqlLinkRepositoryAdapter{db}, nil
	}

	return nil, errors.New("Error creating PostgreSQL backend adapter for link repository")
}

func (a *postgresqlLinkRepositoryAdapter) GetByID(id string) (*domain.Link, error) {
	sl := &sql.Link{}

	err := a.db.Get(sl, "SELECT * FROM link WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	dl := &domain.Link{}

	errs := model.Copy(dl, sl)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, err
	}

	return dl, nil
}

func (a *postgresqlLinkRepositoryAdapter) Create(l *domain.Link) (*string, error) {
	guid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	sguid := guid.String()
	now := time.Now()

	var id string

	strAllowedIPs := strings.Join(l.AllowedIPs[:], ",")

	err = a.db.QueryRow(
		`INSERT INTO link (id, network_id, from_host_id, to_host_id, allowed_ips, persistent_keepalive,
			created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		sguid, l.NetworkID, l.FromHostID, l.ToHostID, strAllowedIPs, l.PersistentKeepalive, now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (a *postgresqlLinkRepositoryAdapter) Update(l *domain.Link) (*string, error) {
	now := time.Now()

	var id string

	strAllowedIPs := strings.Join(l.AllowedIPs[:], ",")

	err := a.db.QueryRow(
		`UPDATE link SET
			allowed_ips = $1,
			persistent_keepalive = $2,
			updated_at = $3
			WHERE id = $4
			RETURNING id`,
		strAllowedIPs, l.PersistentKeepalive, now, l.ID).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (a *postgresqlLinkRepositoryAdapter) DeleteByID(id string) (*string, error) {
	_, err := a.db.Exec("DELETE FROM link WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (a *postgresqlLinkRepositoryAdapter) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	page := &domain.Page{
		Page:       pageInfo.Page,
		PerPage:    pageInfo.PerPage,
		TotalCount: 0,
		PageCount:  0,
	}

	if page.PerPage > maxQueryRows {
		page.PerPage = maxQueryRows
	}
	rows, err := a.db.Queryx(
		"SELECT id, created_at, updated_at, network_id, from_host_id, to_host_id, persistent_keepalive, allowed_ips, "+
			"COUNT(*) OVER() AS total_count "+
			"FROM link "+
			"WHERE network_id = $1 "+
			"ORDER BY created_at DESC LIMIT $2 OFFSET $3", id, page.PerPage, (page.Page-1)*page.PerPage)
	if err != nil {
		return nil, page, err
	}

	receiver := struct {
		sql.Link
		StrAllowedIPs string `db:"allowed_ips"`
		TotalCount    int    `db:"total_count"`
	}{}

	linkList := []*domain.Link{}

	for rows.Next() {
		err = rows.StructScan(&receiver)
		if err != nil {
			return nil, page, err
		}

		link := &domain.Link{}

		errs := model.Copy(link, receiver.Link)
		if errs != nil {
			for _, e := range errs {
				err = multierror.Append(err, e)
			}
			return nil, page, err
		}

		link.AllowedIPs = strings.Fields(strings.Replace(receiver.StrAllowedIPs, ",", " ", -1))

		linkList = append(linkList, link)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return linkList, page, nil
}

func (a *postgresqlLinkRepositoryAdapter) FindAllByHostID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	page := &domain.Page{
		Page:       pageInfo.Page,
		PerPage:    pageInfo.PerPage,
		TotalCount: 0,
		PageCount:  0,
	}

	if page.PerPage > maxQueryRows {
		page.PerPage = maxQueryRows
	}

	rows, err := a.db.Queryx(
		"SELECT id, created_at, updated_at, network_id, from_host_id, to_host_id, persistent_keepalive, allowed_ips, "+
			"COUNT(*) OVER() AS total_count "+
			"FROM link "+
			"WHERE from_host_id = $1 "+
			"ORDER BY created_at DESC LIMIT $2 OFFSET $3", id, page.PerPage, (page.Page-1)*page.PerPage)
	if err != nil {
		return nil, page, err
	}

	receiver := struct {
		sql.Link
		StrAllowedIPs string `db:"allowed_ips"`
		TotalCount    int    `db:"total_count"`
	}{}

	linkList := []*domain.Link{}

	for rows.Next() {
		err = rows.StructScan(&receiver)
		if err != nil {
			return nil, page, err
		}

		link := &domain.Link{}

		errs := model.Copy(link, receiver.Link)
		if errs != nil {
			for _, e := range errs {
				err = multierror.Append(err, e)
			}
			return nil, page, err
		}

		link.AllowedIPs = strings.Fields(strings.Replace(receiver.StrAllowedIPs, ",", " ", -1))

		linkList = append(linkList, link)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return linkList, page, nil
}
