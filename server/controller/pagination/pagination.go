package pagination

const (
	PAGINATION_DEFAULT_PAGE     = 1
	PAGINATION_DEFAULT_PER_PAGE = 10
)

// Input : Query parameters for pagination data
type Input struct {
	Page    int `query:"page" validate:"omitempty,numeric"`
	PerPage int `query:"perPage" validate:"omitempty,numeric,max=100"`
}

// Page : Response body for a paginated list request
type Page struct {
	Page       int         `json:"page"`
	PerPage    int         `json:"perPage"`
	PageCount  int         `json:"pageCount"`
	TotalCount int         `json:"totalCount"`
	Items      interface{} `json:"items"`
}
