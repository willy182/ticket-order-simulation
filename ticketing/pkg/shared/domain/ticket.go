package domain

import (
	"time"
)

// Ticket model
type Ticket struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Title     string    `gorm:"column:title;type:varchar(255)" json:"title"`
	Quota     int       `gorm:"column:quota;type:integer" json:"quota"`
	Price     float64   `gorm:"column:price;type:numeric(18,2)" json:"price"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName return table name of Ticket model
func (Ticket) TableName() string {
	return "tickets"
}
