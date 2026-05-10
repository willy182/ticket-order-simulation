package domain

import "github.com/golangid/candi/candishared"

// FilterTransaction model
type FilterTransaction struct {
	candishared.Filter
	ID       *int64   `json:"id"`
	TicketID *int     `json:"ticket_id,omitempty"`
	Status   string   `json:"status,omitempty"`
	Preloads []string `json:"-"`
}
