# Tech Test - Golang (3 Jam)

Technical test dengan 2 task Web API menggunakan **Golang + Echo Framework + SQLite**.

## Tech Stack

- **Language:** Go 1.22+
- **Framework:** Echo v4
- **Database:** SQLite (go-sqlite3)

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

### SQL Scripts

File SQL tersedia di folder `sql/` untuk dijalankan manual:

```bash
sqlite3 persons.db < sql/01_create_table.sql
sqlite3 persons.db < sql/02_insert_data.sql
sqlite3 persons.db < sql/03_select_and_stored_procedure.sql
```

| File | Deskripsi |
|------|-----------|
| `sql/01_create_table.sql` | Create table script |
| `sql/02_insert_data.sql` | Insert data script (4 data awal) |
| `sql/03_select_and_stored_procedure.sql` | Select by name + stored procedure equivalent (PostgreSQL, MySQL, SQL Server) |

### API Endpoints

#### 1. Get All Persons

Select semua data person dari database.

```
GET /persons
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
GET /GetCountry/{name}
```

**Contoh:** `GET /GetCountry/Adam`

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
POST /person
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
DELETE /person/{name}
```

**Contoh:** `DELETE /person/Budi`

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
GET /GetCurrentTime/{timezone}
```

**Contoh:** `GET /GetCurrentTime/Europe/Amsterdam`

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
go run main.go

# Atau build dulu
go build -o tech-test . && ./tech-test
```

Server berjalan di `http://localhost:8080` (default), atau set via env:

```bash
PORT=3000 ./tech-test
```

---

## Contoh Test via cURL

```bash
# Select all persons
curl http://localhost:8080/persons

# Get country by name
curl http://localhost:8080/GetCountry/Adam

# Create new person
curl -X POST http://localhost:8080/person \
  -H "Content-Type: application/json" \
  -d '{"name":"Budi","country":"Indonesia"}'

# Delete person
curl -X DELETE http://localhost:8080/person/Budi

# Get current time
curl http://localhost:8080/GetCurrentTime/Asia%2FJakarta
```
