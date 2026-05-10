package helper

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Filter data
type Filter struct {
	Limit            int    `json:"limit" default:"10"`
	Page             int    `json:"page" default:"1"`
	Offset           int    `json:"-"`
	Search           string `json:"search,omitempty"`
	OrderBy          string `json:"order_by,omitempty"`
	Sort             string `json:"sort,omitempty" default:"desc" lower:"true"`
	ShowAll          bool   `json:"show_all"`
	AllowEmptyFilter bool   `json:"-"`
}

// CalculateOffset method
func (f *Filter) CalculateOffset() int {
	f.Offset = (f.Page - 1) * f.Limit
	return f.Offset
}

// GetPage method
func (f *Filter) GetPage() int {

	return f.Page
}

// IncrPage method
func (f *Filter) IncrPage() {
	f.Page++
}

// GetLimit method
func (f *Filter) GetLimit() int {
	return f.Limit
}

// Meta model
type Meta struct {
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	TotalRecords int64 `json:"total_records,omitempty"`
	TotalPages   int   `json:"total_pages,omitempty"`
}

// NewMeta create new meta for slice data
func NewMeta(page, limit int, totalRecords int64, showAll bool) (m Meta) {
	m.Page, m.Limit, m.TotalRecords = page, limit, totalRecords
	m.CalculatePages()
	if showAll {
		m.TotalPages = 1
	}
	return m
}

// CalculatePages meta method
func (m *Meta) CalculatePages() {
	m.TotalPages = int(math.Ceil(float64(m.TotalRecords) / float64(m.Limit)))
}

// GenerateInvoice menghasilkan format INV/yymmdd/random
func GenerateInvoice() string {
	now := time.Now()

	// Format Tanggal: yymmdd
	datePart := now.Format("060102")

	// Cara Baru: Membuat generator lokal yang aman dan tidak global
	// Ini menghindari ketergantungan pada rand.Seed() global
	r := rand.New(rand.NewSource(now.UnixNano()))

	// Generate angka random 5 digit (10000 - 99999)
	randomNumber := r.Intn(90000) + 10000

	return fmt.Sprintf("INV/%s/%d", datePart, randomNumber)
}

func GenerateTicketCode(length int) string {
	// Membuat local generator yang unik berdasarkan waktu
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, length)
	for i := range result {
		// Mengambil karakter acak dari CHARSET
		result[i] = CHARSET[r.Intn(len(CHARSET))]
	}

	return string(result)
}
