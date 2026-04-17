# Super Indo Product API

REST API untuk manage data produk Super Indo, dilengkapi JWT authentication. Tech stack Go + Gin, PostgreSQL, Redis, dan RSA256 JWT.

## Tech Stack

- Go 1.21+ (Gin)
- PostgreSQL 15+ (raw query, tanpa ORM)
- Redis 7+
- JWT RS256 (RSA key pair)
- Viper (config)
- Docker

## Project Structure

Struktur project pakai konsep Hexagonal Architecture / Ports & Adapters, di-group per module.

```
.
├── api/
│   ├── auth/
│   │   ├── handler.go
│   │   └── dto.go
│   └── product/
│       ├── handler.go
│       ├── handler_test.go
│       └── dto.go
│
├── cmd/api/
│   ├── main.go
│   └── bootstrap.go
│
├── internal/
│   ├── auth/
│   │   ├── domain/
│   │   │   └── user.go
│   │   ├── port/
│   │   │   ├── inbound.go
│   │   │   └── outbound.go
│   │   └── usecase/
│   │       ├── register.go
│   │       └── login.go
│   └── product/
│       ├── domain/
│       │   ├── product.go
│       │   └── product_test.go
│       ├── port/
│       │   ├── inbound.go
│       │   └── outbound.go
│       └── usecase/
│           ├── create_product.go
│           ├── get_products.go
│           └── get_product_by_id.go
│
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── credentials/
│   │   ├── jwtRS256.key
│   │   └── jwtRS256.key.pub
│   ├── middleware/
│   │   └── auth.go
│   └── infrastructure/
│       ├── adapter/
│       │   ├── product_repository.go
│       │   └── user_repository.go
│       ├── jwt/
│       │   └── jwt.go
│       ├── postgres/
│       │   ├── postgres.go
│       │   ├── migration.go
│       │   ├── seeder.go
│       │   └── migrations/
│       └── redis/
│           └── redis.go
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

### Alur Dependency

```
api (handler)
    │
    ▼
internal/*/port (inbound)  ←── usecase (implementasi)
                                    │
                                    ▼
internal/*/port (outbound) ←── pkg/infrastructure/adapter (implementasi)
```

Use case cuma bergantung ke port (interface), gak tau sama sekali soal Postgres, Redis, atau JWT. Semua infrastructure detail ada di adapter.

## Setup & Run

### Pakai Docker (paling gampang)

```bash
make infra-up
make run
make seed
```

Atau full dockerized (app + infra sekaligus):
```bash
make docker-up
```

### Manual

```bash
cp .env.example .env
go mod tidy
make run
make seed
```

### Generate RSA Key Pair

Key pair untuk JWT signing harus ada di `pkg/credentials/` sebelum app bisa jalan.

```bash
mkdir -p pkg/credentials
openssl genrsa -out pkg/credentials/jwtRS256.key 4096
openssl rsa -in pkg/credentials/jwtRS256.key -pubout -out pkg/credentials/jwtRS256.key.pub
```

### Makefile

```bash
make run            # jalankan app
make seed           # seed produk + user
make test           # unit test + coverage
make build          # compile binary
make infra-up       # docker postgres + redis
make infra-down     # stop docker
make docker-up      # full stack (app + infra)
make docker-down    # stop semua
```

## Environment Variables

| Variable | Default | Keterangan |
|----------|---------|------------|
| APP_PORT | 8080 | Port server |
| DB_HOST | localhost | Host PostgreSQL |
| DB_PORT | 5432 | Port PostgreSQL |
| DB_USER | postgres | User database |
| DB_PASSWORD | postgres | Password database |
| DB_NAME | superindo | Nama database |
| REDIS_ADDR | localhost:6379 | Alamat Redis |
| REDIS_PASSWORD | | Password Redis |
| JWT_PRIVATE_KEY_PATH | pkg/credentials/jwtRS256.key | Path RSA private key |
| JWT_PUBLIC_KEY_PATH | pkg/credentials/jwtRS256.key.pub | Path RSA public key |

## Migration

SQL migration ada di `pkg/infrastructure/postgres/migrations/`. File ini di-embed pakai `go:embed` dan otomatis jalan waktu app start.

Table yang dibuat:
- `products` — data produk
- `users` — data user untuk authentication

## API

Base URL: `http://localhost:8080`

### Auth

#### POST /auth/register

Registrasi user baru. Password minimal 6 karakter.

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

Response `201`:
```json
{
  "status": 201,
  "message": "registrasi berhasil",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2024-01-15T10:00:00+07:00",
    "updated_at": "2024-01-15T10:00:00+07:00"
  }
}
```

#### POST /auth/login

Login dan dapatkan JWT token. Token berlaku 24 jam.

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@superindo.co.id",
    "password": "admin123"
  }'
```

Response `200`:
```json
{
  "status": 200,
  "message": "login berhasil",
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "admin@superindo.co.id",
      "name": "Admin Super Indo",
      "created_at": "2024-01-15T10:00:00+07:00",
      "updated_at": "2024-01-15T10:00:00+07:00"
    }
  }
}
```

### Product

#### POST /product (Protected)

Tambah produk baru. Memerlukan JWT token di header Authorization.

Tipe yang valid: `Sayuran`, `Protein`, `Buah`, `Snack`.

```bash
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
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

#### GET /product (Public)

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
curl http://localhost:8080/product
curl "http://localhost:8080/product?search=bayam"
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

#### GET /product/:id (Public)

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

## Seed Data

`make seed` akan mengisi data awal:

**20 produk** (Sayuran, Protein, Buah, Snack)

**3 user:**

| Email | Password |
|-------|----------|
| admin@superindo.co.id | admin123 |
| kasir@superindo.co.id | kasir123 |
| manager@superindo.co.id | manager123 |

## Caching

Pakai Redis dengan TTL 5 menit. List produk di-cache berdasarkan kombinasi query param, detail produk di-cache per ID. Cache list otomatis di-invalidasi waktu ada produk baru masuk.

App tetap jalan normal kalau Redis mati, cuma tanpa cache aja.

## Test

```bash
make test
```
