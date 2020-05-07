package pagination

// Input : Query parameters for pagination data
type Input struct {
	Page    int `query:"page" validate:"numeric"`
	PerPage int `query:"perPage" validate:"numeric,max=100"`
}

// Page : Response body for a paginated list request
type Page struct {
	Page       int         `json:"page"`
	PerPage    int         `json:"perPage"`
	PageCount  int         `json:"pageCount"`
	TotalCount int         `json:"totalCount"`
	Items      interface{} `json:"items"`
}
