# opnlab
Experimental Go API server for homelab

Make sure we are feature complete between API and dashboard

## Setup

- Air: https://github.com/cosmtrek/air
    - `go install github.com/cosmtrek/air@latest`
- Migrate (if you are running the database migrations outside starting the server): https://github.com/golang-migrate/migrate
    - `go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
- Run `go mod tidy` in root folder to install dependencies

## Running

Start server by running `air` in base folder. Server will hot-reload on source changes (not database changes).

## Database Manual Migrations

- Up: `migrate -source file://migrations -database sqlite://sqlite.db up`
- Down: `migrate -source file://migrations -database sqlite://sqlite.db down`
