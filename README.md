# Ticket System API

A REST API in Go where users register, log in, create tickets, and view or update the status of their own tickets.

## Live Deployment

| | URL |
|---|---|
| Base URL | `https://ticket-system-afd9.onrender.com` |
| Health check | `https://ticket-system-afd9.onrender.com/health/health` |

## Tech Stack

- Go
- JWT authentication (HS256)
- bcrypt password hashing
- In-memory storage
- Docker

## Run Locally

**With Go**

```bash
go mod download
go run .
```

**With Docker**

```bash
docker build -t ticket-system .
docker run -p 8080:8080 ticket-system
curl http://localhost:8080/health
```

The server runs on port `8080`.

## Environment Variables

See `.env.example`.

| Variable     | Default   | Description               |
|--------------|-----------|---------------------------|
| `PORT`       | `8080`    | HTTP port                 |
| `JWT_SECRET` | dev value | Secret used to sign JWTs  |

## API

Protected endpoints require `Authorization: Bearer <token>`.

| Method | Endpoint               | Auth | Description              |
|--------|------------------------|------|--------------------------|
| GET    | `/health`              | No   | Health check             |
| POST   | `/auth/register`       | No   | Register a user          |
| POST   | `/auth/login`          | No   | Log in and receive a JWT |
| POST   | `/tickets`             | Yes  | Create a ticket          |
| GET    | `/tickets`             | Yes  | List your tickets        |
| GET    | `/tickets/{id}`        | Yes  | Get one of your tickets  |
| PATCH  | `/tickets/{id}/status` | Yes  | Update a ticket's status |

### Request and response bodies

**`POST /auth/register`** → `201`

```json
{ "email": "user@example.com", "password": "secret123" }
```

**`POST /auth/login`** → `200`

```json
{ "email": "user@example.com", "password": "secret123" }
```

```json
{ "token": "<jwt>", "token_type": "Bearer" }
```

**`POST /tickets`** → `201`

```json
{ "title": "Login bug", "description": "500 error on login page" }
```

**`PATCH /tickets/{id}/status`** → `200`

```json
{ "status": "in_progress" }
```

**Ticket object**

```json
{
  "id": "d6b7b98b-692e-461e-a389-4804d6058cbd",
  "user_id": "65b9b2d2-d917-48fe-b7d2-c3022478b84f",
  "title": "Login bug",
  "description": "500 error on login page",
  "status": "open",
  "created_at": "2026-09-16T11:48:14Z",
  "updated_at": "2026-09-16T11:48:14Z"
}
```

Errors are returned as:

```json
{ "error": "message" }
```

| Status | Meaning |
|--------|---------|
| 400 | Invalid input |
| 401 | Missing or invalid token, or wrong credentials |
| 404 | Ticket not found or not owned by the user |
| 405 | Wrong HTTP method |
| 409 | Email already registered, or status change not allowed |

## Ticket Status Flow

```
open → in_progress → closed
```

- Status can only move forward.
- A `closed` ticket cannot be changed.

## Assumptions

- Data is stored in memory and resets when the service restarts.
- User and ticket IDs are UUIDs.
- Requesting another user's ticket returns `404`, so the API does not reveal which ticket IDs exist.
- Email is case-insensitive. Password must be 6–72 characters.
- `title` is required and `description` is optional. New tickets start as `open`.
- `open → closed` is allowed directly.
- JWTs expire after 24 hours.
