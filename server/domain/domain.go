package domain

const (
	DefaultPage    = 1
	DefaultPerPage = 100
	MaxPerPage     = 1000
)

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
