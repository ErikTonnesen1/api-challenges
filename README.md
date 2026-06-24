# Go API Challenges

A progressive series of Go API challenges for two developers learning to build servers with Go and Gin. Each challenge is designed to be completable in a few hours.

---

## Challenge 1 — Hello, Server

**Difficulty:** Beginner
**Estimated Time:** 1–2 hours

### The Problem

Build your first HTTP server in Go using the standard `net/http` package (no frameworks). This challenge is about getting comfortable with the basics: routing, handlers, and JSON responses.

### What to Build

Create a simple REST API for a to-do list stored **in memory** (no database). The server must support the following endpoints:

| Method | Path         | Description              |
|--------|--------------|--------------------------|
| GET    | `/todos`     | Return all to-do items   |
| POST   | `/todos`     | Create a new to-do item  |
| GET    | `/todos/:id` | Return a single to-do    |

A to-do item has this shape:

```json
{
  "id": 1,
  "title": "Buy groceries",
  "done": false
}
```

### Requirements

- Use only the Go standard library (`net/http`, `encoding/json`)
- IDs should be auto-incremented integers
- All responses must be `Content-Type: application/json`
- Return appropriate HTTP status codes (`200`, `201`, `404`)
- Data does not need to persist between server restarts

### Expected Behavior

```bash
# Create a to-do
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries"}'
# → 201 Created
# → {"id": 1, "title": "Buy groceries", "done": false}

# Get all todos
curl http://localhost:8080/todos
# → 200 OK
# → [{"id": 1, "title": "Buy groceries", "done": false}]

# Get one
curl http://localhost:8080/todos/1
# → 200 OK
# → {"id": 1, "title": "Buy groceries", "done": false}

# Not found
curl http://localhost:8080/todos/99
# → 404 Not Found
# → {"error": "not found"}
```

### Stretch Goals

- Add a `DELETE /todos/:id` endpoint
- Add a `PATCH /todos/:id` endpoint to toggle `done`

---

## Challenge 2 — Leveling Up with Gin

**Difficulty:** Beginner–Intermediate
**Estimated Time:** 2–3 hours

### The Problem

