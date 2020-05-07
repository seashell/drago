package domain

import "time"

// BaseModel :
type BaseModel struct {
	ID        *string    `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

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
