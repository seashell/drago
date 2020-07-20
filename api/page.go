package api

const (
	pageQueryParamKey    = "page"
	perPageQueryParamKey = "perPage"
)

type Page struct {
	Page       *int `json:"page,omitempty"`
	PerPage    *int `json:"perPage,omitempty"`
	PageCount  *int `json:"pageCount,omitempty"`
	TotalCount *int `json:"totalCount,omitempty"`
}
