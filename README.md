# Social

The social networking component of the ODIN application. Requires Go 1.27.0 or later.

## Run locally

`DATABASE_URL` and `JWT_SECRET_KEY` must be set (see [Configuration](#configuration)); the API pings the database on startup.

On Windows (PowerShell):

```powershell
$env:DATABASE_URL = "postgresql://user:password@host/dbname?sslmode=require"
$env:JWT_SECRET_KEY = "<same secret as Auth>"
go run ./cmd/api
```

In a second terminal:

```powershell
Invoke-RestMethod http://127.0.0.1:8080/health
Invoke-RestMethod http://127.0.0.1:8080/ready
```

On Linux or macOS (Bash/Zsh):

```sh
export DATABASE_URL=postgresql://user:password@host/dbname?sslmode=require
export JWT_SECRET_KEY='<same secret as Auth>'
go run ./cmd/api
```

Check the endpoint with:

```sh
curl http://127.0.0.1:8080/health
curl http://127.0.0.1:8080/ready
```

## Configuration

`GET /health` reports that the HTTP server is responding. `GET /ready` checks
the PostgreSQL connection on every request and returns `503` if it is unavailable.

## Posts

`POST /api/v1/posts` requires a valid JWT signed by Auth using the same
`JWT_SECRET_KEY`. The JWT `id` claim becomes `post.author_id`. The route creates
a post in the `post` table and returns `201 Created`, the created post, and a
`Location` header.
Auth tokens may also include `email`; post creation does not use it.

```sh
curl -X POST http://127.0.0.1:8080/api/v1/posts \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -d '{"content":"Hello","media_url":"https://example.com/image.jpg"}'
```

`content` cannot be empty. `media_url` is optional and limited to 255 characters
by the SQL column. Do not include `author_id` in the request body; unknown JSON
fields are rejected.

By default, the server listens on `127.0.0.1:8080`. To change the port on Windows (PowerShell):

```powershell
$env:HTTP_ADDR = "127.0.0.1:8081"
go run ./cmd/api
```

On Linux or macOS (Bash/Zsh):

```sh
export HTTP_ADDR=127.0.0.1:8081
go run ./cmd/api
```

`DATABASE_URL` is required (the API connects to Postgres and pings it on startup).
`JWT_SECRET_KEY` is also required and must match the secret used by Auth.
`.env.example` documents the environment variables. `.env` files are not loaded automatically
by `go run`/`go build` — export them yourself, or use `make docker-dev` below, which reads `.env`.
Configuration is read from the process environment.

## Checks and build

```sh
go test ./...
go vet ./...
go build -o bin/ ./cmd/...
```

The API executable is `bin/api.exe` on Windows and `bin/api` on Linux and macOS.
If Make is installed, `make run`, `make test`, `make vet`,
`make build`, and `make fmt` are also available.

## Docker

A development `Dockerfile` runs the API with hot-reload (via [air](https://github.com/air-verse/air)),
meant to be plugged as-is into the project's docker-compose. Fill in `DATABASE_URL` and
`JWT_SECRET_KEY` in `.env` first
(see `.env.example`).

Build the image:

```sh
make docker-build
```

Run it locally with hot-reload (mounts the source code and exposes the app on port 8090):

```sh
make docker-dev
```

Then check it with `curl http://127.0.0.1:8090/health`. Editing any `.go` file rebuilds and restarts
the server automatically inside the container. The container itself listens on `0.0.0.0:8080`
(overridden from `.env`'s `HTTP_ADDR` by the `docker-dev` target); only the host-side port is `8090`.

## Structure

- `cmd/api`: HTTP server startup.
- `internal/router`: route registration.
- `internal/handler`: HTTP handlers for health, readiness, and posts.
- `internal/dto`, `internal/model`, `internal/service`, `internal/repository`: post data and creation.
- `internal/repository/implementation`: PostgreSQL post repository.
- `internal/service/implementation`: post service logic.
- `config`: configuration loading and validation.
- `db`: reserved for data access.
- `cmd/migrate`: reserved for migrations; it returns an explicit error until the database is configured.
