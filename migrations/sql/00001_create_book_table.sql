-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated = (NOW() AT TIME ZONE 'utc');
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE EXTENSION "uuid-ossp";
CREATE TABLE IF NOT EXISTS books (
    "id" UUID DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    "title" TEXT NOT NULL DEFAULT '',
    "year" VARCHAR(4),
    "tags" TEXT[],
    "details" JSON NOT NULL,
    "created" TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    "updated" TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);

CREATE INDEX IF NOT EXISTS books_year_idx ON books("year", "created" DESC NULLS LAST);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON books
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX books_year_idx;
DROP TABLE books;
DROP EXTENSION "uuid-ossp";
-- +goose StatementEnd
