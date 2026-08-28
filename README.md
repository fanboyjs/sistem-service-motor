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

3. Isi `JWT_SECRET` di `.env` dengan random string kuat (contoh: `openssl rand -hex 32`).

> **Penting:** jika `JWT_SECRET` tidak diisi, kode jatuh ke default `secret` yang **tidak aman untuk production** — siapapun bisa memalsukan token JWT. Set `JWT_SECRET` wajib dilakukan di `.env` lokal dan sebagai GitHub secret `JWT_SECRET` (yang diset sebagai variabel secret di Railway, lihat tabel secrets). Nilai yang sama dipakai untuk memverifikasi token.

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

### Hot reload (Air)

Untuk development, install [Air](https://github.com/air-verse/air) lalu jalankan:

```bash
go install github.com/air-verse/air@latest
```

```bash
air                       # production (memakai .env)
APP_ENV=testing air       # testing (memakai .env.testing)
```

> Windows (cmd/PowerShell): `$env:APP_ENV="testing"; air`
>
> Jika `air` tidak ditemukan, pastikan `$(go env GOPATH)/bin` ada di PATH.

Air otomatis me-restart server saat file `.go` atau `.env` berubah. Konfigurasi di `.air.toml`.

### Migrasi per branch

Jalankan migrasi dengan script `scripts/migrate.sh` (membaca `DATABASE_URL` dari file env yang sesuai):

```bash
./scripts/migrate.sh testing        # migrate up ke branch testing
./scripts/migrate.sh production     # migrate up ke branch production
./scripts/migrate.sh testing down   # rollback 1 langkah
```

Atau manual dengan direct URL:

```bash
migrate -path migrations -database "postgresql://<user>:<password>@ep-xxx.c-3.ap-southeast-1.aws.neon.tech/<dbname>?sslmode=require" up
```

## CI/CD (GitHub Actions)

Alur otomatis di `.github/workflows/ci.yml` — cabang GitHub memetakan branch Neon:

| Event                          | Aksi                                             |
| ------------------------------ | ------------------------------------------------ |
| PR ke `main` / `test`          | Test (vet + build)                               |
| Push ke branch `test`          | Test + migrasi ke branch **testing** Neon        |
| Push ke branch `main`          | Test + migrasi branch **production** + deploy Railway |

Alur kerja harian:

```text
fitur → push ke test          → CI migrasi Neon testing, tes manual dari laptop
      → PR test → main        → job test jalan
      → merge ke main         → CI migrasi Neon production + deploy Railway
```

App testing tetap dijalankan manual dari local (`APP_ENV=testing go run ./cmd/api`).

### Secrets GitHub yang diperlukan

Di **Settings → Secrets and variables → Actions** — **wajib dibuat**, tanpa ini job migrate/deploy gagal:

| Secret                 | Isi                                              |
| ---------------------- | ------------------------------------------------ |
| `NEON_TESTING_URL`     | Direct URL branch testing                        |
| `NEON_PRODUCTION_URL`  | Direct URL branch production                     |
| `RAILWAY_TOKEN`        | Public API token Railway (**wajib untuk deploy**) |
| `JWT_SECRET`           | Random string kuat untuk memverifikasi token JWT (**wajib untuk deploy**) |

> Migrasi di CI memakai binary [golang-migrate v4.19.1](https://github.com/golang-migrate/migrate/releases) yang diunduh di runner, bukan container Docker.

## Deploy ke Railway

Deploy otomatis ke [Railway](https://railway.app) saat **merge ke branch `main`**. Job `deploy-railway` menjalankan Railway CLI (`railway up --service=api`) dari CI, membangun image dari `Dockerfile` (sesuai `railway.json`) lalu men-deploy ke service **api** pada project Railway yang sudah ter-link. Config build & healthcheck ada di `railway.json`.

> **Catatan storage:** filesystem Railway bersifat **ephemeral** — file upload (`uploads/`) hilang saat redeploy/restart. Gunakan object storage (S3/R2) jika butuh file persisten.

### Setup satu kali (manual)

1. **Buat Project di Railway** — tidak perlu menambahkan plugin Postgres karena database tetap pakai Neon.

2. **Tambahkan service dari GitHub repo** `fandipratamaa/sistem-service-motor` dengan builder **Dockerfile** (nama service: `api`).

3. **Link project Railway dari lokal** agar `railway up` tahu target-nya:
   ```bash
   railway link
   ```
   (pilih project & service yang sudah dibuat di atas).

4. **Set env var di service `api`** pada dashboard Railway:
   - `APP_ENV=production`
   - `APP_PORT=8080`
   - `JWT_EXPIRY=24h`
   - `UPLOAD_DIR=uploads`
   - `PUBLIC_BASE_URL` = Public URL Railway (set setelah deploy pertama)
   - `DATABASE_URL` = `NEON_PRODUCTION_URL` (variabel secret)
   - `JWT_SECRET` = random string kuat (variabel secret)
   - `PORT` otomatis disuntikkan Railway.

5. **Buat Public API token Railway** (Dashboard → Account → Tokens). Simpan sebagai GitHub secret `RAILWAY_TOKEN`.

6. Pastikan GitHub secret berikut sudah dibuat:
   - `NEON_PRODUCTION_URL` — Direct URL branch production Neon
   - `JWT_SECRET` — random string kuat (sama dengan yang dipakai lokal)
   - `RAILWAY_TOKEN` — public API token Railway

7. Jalankan sekali job `deploy-railway` (mis. push dummy ke `main`). Railway akan membangun image dari Dockerfile dan deploy service `api`.

8. Setelah deploy pertama berhasil, ambil **Public URL** Railway dan isi `PUBLIC_BASE_URL` di env service.

App tersedia di Public URL yang diberikan Railway (domain `*.up.railway.app`). Untuk domain custom, buka Project → Service → Settings → Domains di dashboard Railway.

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

- `POST /api/register` — body `{ "name", "email", "password" }` → `{ "message", "data": { "token", "email" } }`
- `POST /api/login` — body `{ "email", "password" }` → `{ "message", "data": { "token" } }`
- `GET /api/user-info`
