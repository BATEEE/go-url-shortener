# URL Shortener Service

## System Description

A small URL shortener service that converts long URLs into short, shareable codes.

Main use cases:

- User submits a long URL → receives a short code
- Access short code → redirect to the original URL
- Track number of clicks
- Manage links created by a user

## How to run the project

### Prerequisites

- Go 1.25+
- SQLite3

### Clone & Run

```bash
git clone <repo-url>
cd simple-shortener
go mod download
go run main.go
```

Server runs at: http://localhost:8080

### Docker (Optional)

```bash
docker build -t url-shortener .
docker run -p 8080:8080 url-shortener
```

## API Documentation

1. Create User

```bash
POST /users
Content-Type: application/json

{
  "email": "user@example.com"
}
```

2. Create Short Link

```bash
POST /users/:user_id/shorten
Content-Type: application/json

{
  "url": "https://example.com/very/long/url",
  "short_code": "custom123"  # optional
}
```

3. Redirect

```bash
GET /:code
→ 302 Redirect to original URL
```

4. Link Info

```bash
GET /links/:code/info
```

5. List user links

```bash
GET /users/:user_id/links
```

## Design & Technical Decisions

### Architecture: Clean Architecture

handler (HTTP) → service (business logic) → repository (data access)

Reasons:

- Separates concerns, easier to test
- Swap DB or web framework easily
- Better maintainability and scalability

Database: SQLite + GORM

Why SQLite:

- Simple, no separate DB server required
- Sufficient for a demo / prototype
- File-based and easy to back up
- GORM provides auto-migration

Trade-offs:

- Not ideal for heavy production load (concurrent writes)
- Limited to single-file DB
- For production, consider PostgreSQL or MySQL

### API: REST (using Gin)

Why REST:

- Simple and familiar
- Easy to test with curl/Postman
- Good fit for CRUD-style operations

Short code generation:

Option A — Custom short code (user-provided)

- Validate length and allowed characters (3–32 chars, alphanumeric, dash, underscore)

Option B — Random 6-character code

- Alphabet: a-zA-Z0-9 (62 chars)
- 62^6 ≈ 56 billion combinations
- Retry up to 5 times on collision
- Database UNIQUE constraint prevents duplicates

Duplicate original URL handling:

- If user creates the same URL again, return the existing short code for that (user, original_url) pair to save storage and be consistent.

Performance optimizations (current):

- Index on short_code (unique) for fast lookup
- Composite index (user_id, original_url) for duplicate checks
- Increment clicks asynchronously (goroutine) for non-blocking redirect

Security considerations (implemented)

- Basic URL validation
- Input validation using Gin binding
- GORM parameterized queries to avoid SQL injection

Security (to add later for production)

- Rate limiting
- Authentication / Authorization
- HTTPS enforcement (TLS)
- CORS configuration

### User Authentication Required

Reasons:

- Reduce spam/abuse (anonymous users could create unlimited links)
- Users can manage their links (list, and later delete)
- Enables per-user tracking and analytics
- Aligns with a business model (plans and billing per user)

Trade-offs:

- Higher friction for first-time users
- Users must create an account before using the service

## Challenges & Solutions

### Challenge 1: Random code collision

Problem: Two concurrent requests may generate the same short code.  
Solution:

- Retry up to 5 times with a new random code
- Database UNIQUE constraint on short_code as a safety net
- Low probability of collision with 62^6 combinations

### Challenge 2: Async increment can miss clicks

Problem: A goroutine used to increment click counts may fail silently.  
Solution:

- Log errors for monitoring (e.g., log.Printf(...))
- Accept eventual consistency for click counts (non-blocking redirect)
- Design trade-off: do not block user-facing redirects for analytics updates

### Challenge 3: Handling duplicate original URLs

Problem: A user creating the same URL multiple times wastes storage and creates duplicate short codes.  
Solution:

- Check for an existing link using GetByOriginalURL(userID, originalURL)
- Return the existing short code if found
- Use a composite index (user_id, original_url) to make the lookup fast

## Database Schema

users:

- id BIGINT PK AUTOINCREMENT
- email TEXT UNIQUE NOT NULL
- created_at TIMESTAMP
- updated_at TIMESTAMP

links:

- id BIGINT PK AUTOINCREMENT
- user_id BIGINT NOT NULL (INDEX)
- short_code VARCHAR(32) UNIQUE NOT NULL
- original_url TEXT NOT NULL (composite index with user_id)
- clicks INT DEFAULT 0
- created_at TIMESTAMP
- updated_at TIMESTAMP

INDEX idx_user_url ON links(user_id, original_url)

## What is currently missing

High priority

- Unit tests for handlers, services, and repositories (write tests for core flows: create link, redirect, duplicate handling)
- Basic CI: GitHub Actions that run go test and go vet on push/PR
- Dockerfile already present? If not, add a simple Dockerfile and a short README section for Docker
- Improved error logs and structured logging (use log package or a small logging helper)
- Input validation tests and API examples (curl examples)
- Add a README section about environment variables (DB path, port)

Medium priority

- Simple rate limiting middleware (per-IP, e.g., allow X requests per minute) — can use a lightweight library or an in-memory token bucket for the prototype
- Link expiration flag and a small background cleanup job (run via a goroutine + time.Ticker) — simpler than a full job queue
- Simple authentication for the management APIs (API key or a development JWT) — start with a very small middleware and document it
- Better URL validation (check scheme http/https, optional DNS lookup)

Lower priority

- Basic HTML admin page to list user links (static template) — useful to demo without building a full frontend
- QR code generation (use a simple Go library) for each short link

## If there is more time (next-phase, still intern-friendly)

- Add integration tests that run against a temporary SQLite file
- Create a GitHub Actions workflow to run lint, tests, and build a Docker image
- Add a simple healthcheck endpoint (GET /health) and basic readiness check
- Implement a background job for stats aggregation (daily click totals) stored in a small table
- Add basic monitoring by shipping logs to stdout in structured JSON (easy to read and compatible with container platforms)
- Add simple caching in-memory (LRU cache) to reduce DB reads for very hot links — implement with a small library or a simple map+mutex with TTL
- Provide a migration guide for moving from SQLite to PostgreSQL (document steps and environment variables)

## Production checklist

1. Tests & CI

   - Add unit tests for core logic
   - Add integration tests that run with a temp SQLite file
   - Add GitHub Actions to run tests on push/PR

2. Configuration & deployment basics

   - Move configuration to environment variables (DB path, port, allowed origins, admin API key)
   - Provide a production-ready Dockerfile
   - Add a simple Deploy README with steps to run with Docker Compose or a single container

3. Reliability & observability

   - Improve logging (structured logs, log levels)
   - Add a /health endpoint
   - Ensure DB file is stored in a configurable location and backed up (regular file copy or sqlite3 .dump script cron)

4. Security & stability

   - Add basic auth for link management APIs (API key or JWT)
   - Add simple rate limiting middleware
   - Enforce https at deployment (document using a reverse proxy such as nginx or a managed platform that provides TLS)

5. Operational runbook (document)
   - How to start/stop the service
   - How to restore from backup
   - How to run migrations
   - How to rotate admin API keys

## Limitations & Future Improvements

- Authentication (API key or JWT)
- Rate limiting
- Link expiration and cleanup job
- Unit and integration tests (high priority)
- Basic API docs and examples (curl)
- Simple admin UI (optional)

## Configuration

The application uses environment variables for configuration:

```bash
# .env
BASE_URL=http://localhost:8080
PORT=:8080
DB_PATH=./shortener.db
```
