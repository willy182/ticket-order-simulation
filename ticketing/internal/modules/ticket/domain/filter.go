package domain

import "ticketing/pkg/helper"

// FilterTicket model
type FilterTicket struct {
	helper.Filter
	ID       *int     `json:"id"`
	Preloads []string `json:"-"`
}
