package domain

// PageInfo :
type PageInfo struct {
	Page    int
	PerPage int
}

// Page :
type Page struct {
	Page       int
	PerPage    int
	PageCount  int
	TotalCount int
}
