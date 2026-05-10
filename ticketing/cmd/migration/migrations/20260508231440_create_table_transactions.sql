-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
	"id" BIGSERIAL NOT NULL PRIMARY KEY,
	"invoice_number" VARCHAR(100) UNIQUE,
	"customer_name" VARCHAR(255) NOT NULL,
	"customer_email" VARCHAR(150) NOT NULL,
	"customer_phone" VARCHAR(15) NOT NULL,
	"ticket_code" VARCHAR(25) UNIQUE,
	"status" VARCHAR(7) NOT NULL,
	"qty" INTEGER NOT NULL,
	"total_amount" NUMERIC(18,2) NOT NULL,
	"ticket_id" INTEGER NOT NULL REFERENCES tickets(id) ON DELETE RESTRICT,
	"created_at" TIMESTAMPTZ(6),
	"updated_at" TIMESTAMPTZ(6)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
