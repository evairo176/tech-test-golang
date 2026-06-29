# Tech Test - Golang (3 Jam)

Technical test dengan 2 task Web API menggunakan Golang + SQLite.

## Task 1: Web API (SQL)

CRUD Person table + endpoint GetCountry by name.

| Name | Country |
|------|---------|
| Adam | Kuala Lumpur |
| John | Singapore |
| Henry | Singapore |
| Dominic | Thailand |

**Endpoint:**
- `GET /GetCountry/{name}` → Returns country string

## Task 2: Web API (Integration)

Consume timeapi.io untuk mendapatkan current time by timezone.

**Endpoint:**
- `GET /GetCurrentTime/{timezone}` → Returns JSON time response

## Cara Run

```bash
go run main.go
# atau
go build -o tech-test . && ./tech-test
```

Server berjalan di `http://localhost:8080`

## Contoh Request

```bash
# Task 1
curl http://localhost:8080/GetCountry/Adam
# Output: Kuala Lumpur

# Task 2
curl http://localhost:8090/GetCurrentTime/Europe%2FAmsterdam
# Output: JSON { year, month, day, hour, minute, seconds, ... }
```
