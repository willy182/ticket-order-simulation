package domain

import (
	shareddomain "ticketing/pkg/shared/domain"
)

// RequestTransaction model
type RequestTransaction struct {
	ID            int64   `json:"id"`
	CustomerName  string  `json:"customer_name"`
	CustomerEmail string  `json:"customer_email"`
	CustomerPhone string  `json:"customer_phone"`
	Status        string  `json:"status"`
	Qty           int     `json:"qty"`
	TicketID      int     `json:"ticket_id"`
	InvoiceNumber *string `json:"-"`
	TotalAmount   float64
}

// Deserialize to db model
func (r *RequestTransaction) Deserialize() (res shareddomain.Transaction) {
	res.CustomerName = r.CustomerName
	res.CustomerEmail = r.CustomerEmail
	res.CustomerPhone = r.CustomerPhone
	res.Status = r.Status
	res.Qty = r.Qty
	res.TicketID = r.TicketID
	if r.InvoiceNumber != nil {
		res.InvoiceNumber = r.InvoiceNumber
	}
	res.TotalAmount = r.TotalAmount
	return
}

// ReqSendEmail model
type ReqSendEmail struct {
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	Status        string `json:"status"`
	TicketTitle   string `json:"ticket_title"`
}
