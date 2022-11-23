# Migrations

Provides migrations using [goose](https://github.com/pressly/goose) which is a golang
SQL migration tool. All migrations are created incrementally and stored in the `sql`
folder.

Each new migration must have an `up` and `down` section indicated by comments above
the sql `-- +goose Up` and `-- +goose Down`; these control what sql is run for the
`up` and `down` commands respectively.

Goose creates and manages a table storing migration information in a table `goose_db_version`.
