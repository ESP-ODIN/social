# Social

The social networking component of the ODIN application. Requires Go 1.27.0 or later.

## Run locally

On Windows (PowerShell):

```powershell
go run ./cmd/api
```

In a second terminal:

```powershell
Invoke-RestMethod http://127.0.0.1:8080/health
```

On Linux or macOS (Bash/Zsh), run the same Go command and check the endpoint with:

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

`.env.exemple` documents the environment variables. `.env` files are not loaded automatically.
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

## Structure

- `cmd/api`: HTTP server and routes.
- `config`: configuration loading and validation.
- `db`: reserved for data access.
- `cmd/migrate`: reserved for migrations; it returns an explicit error until the database is configured.
