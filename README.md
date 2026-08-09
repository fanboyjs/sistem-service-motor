# Go Gin PostgreSQL REST API

Starter template for a REST API using Go, Gin, PostgreSQL and pgx.

## Architecture

Request -> Handler -> Service -> Repository -> PostgreSQL

## Run

1. Copy `.env.example` to `.env`.
2. Start PostgreSQL:

```bash
docker compose up -d
```

3. Run migrations:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/my_api?sslmode=disable" up
```

4. Start API:

```bash
go mod tidy
go run ./cmd/api
```

API:

`GET /api/v1/users`
