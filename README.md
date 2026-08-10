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
| Push ke branch `main`          | Test + migrasi branch **production** + deploy VPS |

Alur kerja harian:

```text
fitur → push ke test          → CI migrasi Neon testing, tes manual dari laptop
      → PR test → main        → job test jalan
      → merge ke main         → CI migrasi Neon production + deploy VPS
```

App testing tetap dijalankan manual dari local (`APP_ENV=testing go run ./cmd/api`).

### Secrets GitHub yang diperlukan

Di **Settings → Secrets and variables → Actions** — **wajib dibuat**, tanpa ini job migrate/deploy gagal:

| Secret                 | Isi                                              |
| ---------------------- | ------------------------------------------------ |
| `NEON_TESTING_URL`     | Direct URL branch testing                        |
| `NEON_PRODUCTION_URL`  | Direct URL branch production                     |
| `VPS_HOST`             | IP VPS Rumahweb (**wajib untuk deploy**)         |
| `VPS_USER`             | User deploy di VPS (**wajib untuk deploy**)      |
| `VPS_SSH_KEY`          | Private key SSH (**wajib untuk deploy**)         |
| `VPS_DATABASE_URL`     | Direct URL production (**wajib untuk deploy**)   |

> Migrasi di CI memakai binary [golang-migrate v4.19.1](https://github.com/golang-migrate/migrate/releases) yang diunduh di runner, bukan container Docker.

> **Deploy VPS otomatis di-skip** sampai keempat secret VPS di-set. Begitu VPS siap dan secret dibuat, job `deploy` langsung aktif tanpa ubah kode.

## Deploy ke VPS (cth:Rumahweb)

Setup satu kali (manual):

1. **Folder app** di VPS:

```bash
mkdir -p /home/<user>/apps/sistem-service-motor
```

2. **SSH key** — generate di lokal:

```bash
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519_deploy -C "deploy"
```

Salin pubkey (`~/.ssh/id_ed25519_deploy.pub`) ke `/home/<user>/.ssh/authorized_keys` di VPS. Isi private key ke secret `VPS_SSH_KEY`.

3. **Systemd unit** `/etc/systemd/system/sistem-service-motor.service`:

```ini
[Unit]
Description=Sistem Service Motor API
After=network.target

[Service]
Type=simple
User=<user>
WorkingDirectory=/home/<user>/apps/sistem-service-motor
EnvironmentFile=/home/<user>/apps/sistem-service-motor/.env
ExecStart=/home/<user>/apps/sistem-service-motor/server
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```

4. **Aktifkan service**:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now sistem-service-motor
```

5. (Opsional) Pasang nginx/caddy sebagai reverse proxy + domain ke port `8080`.

### Deploy ke production

Setelah semua setup di atas, deploy & migrasi production berjalan otomatis saat **merge ke branch `main`**. Tidak perlu tag versi.

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
