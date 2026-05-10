package domain

import (
	"time"
)

// Transaction model
type Transaction struct {
	ID            int64     `gorm:"column:id;primary_key" json:"id"`
	InvoiceNumber *string   `gorm:"column:invoice_number;type:varchar(100)" json:"invoice_number"`
	CustomerName  string    `gorm:"column:customer_name;type:varchar(255)" json:"customer_name"`
	CustomerEmail string    `gorm:"column:customer_email;type:varchar(150)" json:"customer_email"`
	CustomerPhone string    `gorm:"column:customer_phone;type:varchar(15)" json:"customer_phone"`
	TicketCode    *string   `gorm:"column:ticket_code;type:varchar(25)" json:"ticket_code"`
	Status        string    `gorm:"column:status;type:varchar(7)" json:"status"`
	Qty           int       `gorm:"column:qty;type:integer" json:"qty"`
	TotalAmount   float64   `gorm:"column:total_amount;type:numeric(18,2)" json:"total_amount"`
	TicketID      int       `gorm:"column:ticket_id" json:"ticket_id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
	TicketData    Ticket    `gorm:"foreignKey:TicketID"`
}

// TableName return table name of Transaction model
func (Transaction) TableName() string {
	return "transactions"
}
