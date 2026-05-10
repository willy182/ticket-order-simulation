package domain

import (
	"ticketing/pkg/helper"
	shareddomain "ticketing/pkg/shared/domain"
	"time"
)

// ResponseTicketList model
type ResponseTicketList struct {
	Meta helper.Meta      `json:"meta"`
	Data []ResponseTicket `json:"data"`
}

// ResponseTicket model
type ResponseTicket struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Quota     int     `json:"quota"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// Serialize from db model
func (r *ResponseTicket) Serialize(source *shareddomain.Ticket) {
	r.ID = source.ID
	r.Title = source.Title
	r.Quota = source.Quota
	r.Price = source.Price
	r.CreatedAt = source.CreatedAt.Format(time.RFC3339)
	r.UpdatedAt = source.UpdatedAt.Format(time.RFC3339)
}
