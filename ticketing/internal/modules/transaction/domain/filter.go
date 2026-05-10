package domain

import "github.com/golangid/candi/candishared"

// FilterTransaction model
type FilterTransaction struct {
	candishared.Filter
	ID        *int64   `json:"id"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	Preloads  []string `json:"-"`
}
