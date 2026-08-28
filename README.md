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

> **Penting:** jika `JWT_SECRET` tidak diisi, kode jatuh ke default `secret` yang **tidak aman untuk production** — siapapun bisa memalsukan token JWT. Set `JWT_SECRET` wajib dilakukan di `.env` lokal dan sebagai GitHub secret `JWT_SECRET` (yang otomatis disinkronkan ke Koyeb Secret, lihat tabel secrets). Nilai yang sama dipakai untuk memverifikasi token.

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
| Push ke branch `main`          | Test + migrasi branch **production** + deploy Koyeb |

Alur kerja harian:

```text
fitur → push ke test          → CI migrasi Neon testing, tes manual dari laptop
      → PR test → main        → job test jalan
      → merge ke main         → CI migrasi Neon production + deploy Koyeb
```

App testing tetap dijalankan manual dari local (`APP_ENV=testing go run ./cmd/api`).

### Secrets GitHub yang diperlukan

Di **Settings → Secrets and variables → Actions** — **wajib dibuat**, tanpa ini job migrate/deploy gagal:

| Secret                 | Isi                                              |
| ---------------------- | ------------------------------------------------ |
| `NEON_TESTING_URL`     | Direct URL branch testing                        |
| `NEON_PRODUCTION_URL`  | Direct URL branch production                     |
| `KOYEB_API_TOKEN`      | Personal Access Token Koyeb (**wajib untuk deploy**) |
| `JWT_SECRET`           | Random string kuat untuk memverifikasi token JWT (**wajib untuk deploy**) |

> Migrasi di CI memakai binary [golang-migrate v4.19.1](https://github.com/golang-migrate/migrate/releases) yang diunduh di runner, bukan container Docker.

## Deploy ke Koyeb

Deploy otomatis ke [Koyeb](https://www.koyeb.com) saat **merge ke branch `main`**. Job `deploy-koyeb` membangun image dari `Dockerfile` (builder `docker`) lalu membuat/memperbarui App & Service **sistem-service-motor** dengan env var dari Koyeb Secrets.

> **Catatan storage:** filesystem Koyeb bersifat **ephemeral** — file upload (`uploads/`) hilang saat redeploy/restart. Gunakan object storage (S3/R2) jika butuh file persisten.

### Setup satu kali (manual)

1. **Buat Personal Access Token Koyeb** (Dashboard Koyeb → Account → API/Token). Simpan sebagai GitHub secret `KOYEB_API_TOKEN`.

2. Pastikan GitHub secret berikut sudah dibuat:
   - `NEON_PRODUCTION_URL` — Direct URL branch production Neon
   - `JWT_SECRET` — random string kuat (sama dengan yang dipakai lokal)

3. Jalankan sekali job `deploy-koyeb` (mis. push dummy ke `main`). Koyeb akan otomatis:
   - Membuat Koyeb Secrets `DATABASE_URL` (dari `NEON_PRODUCTION_URL`) dan `JWT_SECRET`.
   - Membuat App `sistem-service-motor` + Service `api` yang membangun image dari Dockerfile, mengekspos port `8080`, dan route `/`.

4. Setelah deploy pertama berhasil, ambil **Public URL** untuk diisi sebagai env/QR code publik:
   - `PUBLIC_BASE_URL` (optional) di env service = Public URL Koyeb (mis. `https://sistem-service-motor-<org>.koyeb.app`).

App tersedia di Public URL yang diberikan Koyeb (domain `.koyeb.app`). Untuk domain custom, buka App → Settings → Domains di dashboard Koyeb.

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
