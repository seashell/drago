package repository

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/adapter/repository/sql"
	"github.com/seashell/drago/server/domain"
	"gopkg.in/jeevatkm/go-model.v1"
)

type postgresqlNetworkRepositoryAdapter struct {
	db *sqlx.DB
}

// NewPostgreSQLNetworkRepositoryAdapter :
func NewPostgreSQLNetworkRepositoryAdapter(backend Backend) (domain.NetworkRepository, error) {
	if db, ok := backend.DB().(*sqlx.DB); ok {
		return &postgresqlNetworkRepositoryAdapter{db}, nil
	}

	return nil, errors.New("Error creating PostgreSQL backend adapter for network repository")
}

// GetByID :
func (a *postgresqlNetworkRepositoryAdapter) GetByID(id string) (*domain.Network, error) {
	sn := &sql.Network{}

	err := a.db.Get(sn, "SELECT * FROM network WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	dn := &domain.Network{}

	errs := model.Copy(dn, sn)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, err
	}

	return dn, nil
}

// Create :
func (a *postgresqlNetworkRepositoryAdapter) Create(n *domain.Network) (*string, error) {
	guid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	sguid := guid.String()
	now := time.Now()

	var id string

	err = a.db.QueryRow(
		"INSERT INTO network (id, name, ip_address_range, created_at, updated_at) "+
			"VALUES ($1, $2, $3, $4, $5) RETURNING id",
		sguid, *n.Name, *n.IPAddressRange, now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// Update :
func (a *postgresqlNetworkRepositoryAdapter) Update(n *domain.Network) (*string, error) {
	now := time.Now()

	var id string

	err := a.db.QueryRow(
		"UPDATE network SET "+
			"name = $1, "+
			"ip_address_range = $2, "+
			"updated_at = $3 "+
			"WHERE id = $4 "+
			"RETURNING id",
		*n.Name, *n.IPAddressRange, now, *n.ID).Scan(id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// DeleteByID :
func (a *postgresqlNetworkRepositoryAdapter) DeleteByID(id string) (*string, error) {
	_, err := a.db.Exec("DELETE FROM network WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// FindAll :
func (a *postgresqlNetworkRepositoryAdapter) FindAll(pageInfo domain.PageInfo) ([]*domain.Network, *domain.Page, error) {
	page := &domain.Page{
		Page:       pageInfo.Page,
		PerPage:    pageInfo.PerPage,
		TotalCount: 0,
		PageCount:  0,
	}

	if page.PerPage > MaxQueryRows {
		page.PerPage = MaxQueryRows
	}

	rows, err := a.db.Queryx("SELECT *, COUNT(*) OVER() AS total_count FROM network ORDER BY created_at DESC LIMIT $1 OFFSET $2", page.PerPage, (page.Page-1)*page.PerPage)
	if err != nil {
		return nil, page, err
	}

	receiver := struct {
		sql.Network
		TotalCount int `db:"total_count"`
	}{}

	networkList := []*domain.Network{}

	for rows.Next() {
		err = rows.StructScan(&receiver)
		if err != nil {
			return nil, page, err
		}

		network := &domain.Network{}

		errs := model.Copy(network, receiver.Network)
		if errs != nil {
			for _, e := range errs {
				err = multierror.Append(err, e)
			}
			return nil, page, err
		}

		networkList = append(networkList, network)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}

	return networkList, page, nil
}