Rebuild the to-do API from Challenge 1 using the [Gin](https://github.com/gin-gonic/gin) framework. Then extend it with proper input validation and error handling.

### What to Build

Recreate all endpoints from Challenge 1, then add:

| Method | Path         | Description                    |
|--------|--------------|--------------------------------|
| PUT    | `/todos/:id` | Replace a to-do item entirely  |
| DELETE | `/todos/:id` | Remove a to-do item            |

### Requirements

- Use Gin for routing and handler context
- Validate incoming request bodies — a `POST` with no `title` should return `400 Bad Request`
- Return consistent error response envelopes:

```json
{
  "error": "title is required"
}
```

- Use Gin's `ShouldBindJSON` for request parsing
- Organize your code into at least two files: `main.go` and `handlers.go`

### Expected Behavior

```bash
# Missing title → 400
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{}'
# → 400 Bad Request
# → {"error": "title is required"}

# Delete a todo
curl -X DELETE http://localhost:8080/todos/1
# → 204 No Content

# Delete again
curl -X DELETE http://localhost:8080/todos/1
# → 404 Not Found
# → {"error": "not found"}
```

### Stretch Goals

- Add query param filtering: `GET /todos?done=true`
- Add basic request logging middleware that prints the method, path, and duration

---

## Challenge 3 — Persisting Data with SQLite

**Difficulty:** Intermediate
**Estimated Time:** 2–4 hours

### The Problem

Replace the in-memory store with a real database. You'll connect your Gin API to SQLite using the `database/sql` package with the `mattn/go-sqlite3` driver and learn how to manage schema and perform CRUD operations against an actual database.

### What to Build

Extend the to-do API so that all data persists to a SQLite file (`todos.db`). The API surface stays the same, but the storage layer moves to SQL.

### Requirements

- Use `database/sql` with `github.com/mattn/go-sqlite3`
- Create the `todos` table on server startup if it doesn't exist
- All CRUD endpoints must read from and write to the database
- Wrap your DB access in a simple repository struct (e.g., `TodoRepository`) to keep handlers clean
- Handle SQL errors gracefully — don't let a DB error panic the server

### Schema

```sql
CREATE TABLE IF NOT EXISTS todos (
  id    INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  done  INTEGER NOT NULL DEFAULT 0
);
```

### Expected Behavior

All behavior from Challenge 2 should be preserved, but now restarting the server should retain previously created todos.

```bash
# Create and restart server — data survives
curl -X POST http://localhost:8080/todos -d '{"title": "Persisted!"}' -H "Content-Type: application/json"
# → {"id": 1, "title": "Persisted!", "done": false}

# Restart the server, then:
curl http://localhost:8080/todos
# → [{"id": 1, "title": "Persisted!", "done": false}]
```

### Stretch Goals

- Add pagination: `GET /todos?page=1&limit=10`
- Add a `created_at` timestamp column and return it in responses

---

## Challenge 4 — Authentication with JWT

**Difficulty:** Intermediate–Advanced
**Estimated Time:** 3–5 hours

### The Problem

Secure your API. You'll build a user registration and login system, issue JSON Web Tokens on successful login, and protect your to-do endpoints so that each user only sees their own data.

### What to Build

Add an `auth` layer to the existing API:

| Method | Path             | Description                          |
|--------|------------------|--------------------------------------|
| POST   | `/auth/register` | Create a new user account            |
| POST   | `/auth/login`    | Authenticate and receive a JWT token |

Then protect all `/todos` routes so they require a valid `Authorization: Bearer <token>` header. Users should only be able to read and modify their own todos.

### Requirements

- Use `golang-jwt/jwt` for token creation and verification
- Store users in SQLite with hashed passwords (use `golang.org/x/crypto/bcrypt`)
- Associate todos with a `user_id` foreign key
- Write a Gin middleware (`AuthMiddleware`) that validates the token and sets the user on the request context
- Tokens should expire after 24 hours

### Schema Additions

```sql
CREATE TABLE IF NOT EXISTS users (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  email         TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL
);

-- Add user_id to todos
ALTER TABLE todos ADD COLUMN user_id INTEGER REFERENCES users(id);
```

### Expected Behavior

```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -d '{"email": "dev@example.com", "password": "secret"}' \
  -H "Content-Type: application/json"
# → 201 Created
# → {"message": "registered successfully"}

# Login
curl -X POST http://localhost:8080/auth/login \
  -d '{"email": "dev@example.com", "password": "secret"}' \
  -H "Content-Type: application/json"
# → 200 OK
# → {"token": "<jwt>"}

# Access protected route
curl http://localhost:8080/todos \
  -H "Authorization: Bearer <jwt>"
# → 200 OK — only this user's todos

# No token
curl http://localhost:8080/todos
# → 401 Unauthorized
# → {"error": "authorization required"}
```

### Stretch Goals

- Add a `POST /auth/refresh` endpoint to issue a new token
- Return `403 Forbidden` (not `404`) when a user tries to access another user's todo by ID

---

## Challenge 5 — Background Jobs & Rate Limiting

**Difficulty:** Advanced
**Estimated Time:** 4–6 hours

### The Problem

Production APIs need more than just CRUD. In this challenge you'll add two real-world concerns: a background worker that processes work asynchronously using Go channels and goroutines, and a rate limiter that protects your endpoints from abuse.

### What to Build

**Part A — Background Email Notifications**

When a user marks a to-do as `done`, queue a "notification" job. A background worker goroutine should pick up the job and log a message simulating an email send:

```
[notifier] Sending email to dev@example.com: "Buy groceries" is complete!
```

Use a Go channel as the job queue. The worker should run in a goroutine started at server boot.

**Part B — Rate Limiting**

Add a middleware that limits each IP address to **60 requests per minute**. Requests that exceed the limit should receive a `429 Too Many Requests` response.

```json
{
  "error": "rate limit exceeded, try again later"
}
```

### Requirements

- The notification channel must be buffered (capacity: 100)
- The worker goroutine must handle a server shutdown signal gracefully (use `context.Context` or `os.Signal`)
- Rate limiting must be per-IP
- Use a sliding window or token bucket approach (you may use `golang.org/x/time/rate`)
- Rate limit state is in-memory (no Redis required)
- Both features must work correctly alongside the JWT auth from Challenge 4

### Expected Behavior

```bash
# Mark a todo done
curl -X PATCH http://localhost:8080/todos/1 \
  -H "Authorization: Bearer <jwt>" \
  -d '{"done": true}' \
  -H "Content-Type: application/json"
# → 200 OK

# Server logs:
# [notifier] Sending email to dev@example.com: "Buy groceries" is complete!

# Hammer the API
for i in $(seq 1 65); do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/todos; done
# → First 60: 200 (or 401 without token)
# → Remaining: 429
```

### Stretch Goals

- Make the rate limit configurable via environment variable (`RATE_LIMIT_RPM`)
- Add a `/health` endpoint that is exempt from rate limiting and returns server uptime
- Write a test that spins up the server and verifies the rate limiter behavior

---

## Tips & Resources

- **Project layout:** Keep it simple. `main.go`, `handlers.go`, `repository.go`, `middleware.go` is plenty for these challenges.
- **Testing your API:** [httpie](https://httpie.io/) (`http POST :8080/todos title="test"`) is friendlier than curl for quick iteration.
- **Go module setup:** `go mod init github.com/yourname/go-challenges && go mod tidy`
- **Recommended packages:**
  - Gin: `github.com/gin-gonic/gin`
  - SQLite driver: `github.com/mattn/go-sqlite3`
  - JWT: `github.com/golang-jwt/jwt/v5`
  - Bcrypt: `golang.org/x/crypto/bcrypt`
  - Rate limiter: `golang.org/x/time/rate`
