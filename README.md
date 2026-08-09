# Sistem Service Motor

Backend sistem service motor.

## Architecture

- Golang Gin
- Postgresql (Neon)

## Database (Neon)

Project ini memakai [Neon](https://neon.tech) dengan dua branch:

| Branch     | File env          | Cara menjalankan                              |
| ---------- | ----------------- | --------------------------------------------- |
| Production | `.env`            | `go run ./cmd/api`                            |
| Testing    | `.env.testing`    | `APP_ENV=testing go run ./cmd/api`            |

App akan me-load `.env` lalu menimpanya dengan `.env.<APP_ENV>` jika file tersebut ada.

### Mendapatkan connection string

1. Buka [Neon Console](https://console.neon.tech) dan pilih project.
2. Klik **Connect** di Project Dashboard.
3. Pilih **Branch**, **Database**, dan **Role**, lalu matikan toggle **Connection pooling** untuk mendapatkan URL **direct**.
4. Salin connection string ke `.env` (production) atau `.env.testing` (testing).

> Gunakan URL **direct** (tanpa `-pooler`), bukan pooled. Migrasi dan pgxpool sama-sama butuh direct connection.

## Run

1. Copy `.env.example` ke `.env`:

```bash
cp .env.example .env
```

2. Isi `DATABASE_URL` di `.env` dengan connection string Neon branch production (lihat tabel di atas).

3. Jalankan migrasi (pakai direct URL):

```bash
migrate -path migrations -database "postgresql://<user>:<password>@ep-xxx.c-3.ap-southeast-1.aws.neon.tech/<dbname>?sslmode=require" up
```

4. Start API:

```bash
go mod tidy
go run ./cmd/api
```

Untuk environment testing:

```bash
APP_ENV=testing go run ./cmd/api
```

> Windows (cmd/PowerShell): `$env:APP_ENV="testing"; go run ./cmd/api`

## PostgreSQL lokal (opsional)

Untuk development tanpa Neon, gunakan docker compose:

```bash
docker compose up -d
```

Lalu isi `DATABASE_URL` di `.env` dengan URL postgres lokal, contoh:

```bash
migrate -path migrations -database "postgres://postgres:root@localhost:5432/sistem_service_motor?sslmode=disable" up
```

API:

`GET /api/users`
