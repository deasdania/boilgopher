-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION "uuid-ossp";
CREATE TABLE IF NOT EXISTS book (
    "id" UUID DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    "title" TEXT NOT NULL DEFAULT '',
    "created" TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    "updated" TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
)


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE book;
DROP EXTENSION "uuid-ossp";
-- +goose StatementEnd
