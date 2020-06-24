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

	receiver := struct {
		sql.Link
		StrAllowedIPs string `db:"allowed_ips"`
	}{}

	err := a.db.Get(&receiver, "SELECT * FROM link WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	res := &domain.Link{}

	errs := model.Copy(res, receiver.Link)
	if errs != nil {
		for _, e := range errs {
			err = multierror.Append(err, e)
		}
		return nil, err
	}

	res.AllowedIPs = commaSeparatedStrToSlice(receiver.StrAllowedIPs)

	return res, nil
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
		`INSERT INTO link (id, from_interface_id, to_interface_id, allowed_ips, persistent_keepalive,
			created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		sguid, l.FromInterfaceID, l.ToInterfaceID, strAllowedIPs, l.PersistentKeepalive, now, now).Scan(&id)
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

func (a *postgresqlLinkRepositoryAdapter) FindAll(pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
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
		`SELECT l.*, COUNT(*) OVER() AS total_count 
			FROM link l LEFT JOIN interface if ON l.to_interface_id = if.id
			ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		page.PerPage, (page.Page-1)*page.PerPage)
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

		link.AllowedIPs = commaSeparatedStrToSlice(receiver.StrAllowedIPs)

		linkList = append(linkList, link)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return linkList, page, nil
}

func (a *postgresqlLinkRepositoryAdapter) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
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
		`SELECT l.*, COUNT(*) OVER() AS total_count 
			FROM link l 
			WHERE network_id = $1 
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`, id, page.PerPage, (page.Page-1)*page.PerPage)
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

		link.AllowedIPs = commaSeparatedStrToSlice(receiver.StrAllowedIPs)

		linkList = append(linkList, link)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return linkList, page, nil
}

func (a *postgresqlLinkRepositoryAdapter) FindAllBySourceHostID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
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
		`SELECT l.*, COUNT(*) OVER() AS total_count 
			FROM link l LEFT JOIN interface if ON l.from_interface_id = if.id
			WHERE if.host_id = $1 
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		id, page.PerPage, (page.Page-1)*page.PerPage)
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

		link.AllowedIPs = commaSeparatedStrToSlice(receiver.StrAllowedIPs)

		linkList = append(linkList, link)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return linkList, page, nil
}

func (a *postgresqlLinkRepositoryAdapter) FindAllByTargetHostID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
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
		`SELECT l.*, COUNT(*) OVER() AS total_count 
			FROM link l LEFT JOIN interface if ON l.to_interface_id = if.id
			WHERE if.host_id = $1 
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		id, page.PerPage, (page.Page-1)*page.PerPage)
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

		link.AllowedIPs = commaSeparatedStrToSlice(receiver.StrAllowedIPs)

		linkList = append(linkList, link)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return linkList, page, nil
}

func (a *postgresqlLinkRepositoryAdapter) FindAllBySourceInterfaceID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
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
		`SELECT l.*, COUNT(*) OVER() AS total_count 
			FROM link l LEFT JOIN interface if ON l.from_interface_id = if.id
			WHERE if.id = $1 
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		id, page.PerPage, (page.Page-1)*page.PerPage)
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

		link.AllowedIPs = commaSeparatedStrToSlice(receiver.StrAllowedIPs)

		linkList = append(linkList, link)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return linkList, page, nil
}

func (a *postgresqlLinkRepositoryAdapter) FindAllByTargetInterfaceID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
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
		`SELECT l.*, COUNT(*) OVER() AS total_count 
			FROM link l LEFT JOIN interface if ON l.to_interface_id = if.id
			WHERE if.id = $1 
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		id, page.PerPage, (page.Page-1)*page.PerPage)
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

		link.AllowedIPs = commaSeparatedStrToSlice(receiver.StrAllowedIPs)

		linkList = append(linkList, link)
	}

	page.TotalCount = receiver.TotalCount
	if page.TotalCount > 0 {
		page.PageCount = int(math.Ceil(float64(page.TotalCount) / float64(page.PerPage)))
	}
	return linkList, page, nil
}
