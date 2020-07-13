package repository

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/adapter/repository/sql"
	"github.com/seashell/drago/server/domain"
	"github.com/shurcooL/go-goon"
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

func (a *postgresqlHostRepositoryAdapter) ExistsByID(id string) (bool, error) {
	sql := `SELECT EXISTS (SELECT 1 FROM host h WHERE h.id = $1) AS "exists"`

	res := false

	err := a.db.QueryRow(sql, id).Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}

// GetByID :
func (a *postgresqlHostRepositoryAdapter) GetByID(id string) (*domain.Host, error) {

	query := `SELECT h.* FROM host h WHERE h.id=$1 GROUP BY h.id`

	receiver := struct {
		sql.Host
		StrLabels string `db:"labels"`
	}{}

	err := a.db.Get(&receiver, query, id)
	if err != nil {
		return nil, err
	}

	res := &domain.Host{}

	errs := model.Copy(res, receiver.Host)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, err
	}

	res.Labels = commaSeparatedStrToSlice(receiver.StrLabels)

	return res, nil
}

// Create :
func (a *postgresqlHostRepositoryAdapter) Create(h *domain.Host) (*string, error) {
	guid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	strLabels := strings.Join(h.Labels[:], ",")

	sguid := guid.String()
	now := time.Now()

	var id string

	err = a.db.QueryRow(
		`INSERT INTO host (
			id,
			name,
			advertise_address,
			labels,
			created_at,
			updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		sguid, h.Name, h.AdvertiseAddress, strLabels, now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// CreateWithID :
func (a *postgresqlHostRepositoryAdapter) CreateWithID(h *domain.Host) (*string, error) {

	strLabels := strings.Join(h.Labels[:], ",")

	sguid := h.ID
	now := time.Now()

	if sguid == nil {
		return nil, errors.New("host id is nil")
	}

	if _, err := uuid.Parse(*sguid); err != nil {
		return nil, err
	}

	var id string

	err := a.db.QueryRow(
		`INSERT INTO host (
			id,
			name,
			advertise_address,
			labels,
			created_at,
			updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		sguid, h.Name, h.AdvertiseAddress, strLabels, now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// Update :
func (a *postgresqlHostRepositoryAdapter) Update(h *domain.Host) (*string, error) {
	now := time.Now()

	strLabels := strings.Join(h.Labels[:], ",")

	var id string

	err := a.db.QueryRow(
		`UPDATE host SET
			name = $1,
			labels = $2,
			advertise_address = $3,
			updated_at = $4
			WHERE id = $5
			RETURNING id`,
		h.Name, strLabels, h.AdvertiseAddress, now, h.ID).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// DeleteByID :
func (a *postgresqlHostRepositoryAdapter) DeleteByID(id string) (*string, error) {
	_, err := a.db.Exec("DELETE FROM host WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// FindAll :
func (a *postgresqlHostRepositoryAdapter) FindAll(pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error) {
	page := &domain.Page{
		Page:       pageInfo.Page,
		PerPage:    pageInfo.PerPage,
		TotalCount: 0,
		PageCount:  0,
	}

	if page.PerPage > MaxQueryRows {
		page.PerPage = MaxQueryRows
	}

	rows, err := a.db.Queryx(
		`SELECT h.*, COUNT(*) OVER() AS total_count FROM host h
		GROUP BY h.id 
		ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		page.PerPage, (page.Page-1)*page.PerPage)
	if err != nil {
		return nil, page, err
	}

	receiver := struct {
		sql.Host
		StrLabels  string `db:"labels"`
		TotalCount int    `db:"total_count"`
	}{}

	hostList := []*domain.Host{}

	for rows.Next() {
		err = rows.StructScan(&receiver)
		if err != nil {
			return nil, page, err
		}

		host := &domain.Host{}

		errs := model.Copy(host, receiver.Host)
		if errs != nil {
			for _, e := range errs {
				err = multierror.Append(err, e)
			}
			return nil, page, err
		}

		host.Labels = commaSeparatedStrToSlice(receiver.StrLabels)

		hostList = append(hostList, host)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return hostList, page, nil
}

// FindAllByNetworkID :
func (a *postgresqlHostRepositoryAdapter) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error) {
	page := &domain.Page{
		Page:       pageInfo.Page,
		PerPage:    pageInfo.PerPage,
		TotalCount: 0,
		PageCount:  0,
	}

	if page.PerPage > MaxQueryRows {
		page.PerPage = MaxQueryRows
	}

	rows, err := a.db.Queryx(
		`SELECT h.*, COUNT(*) OVER() AS total_count 
		FROM host h
		LEFT JOIN interface if ON if.host_id = h.id
		LEFT JOIN network net ON if.network_id = net.id
		WHERE if.network_id = $1 
		GROUP BY h.id 
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		id, page.PerPage, (page.Page-1)*page.PerPage)
	if err != nil {
		return nil, page, err
	}

	receiver := struct {
		sql.Host
		StrLabels  string `db:"labels"`
		TotalCount int    `db:"total_count"`
	}{}

	hostList := []*domain.Host{}

	for rows.Next() {
		err = rows.StructScan(&receiver)
		if err != nil {
			return nil, page, err
		}

		host := &domain.Host{}

		errs := model.Copy(host, receiver.Host)
		if errs != nil {
			for _, e := range errs {
				err = multierror.Append(err, e)
			}
			return nil, page, err
		}

		host.Labels = commaSeparatedStrToSlice(receiver.StrLabels)
		hostList = append(hostList, host)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return hostList, page, nil
}

// FindAllByLabels :
func (a *postgresqlHostRepositoryAdapter) FindAllByLabels(labels []string, pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error) {
	page := &domain.Page{
		Page:       pageInfo.Page,
		PerPage:    pageInfo.PerPage,
		TotalCount: 0,
		PageCount:  0,
	}

	if page.PerPage > MaxQueryRows {
		page.PerPage = MaxQueryRows
	}

	for i, _ := range labels {
		labels[i] = "'" + labels[i] + "'"
	}

	slabels := strings.Join(labels, ",")

	query := fmt.Sprintf(`SELECT h.*, COUNT(*) OVER() AS total_count 
							FROM host h
							WHERE string_to_array(h.labels,',') @> array[%s]
							GROUP BY h.id
							ORDER BY created_at DESC LIMIT $1 OFFSET $2`, slabels)

	goon.Dump(query)

	rows, err := a.db.Queryx(query, page.PerPage, (page.Page-1)*page.PerPage)
	if err != nil {
		return nil, page, err
	}

	receiver := struct {
		sql.Host
		StrLabels  string `db:"labels"`
		TotalCount int    `db:"total_count"`
	}{}

	hostList := []*domain.Host{}

	for rows.Next() {
		err = rows.StructScan(&receiver)
		if err != nil {
			return nil, page, err
		}

		host := &domain.Host{}

		errs := model.Copy(host, receiver.Host)
		if errs != nil {
			for _, e := range errs {
				err = multierror.Append(err, e)
			}
			return nil, page, err
		}

		host.Labels = commaSeparatedStrToSlice(receiver.StrLabels)
		hostList = append(hostList, host)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return hostList, page, nil
}
