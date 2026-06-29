# Tech Test - Golang (3 Jam)

Technical test dengan 2 task Web API menggunakan **Golang + Echo Framework + SQLite**.

## Tech Stack

- **Language:** Go 1.22+
- **Framework:** Echo v4
- **Database:** SQLite (go-sqlite3)
- **Migration:** golang-migrate/migrate

---

## Project Structure

```
tech-test-golang/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── config/
│   └── config.go                # Config, DB init, migrations & seeders
├── database/
│   └── migrations/
│       ├── 000001_create_person_table.up.sql
│       └── 000001_create_person_table.down.sql
├── internal/
│   ├── handler/
│   │   └── handler.go           # HTTP handlers (controller layer)
│   ├── model/
│   │   └── model.go             # Data models (Person, TimeResponse)
│   ├── repository/
│   │   └── person_repository.go # Database queries (repository layer)
│   └── response/
│       └── response.go          # Standard API response helper
├── go.mod
├── go.sum
└── README.md
```

---

## Task 1: Web API (SQL)

CRUD Person table + endpoint GetCountry by name.

### Person Table

| Name    | Country       |
|---------|---------------|
| Adam    | Kuala Lumpur  |
| John    | Singapore     |
| Henry   | Singapore     |
| Dominic | Thailand      |

### Migrations

Migration files berada di `database/migrations/`. Dijalankan otomatis saat server start.

- `000001_create_person_table.up.sql` — Create table
- `000001_create_person_table.down.sql` — Drop table (rollback)

### Seeders

Seeder dijalankan otomatis saat server start jika tabel Person kosong. Menginsert 4 data awal (Adam, John, Henry, Dominic).

### API Endpoints

#### 1. Get All Persons

Select semua data person dari database.

```
GET /api/persons
```

**Response:**
```json
{
    "status": "success",
    "code": 200,
    "message": "Found 4 persons",
    "data": [
        { "name": "Adam", "country": "Kuala Lumpur" },
        { "name": "Dominic", "country": "Thailand" },
        { "name": "Henry", "country": "Singapore" },
        { "name": "John", "country": "Singapore" }
    ]
}
```

#### 2. Get Country by Person Name

Cari country berdasarkan nama person. Case-insensitive.

```
GET /api/GetCountry/{name}
```

**Contoh:** `GET /api/GetCountry/Adam`

**Response:**
```json
{
    "status": "success",
    "code": 200,
    "message": "Country for 'Adam' found",
    "data": {
        "name": "Adam",
        "country": "Kuala Lumpur"
    }
}
```

**Response jika tidak ditemukan:**
```json
{
    "status": "error",
    "code": 404,
    "message": "Person 'Nobody' not found"
}
```

#### 3. Create Person

Insert data person baru ke database.

```
POST /api/person
Content-Type: application/json
```

**Body:**
```json
{
    "name": "Budi",
    "country": "Indonesia"
}
```

**Response:**
```json
{
    "status": "success",
    "code": 201,
    "message": "Person created successfully",
    "data": {
        "name": "Budi",
        "country": "Indonesia",
        "rowsAffected": 1
    }
}
```

**Response jika nama sudah ada:**
```json
{
    "status": "error",
    "code": 409,
    "message": "Person 'Adam' already exists"
}
```

#### 4. Delete Person

Hapus data person berdasarkan nama. Case-insensitive.

```
DELETE /api/person/{name}
```

**Contoh:** `DELETE /api/person/Budi`

**Response:**
```json
{
    "status": "success",
    "code": 200,
    "message": "Person 'Budi' deleted"
}
```

---

## Task 2: Web API (Integration)

Consume timeapi.io untuk mendapatkan current time by timezone.

```
GET /api/GetCurrentTime/{timezone}
```

**Contoh:** `GET /api/GetCurrentTime/Europe/Amsterdam`

> ⚠️ Timezone yang mengandung `/` perlu di-URL-encode menjadi `%2F`.
> Contoh: `Europe/Amsterdam` → `Europe%2FAmsterdam`

**Response:**
```json
{
    "status": "success",
    "code": 200,
    "message": "Current time for 'Europe/Amsterdam'",
    "data": {
        "year": 2026,
        "month": 6,
        "day": 29,
        "hour": 18,
        "minute": 38,
        "seconds": 55,
        "milliSeconds": 224,
        "dateTime": "2026-06-29T18:38:55.2244932",
        "date": "06/29/2026",
        "time": "18:38",
        "timeZone": "Europe/Amsterdam",
        "dayOfWeek": "Monday",
        "dstActive": true
    }
}
```

**Contoh timezone lain:**
- `Asia/Jakarta` → `Asia%2FJakarta`
- `America/New_York` → `America%2FNew_York`
- `UTC` → `UTC`

---

## Cara Run

```bash
# Langsung run
go run ./cmd/server

# Atau build dulu
go build -o tech-test ./cmd/server && ./tech-test
```

Server berjalan di `http://localhost:8080` (default), atau set via env:

```bash
PORT=3000 ./tech-test
```

---

## Contoh Test via cURL

```bash
# Select all persons
curl http://localhost:8080/api/persons

# Get country by name
curl http://localhost:8080/api/GetCountry/Adam

# Create new person
curl -X POST http://localhost:8080/api/person \
  -H "Content-Type: application/json" \
  -d '{"name":"Budi","country":"Indonesia"}'

# Delete person
curl -X DELETE http://localhost:8080/api/person/Budi

# Get current time
curl http://localhost:8080/api/GetCurrentTime/Asia%2FJakarta
```
