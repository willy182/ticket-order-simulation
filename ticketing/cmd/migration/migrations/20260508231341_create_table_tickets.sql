-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tickets (
	"id" SERIAL NOT NULL PRIMARY KEY,
	"title" VARCHAR(255) NOT NULL,
	"quota" INTEGER NOT NULL,
	"price" NUMERIC(18,2) NOT NULL,
	"created_at" TIMESTAMPTZ(6),
	"updated_at" TIMESTAMPTZ(6)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tickets;
-- +goose StatementEnd
