# Go API + React Challenges

A progressive series of challenges for two developers learning to build and connect a Go API to a React frontend. Each challenge is designed to be completable in a few hours.

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

Rebuild the to-do API from Challenge 1 using the [Gin](https://github.com/gin-gonic/gin) framework. Then extend it with proper input validation, error handling, and integration tests.

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
- Write integration tests using `httptest` that verify each route is wired correctly

### Testing

Use Go's `httptest` package to test your routes against a real Gin router without binding to a port. Your integration tests should verify status codes and response shapes — not business logic (that belongs in unit tests).

```go
func TestGetTodos(t *testing.T) {
    router := setupRouter()

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/todos", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response []map[string]any
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.NotNil(t, response)
}
```

Use `TestMain` to set Gin to test mode once for the whole package:

```go
func TestMain(m *testing.M) {
    gin.SetMode(gin.TestMode)
    os.Exit(m.Run())
}
```

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
curl -X POST http://localhost:8080/todos \
  -d '{"title": "Persisted!"}' \
  -H "Content-Type: application/json"
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

## Challenge 6 — Your First React UI

**Difficulty:** Beginner (Frontend)
**Estimated Time:** 2–4 hours

### The Problem

You have a working Go API. Now build a UI for it. This challenge is about learning to think in components and understanding that **UI is a function of state** — when state changes, the UI re-renders automatically.

This will feel different from backend code. You're not writing procedures that run top to bottom; you're describing what the UI should look like *given* the current data, and React figures out how to update the DOM.

### What to Build

A single-page todo app using [Vite](https://vitejs.dev/) + React that connects to your Go API.

- Display all todos in a list
- Add a new todo via a form
- Toggle a todo's `done` state
- Delete a todo

### Setup

```bash
npm create vite@latest todo-ui -- --template react
cd todo-ui
npm install axios
npm run dev
```

Make sure your Go API has CORS enabled. Add this middleware before your routes:

```go
// go get github.com/gin-contrib/cors
import "github.com/gin-contrib/cors"

r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:5173"},
    AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
    AllowHeaders: []string{"Content-Type"},
}))
```

### Requirements

- Use `axios` for all API calls
- Use `useState` to store your todos list locally after fetching
- Use `useEffect` to fetch todos when the component first mounts
- Break the UI into at least three components: `TodoList`, `TodoItem`, `AddTodoForm`
- No CSS frameworks required — basic inline styles or a plain CSS file is fine

### Key Concepts to Understand

**State drives the UI.** When you call `setTodos(...)`, React re-renders automatically. You never touch the DOM directly.

```jsx
const [todos, setTodos] = useState([])

useEffect(() => {
    axios.get('http://localhost:8080/todos')
        .then(res => setTodos(res.data))
}, []) // empty array = run once on mount
```

**Optimistic vs. server-driven updates.** After creating a todo, you have two choices: re-fetch the whole list from the API, or append the returned item to local state. Try both and notice the tradeoff.

### Expected Behavior

- Page loads and displays all existing todos from the API
- Submitting the form creates a new todo and it appears in the list
- Clicking a toggle button updates `done` state via the API and reflects in the UI
- Clicking delete removes the todo from the list

### Stretch Goals

- Add a loading spinner while the initial fetch is in progress
- Add an error message if the API is unreachable
- Filter todos by `done` state with tab buttons (All / Active / Completed)

---

## Challenge 7 — Server State with React Query

**Difficulty:** Intermediate (Frontend)
**Estimated Time:** 2–3 hours

### The Problem

In Challenge 6 you managed server data manually — fetching it, storing it in `useState`, and keeping it in sync with the API by hand. This gets messy fast. React Query treats server data as a first-class concern with caching, background refetching, and loading/error states built in.

As a backend developer you'll recognize the concepts: stale-while-revalidate, cache invalidation, retry logic. React Query is essentially a state machine for server data living in the browser.

### What to Build

Refactor the Challenge 6 UI to use [TanStack Query](https://tanstack.com/query/latest) (React Query). The UI behavior stays identical — the internals change completely.

### Setup

```bash
npm install @tanstack/react-query
```

Wrap your app in the provider:

```jsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

const queryClient = new QueryClient()

function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <TodoApp />
        </QueryClientProvider>
    )
}
```

### Requirements

- Replace all `useState` + `useEffect` + `axios` fetch logic with `useQuery` for reads
- Replace all manual POST/PATCH/DELETE calls with `useMutation`
- After a successful mutation (create, update, delete), invalidate the todos query so the list refreshes automatically
- Handle loading and error states using the values returned by `useQuery`

### Key Concepts to Understand

**Queries fetch data:**

```jsx
const { data: todos, isLoading, isError } = useQuery({
    queryKey: ['todos'],
    queryFn: () => axios.get('/todos').then(res => res.data)
})
```

**Mutations change data, then invalidate:**

```jsx
const queryClient = useQueryClient()

const deleteMutation = useMutation({
    mutationFn: (id) => axios.delete(`/todos/${id}`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['todos'] })
})
```

**Notice what disappears:** no `useState` for todos, no manual re-fetch after mutations, no loading booleans to manage by hand.

### Expected Behavior

All behavior from Challenge 6 is preserved. Additionally:

- If you open the app in two tabs and create a todo in one, the other tab should eventually reflect it (React Query refetches on window focus by default)
- Loading and error states are handled gracefully

### Stretch Goals

- Add optimistic updates to the toggle mutation so the UI updates instantly before the API responds
- Install the [React Query Devtools](https://tanstack.com/query/latest/docs/framework/react/devtools) and explore the cache

---

## Challenge 8 — Next.js & the BFF Pattern

**Difficulty:** Intermediate–Advanced (Frontend)
**Estimated Time:** 4–6 hours

### The Problem

Migrate your Vite React app into [Next.js](https://nextjs.org/) using the App Router. Then implement a **Backend For Frontend (BFF)** layer using Next.js Route Handlers that sit between your React UI and the Go API.

### What is a BFF?

Your Go API is a domain API — it's designed to serve any consumer (web, mobile, third parties). It returns general-purpose data shapes and doesn't know or care about your UI's specific needs.

A BFF is a thin server layer owned by the frontend team. It:

- Aggregates multiple API calls into one response shaped for the UI
- Handles auth tokens server-side so they never touch the browser
- Translates or filters data so components get exactly what they need

```
Browser → Next.js Route Handlers (BFF) → Go API
```

The browser never talks directly to Go. It talks to Next, and Next talks to Go.

### What to Build

Migrate the todo UI into Next.js, then add Route Handlers as a BFF layer:

| Next.js Route           | Proxies to Go API       |
|-------------------------|-------------------------|
| GET `/api/todos`        | GET `/todos`            |
| POST `/api/todos`       | POST `/todos`           |
| PATCH `/api/todos/:id`  | PATCH `/todos/:id`      |
| DELETE `/api/todos/:id` | DELETE `/todos/:id`     |

Your React components should call `/api/todos` (the BFF), never `localhost:8080` directly.

### Setup

```bash
npx create-next-app@latest todo-next --app --no-src-dir
cd todo-next
npm install @tanstack/react-query axios
```

### Requirements

- Use the App Router (`app/` directory)
- The todo list page should be a **Server Component** that fetches todos at render time — no `useEffect`, no client-side fetch
- Toggle and delete interactions require client state, so extract those into a `'use client'` component
- Write Route Handlers in `app/api/todos/route.ts` and `app/api/todos/[id]/route.ts` that proxy to your Go API
- Keep the Go API URL in an environment variable (`GO_API_URL=http://localhost:8080`)

### Key Concepts to Understand

**Server Components fetch on the server:**

```tsx
// app/todos/page.tsx — no 'use client', runs on the server
export default async function TodosPage() {
    const res = await fetch(`${process.env.GO_API_URL}/todos`)
    const todos = await res.json()

    return <TodoList todos={todos} />
}
```

**Route Handlers are your BFF endpoints:**

```ts
// app/api/todos/route.ts
export async function GET() {
    const res = await fetch(`${process.env.GO_API_URL}/todos`)
    const data = await res.json()
    return Response.json(data)
}

export async function POST(request: Request) {
    const body = await request.json()
    const res = await fetch(`${process.env.GO_API_URL}/todos`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    })
    const data = await res.json()
    return Response.json(data, { status: res.status })
}
```

**Client Components handle interactivity:**

```tsx
'use client'

// Mutations, click handlers, and React Query live here
```

### Expected Behavior

- The todo list renders on the server — view source in the browser and you should see the todo HTML, not a loading spinner
- Creating, toggling, and deleting todos works via the BFF route handlers
- The Go API URL never appears in browser network requests — all calls go to `/api/todos`

### Stretch Goals

- Add a shape transformation in the BFF — return `{ id, title, completed }` instead of `{ id, title, done }` and update the UI to match. Notice how the BFF decouples the UI from the API contract
- Add error handling in the Route Handlers that returns consistent error envelopes regardless of what the Go API returns

---

## Challenge 9 — Auth in the UI

**Difficulty:** Advanced
**Estimated Time:** 4–6 hours

### The Problem

Wire up the JWT authentication from Challenge 4 to your Next.js frontend. This is where the BFF pattern really pays off — instead of storing JWTs in the browser (a security risk), the token lives server-side in an HttpOnly cookie that JavaScript can never read.

### The Auth Flow

```
1. User submits login form
2. Next.js Route Handler receives credentials
3. Route Handler calls Go API → POST /auth/login
4. Go API returns JWT
5. Route Handler stores JWT in an HttpOnly cookie (never exposed to browser JS)
6. Subsequent BFF requests read the cookie and forward the token to the Go API
7. On logout, the cookie is cleared
```

### What to Build

| Page / Route              | Description                                         |
|---------------------------|-----------------------------------------------------|
| `/login`                  | Login form, calls BFF auth route                    |
| `/api/auth/login`         | BFF route — calls Go API, sets HttpOnly cookie      |
| `/api/auth/logout`        | BFF route — clears the cookie                       |
| `/todos`                  | Protected page — redirects to `/login` if no cookie |

### Requirements

- Store the JWT in an HttpOnly cookie via the BFF — never in `localStorage` or accessible JS
- All BFF routes that proxy to protected Go endpoints must read the cookie and forward it as `Authorization: Bearer <token>`
- The `/todos` page must redirect to `/login` if the auth cookie is missing
- Use Next.js middleware (`middleware.ts`) to protect routes at the edge

### Key Concepts to Understand

**Setting an HttpOnly cookie in a Route Handler:**

```ts
// app/api/auth/login/route.ts
import { cookies } from 'next/headers'

export async function POST(request: Request) {
    const body = await request.json()
    const res = await fetch(`${process.env.GO_API_URL}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    })

    const { token } = await res.json()

    cookies().set('auth_token', token, {
        httpOnly: true,    // JS cannot read this
        secure: true,      // HTTPS only in production
        sameSite: 'lax',
        maxAge: 60 * 60 * 24 // 24 hours
    })

    return Response.json({ success: true })
}
```

**Forwarding the token from the BFF to Go:**

```ts
// app/api/todos/route.ts
import { cookies } from 'next/headers'

export async function GET() {
    const token = cookies().get('auth_token')?.value

    const res = await fetch(`${process.env.GO_API_URL}/todos`, {
        headers: { Authorization: `Bearer ${token}` }
    })

    return Response.json(await res.json())
}
```

**Protecting routes with middleware:**

```ts
// middleware.ts
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
    const token = request.cookies.get('auth_token')

    if (!token) {
        return NextResponse.redirect(new URL('/login', request.url))
    }
}

export const config = {
    matcher: ['/todos/:path*']
}
```

### Expected Behavior

```bash
# Visit /todos without being logged in → redirected to /login
# Submit login form → cookie set, redirected to /todos
# Todos load — Go API receives valid Bearer token from BFF
# Logout → cookie cleared, redirected to /login
# Each user only sees their own todos
```

### Stretch Goals

- Add a register page and wire it to `POST /auth/register`
- Show the logged-in user's email in the UI header (decode it from the JWT on the server — never send the raw token to the client)
- Handle token expiry — if the Go API returns `401`, clear the cookie and redirect to `/login`

---

## Tips & Resources

### Go
- **Project layout:** `main.go`, `handlers.go`, `repository.go`, `middleware.go` is plenty for these challenges
- **Testing your API:** [httpie](https://httpie.io/) (`http POST :8080/todos title="test"`) is friendlier than curl
- **Go module setup:** `go mod init github.com/yourname/go-challenges && go mod tidy`
- **Recommended packages:**
  - Gin: `github.com/gin-gonic/gin`
  - CORS: `github.com/gin-contrib/cors`
  - SQLite driver: `github.com/mattn/go-sqlite3`
  - JWT: `github.com/golang-jwt/jwt/v5`
  - Bcrypt: `golang.org/x/crypto/bcrypt`
  - Rate limiter: `golang.org/x/time/rate`

### React / Next.js
- **Vite docs:** https://vitejs.dev
- **TanStack Query docs:** https://tanstack.com/query/latest
- **Next.js App Router docs:** https://nextjs.org/docs/app
- **Key mental shift for backend devs:** You're not writing procedures — you're describing what the UI looks like given the current state. React handles the rest.
- **Server vs. Client Components:** If it needs `onClick`, `useState`, or browser APIs → `'use client'`. Everything else can stay a Server Component.
