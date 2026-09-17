# Social

The social networking component of the ODIN application. Requires Go 1.27.0 or later.

## Run locally

`DATABASE_URL` must be set (see [Configuration](#configuration)); the API pings the database on startup.

On Windows (PowerShell):

```powershell
$env:DATABASE_URL = "postgresql://user:password@host/dbname?sslmode=require"
go run ./cmd/api
```

In a second terminal:

```powershell
Invoke-RestMethod http://127.0.0.1:8080/health
```

On Linux or macOS (Bash/Zsh):

```sh
export DATABASE_URL=postgresql://user:password@host/dbname?sslmode=require
go run ./cmd/api
```

Check the endpoint with:

```sh
curl http://127.0.0.1:8080/health
```

## Configuration

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
`.env.exemple` documents the environment variables. `.env` files are not loaded automatically
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
meant to be plugged as-is into the project's docker-compose. Fill in `DATABASE_URL` in `.env` first
(see `.env.exemple`).

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

- `cmd/api`: HTTP server and routes.
- `config`: configuration loading and validation.
- `db`: reserved for data access.
- `cmd/migrate`: reserved for migrations; it returns an explicit error until the database is configured.
