package domain

import (
	shareddomain "ticketing/pkg/shared/domain"
)

// RequestTicket model
type RequestTicket struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Quota int     `json:"quota"`
	Price float64 `json:"price"`
}

// Deserialize to db model
func (r *RequestTicket) Deserialize() (res shareddomain.Ticket) {
	res.Title = r.Title
	res.Quota = r.Quota
	res.Price = r.Price
	return
}
