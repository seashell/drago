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

type postgresqlHostRepositoryAdapter struct {
	db *sqlx.DB
}

// NewPostgreSQLHostRepositoryAdapter :
func NewPostgreSQLHostRepositoryAdapter(backend Backend) (domain.HostRepository, error) {
	if db, ok := backend.DB().(*sqlx.DB); ok {
		return &postgresqlHostRepositoryAdapter{db}, nil
	}
	return nil, errors.New("Error creating PostgreSQL backend adapter for host repository")
}

func (a *postgresqlHostRepositoryAdapter) GetByID(id string) (*domain.Host, error) {
	sh := &sql.Host{}

	err := a.db.Get(sh,
		"SELECT h.*, array_agg(l.id) AS link_ids "+
			"FROM host h "+
			"LEFT JOIN link l ON l.to_host_id = h.id OR l.from_host_id = h.id "+
			"WHERE id=$1 "+
			"GROUP BY h.id",
		id)
	if err != nil {
		return nil, err
	}

	dh := &domain.Host{}

	errs := model.Copy(dh, sh)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, err
	}

	return dh, nil
}

func (a *postgresqlHostRepositoryAdapter) Create(h *domain.Host) (id *string, err error) {
	guid, err := uuid.NewRandom()
	if err != nil {
		return
	}

	sguid := guid.String()
	now := time.Now()

	err = a.db.QueryRow(
		"INSERT INTO host (id, network_id, name, address, advertise_address, listen_port, public_key, "+
			"table, dns, mtu, pre_up, post_up, pre_down, post_down, created_at, updated_at) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id",
		sguid, *h.NetworkID, *h.Name, *h.IPAddress, *h.AdvertiseAddress, *h.ListenPort, *h.PublicKey,
		*h.Table, *h.DNS, *h.MTU, *h.PreUp, *h.PostUp, *h.PreDown, *h.PostDown, now, now).Scan(id)
	if err != nil {
		return
	}

	return
}

func (a *postgresqlHostRepositoryAdapter) Update(h *domain.Host) (id *string, err error) {
	now := time.Now()

	err = a.db.QueryRow(
		"UPDATE host SET "+
			"name = $1, "+
			"ip_address = $2, "+
			"advertise_address = $3, "+
			"listen_port = $4, "+
			"public_key = $5, "+
			"table = $6, "+
			"dns = $7, "+
			"mtu = $8, "+
			"pre_up = $9, "+
			"post_up = $10, "+
			"pre_down = $11, "+
			"post_down = $12, "+
			"network_id = $13, "+
			"updated_at = $14 "+
			"WHERE id = $15 "+
			"RETURNING id",
		*h.Name, *h.IPAddress, *h.AdvertiseAddress, *h.ListenPort, *h.PublicKey, *h.Table, *h.DNS,
		*h.MTU, *h.PreUp, *h.PostUp, *h.PreDown, *h.PostDown, *h.NetworkID, now, *h.ID).Scan(id)
	if err != nil {
		return
	}

	return
}

func (a *postgresqlHostRepositoryAdapter) DeleteByID(id string) error {
	_, err := a.db.Exec("DELETE FROM host WHERE id = $1", id)
	if err != nil {
		return err
	}

	return err
}

func (a *postgresqlHostRepositoryAdapter) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error) {
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
		"SELECT h.*, COUNT(*) OVER() AS total_count, array_agg(l.id) AS link_ids "+
			"FROM host h "+
			"LEFT JOIN link l ON l.to_host_id = h.id OR l.from_host_id = h.id "+
			"WHERE h.network_id = $1 "+
			"GROUP BY h.id "+
			"ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		id, page.PerPage, page.Page)
	if err != nil {
		return nil, page, err
	}

	receiver := struct {
		sql.Host
		TotalCount int `db:"total_count"`
	}{}

	hostList := []*domain.Host{}

	for rows.Next() {
		err = rows.StructScan(&receiver)
		if err != nil {
			return nil, page, err
		}

		host := &domain.Host{}

		errs := model.Copy(host, receiver)
		if errs != nil {
			for _, e := range errs {
				err = multierror.Append(err, e)
			}
			return nil, page, err
		}

		hostList = append(hostList, host)
	}

	page.TotalCount = receiver.TotalCount
	page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))

	return hostList, page, nil
}
