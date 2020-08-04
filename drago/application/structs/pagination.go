package structs

const (
	PaginationDefaultPage    = 1
	PaginationDefaultPerPage = 10
)

// PageInput :
type PageInput struct {
	Page    int `query:"page" validate:"omitempty,numeric"`
	PerPage int `query:"perPage" validate:"omitempty,numeric,max=100"`
}

// PageOutput :
type PageOutput struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	PageCount  int `json:"pageCount"`
	TotalCount int `json:"totalCount"`
}
