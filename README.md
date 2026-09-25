# Social

The social networking component of ODIN. Requires Go 1.27.0 or later.

## Authentication and local dependency

Auth alone signs RS256 JWTs. Social uses
[authkit-go](https://github.com/ESP-ODIN/authkit-go) to validate them through
Auth's public JWKS, including issuer, audience, expiry and issued-at.
Social has no private key or shared JWT secret.

The inspected GitHub commit `fb87d0a2296b` currently contains only a README,
with no Go package. Development therefore uses the existing sibling checkout:

```go
require github.com/ESP-ODIN/authkit-go v0.0.0-20260922093749-fb87d0a2296b
replace github.com/ESP-ODIN/authkit-go => ../authkit-go
```

Here that path resolves to `C:\Users\benis\IdeaProjects\authkit-go`.
Keep the repositories side by side. No library code is copied into Social.
The pseudo-version identifies the existing GitHub commit; no tag or release
was created. Once the Go package is available on GitHub, remove the replacement,
select the available commit with Go, and run `go mod tidy`.
Until then, a Social-only checkout/CI job needs the sibling library checkout.

One authenticator is initialized at startup and shared across requests so its
JWKS cache is reused. Only `POST /api/v1/posts` is protected.

Client -> Authorization: Bearer JWT -> authkit-go RequireAuth -> RS256/JWKS validation
-> Identity.ID -> Post Handler -> Post Service -> Post Repository -> post.author_id.

The JWT `id` is Auth's USER.id, never `sub`. No Auth database access or
cross-database foreign key is involved.

## Configuration

The API loads `.env` from its working directory using godotenv; existing process
environment variables take precedence. See [.env.example](.env.example).

| Variable | Purpose | Example for both services running on the host |
| --- | --- | --- |
| HTTP_ADDR | Social listener; code default is 127.0.0.1:8080 | 127.0.0.1:8081 |
| DATABASE_URL | Required Social PostgreSQL connection | postgresql://user:password@host/dbname?sslmode=require |
| AUTH_JWKS_URL | Required public JWKS endpoint reachable by Social | http://localhost:8080/.well-known/jwks.json |
| AUTH_ISSUER | Required exact JWT issuer | http://localhost:8080 |
| AUTH_AUDIENCE | Required JWT audience | odin-api |

The example uses 8081 for Social to leave 8080 for Auth. The API connects to and
pings its own database on startup. The existing `post` table must be provisioned;
`cmd/migrate` remains a placeholder.

`GET /health` remains public and returns 200.
`GET /ready` remains public, pings PostgreSQL, and returns 200 or 503.

## Run Auth and Social locally

In the sibling Auth checkout, configure its own `.env` according to its README:
Auth's DATABASE_URL, HTTP_ADDR=127.0.0.1:8080, JWT_PRIVATE_KEY_PATH pointing to its
existing RSA private key, JWT_KEY_ID, JWT_ISSUER=http://localhost:8080 and
JWT_AUDIENCE=odin-api. The private key stays in Auth.

In a PowerShell terminal:

```powershell
Set-Location C:\Users\benis\IdeaProjects\auth
go run ./cmd/api
```

In another terminal, prepare Social's `.env` from `.env.example` if it does not
already exist, and fill in Social's DATABASE_URL:

```powershell
Set-Location C:\Users\benis\IdeaProjects\social
if (-not (Test-Path .env)) { Copy-Item .env.example .env }
# Edit .env with the Social database URL and the three AUTH_* values above.
$env:HTTP_ADDR = '127.0.0.1:8081'
go run ./cmd/api
```

## Obtain a JWT and create a post

The inspected Auth source exposes `POST /api/v1/auth/register`; it has no login
route. Registration needs Auth's existing database schema and role ID 1.
This creates a new account and obtains its token from `data.token`:

```powershell
$email = 'social-' + [guid]::NewGuid().ToString('N') + '@example.com'
$registration = @{ email = $email; password = 'Local-test-password-123!' } | ConvertTo-Json
$auth = Invoke-RestMethod -Method Post -Uri 'http://localhost:8080/api/v1/auth/register' -ContentType 'application/json' -Body $registration
$token = $auth.data.token

$body = @{ content = 'Mon premier post ODIN'; media_url = 'https://example.com/image.jpg' } | ConvertTo-Json
$post = Invoke-RestMethod -Method Post -Uri 'http://localhost:8081/api/v1/posts' -Headers @{ Authorization = "Bearer $token" } -ContentType 'application/json' -Body $body
$post
if ($post.author_id -ne $auth.data.user.id) { throw 'Unexpected post author' }

Invoke-RestMethod http://localhost:8081/health
Invoke-RestMethod http://localhost:8081/ready
```

Successful creation returns 201, the created post and a Location header.
The body accepts `content` and optional `media_url` (at most 255 characters).
Content cannot be blank. Unknown fields, including `author_id`, are rejected.
Invalid JSON or business data returns 400; the existing 1 MiB body limit returns
413. Missing or invalid authentication returns a generic 401; unexpected
repository errors return a generic 500 without SQL details.

## API documentation (Swagger)

The Swagger UI is served at `http://127.0.0.1:8080/swagger/index.html` (port `8090` with
`make docker-dev`). The spec is generated from annotations on the handlers: after adding or
changing a route, annotate its handler and run `make swagger`.

See [docs/SWAGGER.md](docs/SWAGGER.md) for authentication in the UI, how to document a new route,
and troubleshooting.

## Checks and build

```sh
gofmt -w cmd config db internal
go mod tidy
go test ./...
go vet ./...
go test -race ./...
go build -o bin/ ./cmd/...
```

The API executable is `bin/api.exe` on Windows and `bin/api` on Linux and macOS.
If Make is installed, `make run`, `make test`, `make vet`,
`make build`, `make fmt`, and `make swagger` are also available.
Tests use ephemeral RSA keys, an httptest JWKS server and fake services or
repositories; no running Auth, Internet or Neon is required after dependencies
are available. The race detector requires CGO and a compatible C compiler.
The API executable is bin/api.exe on Windows, bin/api on Linux/macOS.
Make targets run, build, test, vet and fmt are also available.

## Docker

The current development Dockerfile uses Air. `make docker-dev` publishes Social
on host port 8090 and listens on 0.0.0.0:8080 inside the container.

No shared compose/network is defined in the inspected Social or Auth repositories.
Inside Docker, localhost refers to the Social container, not Auth on the host.
Set AUTH_JWKS_URL to the address actually reachable in your Docker setup.
AUTH_ISSUER must still match the token's issuer exactly, even if the JWKS
transport address differs. No Auth service hostname is assumed here.

The temporary ../authkit-go replacement is outside the Dockerfile's Social-only
build context and is not included by the current runtime bind mount. Therefore
the current `make docker-build` / `make docker-dev` flow cannot use this local
replacement as-is. Prefer running both services directly during this stage.
Restore the ordinary GitHub dependency once its Go code is available before
using the existing Docker flow, or explicitly arrange a context/mount containing
both repositories in your infrastructure. No new compose infrastructure is added.

## Structure

- `cmd/api`: HTTP server startup.
- `internal/router`: route registration.
- `internal/handler`: HTTP handlers for health, readiness, and posts.
- `internal/dto`, `internal/model`, `internal/service`, `internal/repository`: post data and creation.
- `internal/repository/implementation`: PostgreSQL post repository.
- `internal/service/implementation`: post service logic.
- `docs`: OpenAPI spec generated by `make swagger` (do not edit by hand) and the
  [Swagger guide](docs/SWAGGER.md).
- `config`: configuration loading and validation.
- `db`: reserved for data access.
- `cmd/migrate`: reserved for migrations; it returns an explicit error until the database is configured.
