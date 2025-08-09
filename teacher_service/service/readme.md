---

# Teacher Service — Redis-Focused Guide

A Go microservice (Gorilla Mux + PostgreSQL) that uses **Redis** as a read-through cache for teacher records.

## What Redis does here (in one minute)

* **Purpose:** Speed up reads and reduce PostgreSQL load.
* **Strategy:** *Lazy (read-through) caching.*

  * `GET /teacher/{id}` → check Redis first

    * **HIT:** return cached JSON
    * **MISS:** query Postgres → cache JSON in Redis → return
  * `POST /teacher` → **do not** cache (no one has requested it yet)
  * `PUT /teacher` → update Postgres **then** `DEL teacher:{id}` in Redis
  * `DELETE /teacher/{id}` → delete in Postgres **then** `DEL teacher:{id}` in Redis

**Cache key format:** `teacher:{id}`
**Value:** JSON (e.g., `{"id":1,"name":"Ayesha Khan","email":"ayesha.khan@example.com"}`)
**TTL:** none by default (you can add one if you want)

---

## Prerequisites

* Docker Desktop installed (you said you have this ✅)
* Go 1.21+
* `curl` for quick testing

---

## Environment variables

Create a `.env` in the service root:

```env
# App
PORT=8080

# Postgres (host port may be 5433 if you remapped)
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=teacherdb

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=yourStrongRedisPassword
```

> If you didn’t change Postgres mapping and 5432 is free, set `DB_PORT=5432`.
> If you set a Redis password (recommended), set `REDIS_PASSWORD` and use it in your Go client.

---

## Docker Compose (Redis + Postgres only)

> This binds services to **localhost** (not publicly exposed) and sets a **Redis password** for security.

```yaml
services:
  postgres:
    image: postgres:16
    container_name: teacher_postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: teacherdb
    ports:
      - "127.0.0.1:5433:5432" # host:container (change to 5432:5432 if free)
    restart: unless-stopped

  redis:
    image: redis:7
    container_name: teacher_redis
    command: ["redis-server", "--requirepass", "yourStrongRedisPassword"]
    ports:
      - "127.0.0.1:6379:6379" # local-only bind for security
    restart: unless-stopped
```

Bring them up:

```bash
docker-compose up -d
```

---

## Database setup (once)

Connect and create the table:

```bash
docker exec -it teacher_postgres bash
psql -U postgres -d teacherdb
```

```sql
CREATE TABLE IF NOT EXISTS teachers (
  id SERIAL PRIMARY KEY,
  name  TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE
);
```

`\\q` to exit psql, then `exit` to leave the container.

---

## Run the service

```bash
go run main.go
# Server started at :8080
```

---

## Test the Redis flow (HIT/MISS)

### 1) Create (no cache write here)

```bash
curl -X POST http://localhost:8080/teacher \
  -H "Content-Type: application/json" \
  -d '{"name":"Ayesha Khan","email":"ayesha.khan1@example.com"}'
```

### 2) First GET (MISS → fetch DB → cache)

```bash
curl http://localhost:8080/teacher/1
```

### 3) Check Redis directly (optional)

```bash
docker exec -it teacher_redis redis-cli -a yourStrongRedisPassword
GET teacher:1
# should return the JSON
```

### 4) Second GET (HIT → served from Redis)

```bash
curl http://localhost:8080/teacher/1
```

### 5) Update (invalidate cache)

```bash
curl -X PUT http://localhost:8080/teacher \
  -H "Content-Type: application/json" \
  -d '{"id":1,"name":"Ayesha K.","email":"ayesha.khan1@example.com"}'
```

* Service runs `DEL teacher:1` → next GET will MISS and refresh the cache.

### 6) Delete (invalidate cache)

```bash
curl -X DELETE http://localhost:8080/teacher/1
```

* Service runs `DEL teacher:1` → key removed.

---

## Security (Redis-focused)

* **Don’t expose Redis to the internet.** In compose, bind to `127.0.0.1`.
* **Require a password:** we used `--requirepass` in `command`.
* **Use app-level auth (JWT) on your APIs.**
* **Keep secrets in env vars** (never hardcode).
* **Firewall / Security Groups** if deploying to cloud.
* **Optional:** disable dangerous Redis commands, run on a private network only.

Minimal Go client config with password:

```go
RedisClient = redis.NewClient(&redis.Options{
  Addr:     os.Getenv("REDIS_ADDR"),       // e.g., "localhost:6379"
  Password: os.Getenv("REDIS_PASSWORD"),   // set in .env
})
```

---

## Optional improvements

* **TTL** on cache entries (e.g., 5 minutes):

  ```go
  RedisClient.Set(ctx, key, jsonBytes, 5*time.Minute)
  ```
* **Batch endpoints** (`POST /teachers`) → cache is still per-id.
* **Rate limiting** using Redis INCR with expiry.
* **Cache warming** on startup if you have a hot set.
* **ON CONFLICT DO NOTHING** for idempotent inserts.

---

## Troubleshooting

* **`pq: relation "teachers" does not exist`**
  You didn’t run the `CREATE TABLE`. Create it in `teacherdb`.

* **Port 5432 already in use**
  Map Postgres to `5433:5432` and set `DB_PORT=5433`.

* **`NOAUTH Authentication required` when using redis-cli**
  Use `redis-cli -a yourStrongRedisPassword` or set the password in your Go client.

* **Getting stale data**
  Ensure `PUT`/`DELETE` calls **always** `DEL teacher:{id}`.

---

## Quick API recap

* `POST /teacher` → create (DB only; no cache)
* `GET /teacher/{id}` → try Redis → fallback DB → cache
* `PUT /teacher` → update DB → `DEL teacher:{id}`
* `DELETE /teacher/{id}` → delete DB → `DEL teacher:{id}`

---

That’s it. If you want, I can add **TTL**, **rate limiting middleware**, or a **Dockerfile** for the Go service next.
