# Noteriety

Note stuff

## Development

### Dependencies

- `sqlite3`
- `go`

### Embedding the SQL Migrations

In order to support the `init` command and not duplicate SQL, this project uses
[go-bindata](https://github.com/go-bindata/go-bindata) in order to embed
`sql/migration/*.sql` files. You will need to run `go generate` on your first
build of the project, and whenever there are any subsequent edits to
migration files. The generated `note/bindata` package will be tracked in source
control.

*Waiting for `go:embed` support.*
