# Sistem Service Motor

Backend sistem service motor.

## Architecture

- Golang Gin
- Postgresql

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

`GET /api/users`
