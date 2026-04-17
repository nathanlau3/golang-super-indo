# Super Indo Product API

REST API sederhana buat manage data produk Super Indo. Tech stack utama Go + Gin, PostgreSQL, dan Redis buat caching.

## Tech Stack

- Go 1.21+ (Gin)
- PostgreSQL 15+ (raw query, tanpa ORM)
- Redis 7+
- Viper (config)
- Docker

## Project Structure

Struktur project pakai konsep Hexagonal Architecture / Ports & Adapters, di-group per module.

```
.
├── api/                                        # Presentation layer
│   └── product/
│       ├── handler.go                          #   HTTP handler (Gin)
│       ├── handler_test.go
│       └── dto.go                              #   Request & response struct
│
├── cmd/api/
│   └── main.go                                 # Entry point, wiring dependency
│
├── internal/product/                           # Core business logic
│   ├── domain/
│   │   ├── product.go                          #   Entity, value object, error
│   │   └── product_test.go
│   ├── port/
│   │   ├── inbound.go                          #   Use case interface (dipanggil handler)
│   │   └── outbound.go                         #   Repository interface (diimplementasi adapter)
│   └── usecase/
│       ├── create_product.go                   #   Tiap use case 1 file, 1 method Execute()
│       ├── get_products.go
│       └── get_product_by_id.go
│
├── pkg/
│   ├── config/
│   │   └── config.go                           # Viper config loader
│   └── infrastructure/
│       ├── adapter/
│       │   └── product_repository.go           # Implementasi repository (Postgres + Redis)
│       ├── postgres/
│       │   ├── postgres.go                     # DB connection
│       │   ├── migration.go                    # Embedded SQL migration
│       │   ├── seeder.go                       # Seed data produk
│       │   └── migrations/                     # File .sql
│       └── redis/
│           └── redis.go                        # Redis connection
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

### Alur Dependency

```
api/product (handler)
    │
    ▼
internal/product/port (inbound)  ←── usecase (implementasi)
                                          │
                                          ▼
internal/product/port (outbound) ←── pkg/infrastructure/adapter (implementasi)
```

Use case cuma bergantung ke port (interface), gak tau sama sekali soal Postgres atau Redis. Semua infrastructure detail ada di adapter.

## Setup & Run

### Pakai Docker (paling gampang)

```bash
# start postgres + redis
make infra-up

# jalankan app
make run

# isi data awal
make seed
```

Atau kalau mau full dockerized (app + infra sekaligus):
```bash
make docker-up
```

### Manual

```bash
cp .env.example .env
# edit .env sesuai kebutuhan

go mod tidy
go run cmd/api/main.go

# seed data
go run cmd/api/main.go seed
```

### Makefile

```bash
make run            # go run
make seed           # jalankan seeder
make test           # unit test + coverage
make build          # compile binary
make infra-up       # docker postgres + redis aja
make infra-down     # stop docker
make docker-up      # full stack (app + infra)
make docker-down    # stop semua
```

## Migration

SQL migration ada di `pkg/infrastructure/postgres/migrations/`. File ini di-embed pakai `go:embed` dan otomatis jalan waktu app start, jadi gak perlu tool migration tambahan.

Kalau mau manual:
```bash
psql -U postgres -d superindo -f pkg/infrastructure/postgres/migrations/20240115100000_create_products_table.up.sql
```

## API

Base URL: `http://localhost:8080`

### POST /product

Tambah produk baru. Tipe yang valid: `Sayuran`, `Protein`, `Buah`, `Snack`.

```bash
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bayam Segar",
    "type": "Sayuran",
    "price": 5500,
    "description": "Bayam hijau segar 250g",
    "stock": 100
  }'
```

Response `201`:
```json
{
  "status": 201,
  "message": "produk berhasil ditambahkan",
  "data": {
    "id": 1,
    "name": "Bayam Segar",
    "type": "Sayuran",
    "price": 5500,
    "description": "Bayam hijau segar 250g",
    "stock": 100,
    "created_at": "2024-01-15T10:00:00+07:00",
    "updated_at": "2024-01-15T10:00:00+07:00"
  }
}
```

### GET /product

Ambil list produk. Support search, filter, sort, dan pagination.

| Param   | Type   | Default | Keterangan |
|---------|--------|---------|------------|
| search  | string | -       | Cari nama produk |
| type    | string | -       | Filter: Sayuran / Protein / Buah / Snack |
| sort_by | string | date    | Sort: name / price / date |
| order   | string | desc    | asc / desc |
| page    | int    | 1       | Halaman |
| limit   | int    | 10      | Per halaman (max 50) |

```bash
# semua produk
curl http://localhost:8080/product

# cari "bayam"
curl "http://localhost:8080/product?search=bayam"

# filter sayuran, sort harga murah dulu
curl "http://localhost:8080/product?type=Sayuran&sort_by=price&order=asc"
```

Response `200`:
```json
{
  "status": 200,
  "message": "success",
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 20,
    "total_page": 2
  }
}
```

### GET /product/:id

Ambil detail produk berdasarkan ID.

```bash
curl http://localhost:8080/product/1
```

Response `200`:
```json
{
  "status": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Bayam Segar",
    "type": "Sayuran",
    "price": 5500,
    "description": "Bayam hijau segar 250g",
    "stock": 100,
    "created_at": "2024-01-15T10:00:00+07:00",
    "updated_at": "2024-01-15T10:00:00+07:00"
  }
}
```

## Caching

Pakai Redis dengan TTL 5 menit. List produk di-cache berdasarkan kombinasi query param, detail produk di-cache per ID. Cache list otomatis di-invalidasi waktu ada produk baru masuk.

App tetap jalan normal kalau Redis mati, cuma tanpa cache aja.

## Test

```bash
make test
```

Atau langsung:
```bash
go test ./... -v -cover
```
