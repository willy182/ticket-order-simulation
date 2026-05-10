package domain

import (
	shareddomain "ticketing/pkg/shared/domain"
	"time"

	"github.com/golangid/candi/candihelper"
	"github.com/golangid/candi/candishared"
)

// ResponseTransactionList model
type ResponseTransactionList struct {
	Meta candishared.Meta      `json:"meta"`
	Data []ResponseTransaction `json:"data"`
}

// ResponseTransaction model
type ResponseTransaction struct {
	ID            int64   `json:"id"`
	InvoiceNumber string  `json:"invoice_number"`
	CustomerName  string  `json:"customer_name"`
	CustomerEmail string  `json:"customer_email"`
	CustomerPhone string  `json:"customer_phone"`
	TicketCode    string  `json:"ticket_code"`
	Status        string  `json:"status"`
	Qty           int     `json:"qty"`
	TotalAmount   float64 `json:"total_amount"`
	TicketID      int     `json:"ticket_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// Serialize from db model
func (r *ResponseTransaction) Serialize(source *shareddomain.Transaction) {
	r.ID = source.ID
	r.InvoiceNumber = candihelper.PtrToString(source.InvoiceNumber)
	r.CustomerName = source.CustomerName
	r.CustomerEmail = source.CustomerEmail
	r.CustomerPhone = source.CustomerPhone
	r.TicketCode = candihelper.PtrToString(source.TicketCode)
	r.Status = source.Status
	r.Qty = source.Qty
	r.TotalAmount = source.TotalAmount
	r.TicketID = source.TicketID
	r.CreatedAt = source.CreatedAt.Format(time.RFC3339)
	r.UpdatedAt = source.UpdatedAt.Format(time.RFC3339)
}
