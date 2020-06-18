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

type postgresqlInterfaceRepositoryAdapter struct {
	db *sqlx.DB
}

// NewPostgreSQLInterfaceRepositoryAdapter :
func NewPostgreSQLInterfaceRepositoryAdapter(backend Backend) (domain.InterfaceRepository, error) {
	if db, ok := backend.DB().(*sqlx.DB); ok {
		return &postgresqlInterfaceRepositoryAdapter{db}, nil
	}
	return nil, errors.New("Error creating PostgreSQL backend adapter for interface repository")
}

func (a *postgresqlInterfaceRepositoryAdapter) GetByID(id string) (*domain.Interface, error) {
	sqlOut := &sql.Interface{}
	err := a.db.Get(sqlOut,
		`SELECT iface.* FROM interface iface
			WHERE iface.id=$1
			GROUP BY iface.id`,
		id)
	if err != nil {
		return nil, err
	}

	res := &domain.Interface{}

	errs := model.Copy(res, sqlOut)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, err
	}

	return res, nil
}

func (a *postgresqlInterfaceRepositoryAdapter) Create(i *domain.Interface) (*string, error) {
	guid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	sguid := guid.String()
	now := time.Now()

	var id string

	err = a.db.QueryRow(
		`INSERT INTO interface (
			id,
			host_id,
			network_id,
			name,
			ip_address,
			listen_port,
			public_key,
			"table",
			dns,
			mtu,
			pre_up,
			post_up,
			pre_down,
			post_down,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) 
		RETURNING id`,
		sguid, i.HostID, i.NetworkID, i.Name, i.IPAddress, i.ListenPort, i.PublicKey,
		i.Table, i.DNS, i.MTU, i.PreUp, i.PostUp, i.PreDown, i.PostDown, now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (a *postgresqlInterfaceRepositoryAdapter) Update(h *domain.Interface) (*string, error) {
	now := time.Now()

	var id string

	err := a.db.QueryRow(
		`UPDATE interface SET
			name = $1,
			ip_address = $2,
			listen_port = $3,
			public_key = $4,
			"table" = $5,
			dns = $6,
			mtu = $7,
			pre_up = $8,
			post_up = $9,
			pre_down = $10,
			post_down = $11,
			network_id = $12,
			updated_at = $13
			WHERE id = $14
			RETURNING id`,
		h.Name, h.IPAddress, h.ListenPort, h.PublicKey, h.Table, h.DNS,
		h.MTU, h.PreUp, h.PostUp, h.PreDown, h.PostDown, h.NetworkID, now, h.ID).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (a *postgresqlInterfaceRepositoryAdapter) DeleteByID(id string) (*string, error) {
	_, err := a.db.Exec(`DELETE FROM interface WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (a *postgresqlInterfaceRepositoryAdapter) FindAll(pageInfo domain.PageInfo) ([]*domain.Interface, *domain.Page, error) {
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
		`SELECT iface.*, COUNT(*) OVER() AS total_count
			FROM interface iface 
			GROUP BY iface.id 
			ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		page.PerPage, (page.Page-1)*page.PerPage)
	if err != nil {
		return nil, page, err
	}

	receiver := struct {
		sql.Interface
		TotalCount int `db:"total_count"`
	}{}

	ifaceList := []*domain.Interface{}

	for rows.Next() {
		err = rows.StructScan(&receiver)
		if err != nil {
			return nil, page, err
		}

		iface := &domain.Interface{}

		errs := model.Copy(iface, receiver.Interface)
		if errs != nil {
			for _, e := range errs {
				err = multierror.Append(err, e)
			}
			return nil, page, err
		}

		ifaceList = append(ifaceList, iface)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return ifaceList, page, nil
}

func (a *postgresqlInterfaceRepositoryAdapter) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Interface, *domain.Page, error) {
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
		`SELECT iface.*, COUNT(*) OVER() AS total_count
			FROM interface iface 
			WHERE iface.network_id = $1
			GROUP BY iface.id 
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		id, page.PerPage, (page.Page-1)*page.PerPage)
	if err != nil {
		return nil, page, err
	}

	receiver := struct {
		sql.Interface
		TotalCount int `db:"total_count"`
	}{}

	ifaceList := []*domain.Interface{}

	for rows.Next() {
		err = rows.StructScan(&receiver)
		if err != nil {
			return nil, page, err
		}

		iface := &domain.Interface{}

		errs := model.Copy(iface, receiver.Interface)
		if errs != nil {
			for _, e := range errs {
				err = multierror.Append(err, e)
			}
			return nil, page, err
		}

		ifaceList = append(ifaceList, iface)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return ifaceList, page, nil
}

func (a *postgresqlInterfaceRepositoryAdapter) FindAllByHostID(id string, pageInfo domain.PageInfo) ([]*domain.Interface, *domain.Page, error) {
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
		`SELECT iface.*, COUNT(*) OVER() AS total_count
			FROM interface iface
			WHERE iface.host_id = $1
			GROUP BY iface.id 
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		id, page.PerPage, (page.Page-1)*page.PerPage)
	if err != nil {
		return nil, page, err
	}

	receiver := struct {
		sql.Interface
		TotalCount int `db:"total_count"`
	}{}

	ifaceList := []*domain.Interface{}

	for rows.Next() {
		err = rows.StructScan(&receiver)
		if err != nil {
			return nil, page, err
		}

		iface := &domain.Interface{}

		errs := model.Copy(iface, receiver.Interface)
		if errs != nil {
			for _, e := range errs {
				err = multierror.Append(err, e)
			}
			return nil, page, err
		}

		ifaceList = append(ifaceList, iface)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return ifaceList, page, nil
}
