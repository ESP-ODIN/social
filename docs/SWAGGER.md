# Swagger / OpenAPI

The Social API is documented with [swag](https://github.com/swaggo/swag): the OpenAPI 2.0 spec is
generated from comment annotations on the handlers, committed in this folder, and served by the API
with [http-swagger](https://github.com/swaggo/http-swagger).

| File | Role |
|---|---|
| `docs.go` | Generated. Registers the spec at startup (imported by `internal/router/swagger_routes.go`). |
| `swagger.json`, `swagger.yaml` | Generated. The same spec, for Postman, the frontend, or other services. |
| `SWAGGER.md` | This guide. |

Do not edit the generated files by hand; they are overwritten on every generation.

## Open the UI

| Run mode | URL |
|---|---|
| `go run ./cmd/api` / `make run` | http://127.0.0.1:8080/swagger/index.html |
| `make docker-dev` | http://127.0.0.1:8090/swagger/index.html |

The raw spec is at `/swagger/doc.json`. The Swagger routes are public, like `/health`.

## Call protected routes

Routes marked with a padlock (for example `POST /api/v1/posts`) require a JWT signed by Auth with
the same `JWT_SECRET_KEY` as this API. The token must use HS256 and contain an `exp` claim and
an `id` claim (the user UUID).

1. Get a token from Auth's login route.
2. Click **Authorize** and enter `Bearer <token>`. The `Bearer ` prefix must be typed: with
   Swagger 2.0 the UI sends the value as-is.
3. Click **Try it out** on a route, then **Execute**.

Requests sent from the UI hit the real database configured in `DATABASE_URL`.

## Document a new route

1. Register the route in `internal/router` as usual.
2. Annotate the handler, right above its declaration:

   ```go
   // GetPost godoc
   //
   //	@Summary	Get a post
   //	@Tags		posts
   //	@Produce	json
   //	@Security	BearerAuth
   //	@Param		id	path		string	true	"Post ID"	format(uuid)
   //	@Success	200	{object}	model.Post
   //	@Failure	401	{string}	string	"unauthorized"
   //	@Failure	404	{object}	dto.ErrorResponse
   //	@Router		/api/v1/posts/{id} [get]
   func GetPost(posts service.PostService) http.HandlerFunc {
   ```

3. Regenerate the spec (see [below](#regenerate-the-spec)) and commit `docs/` with the route.

`internal/handler/post.go` (`CreatePost`) is a complete example with a request body, a response
header and error responses.

### Annotation reference

| Annotation | Meaning |
|---|---|
| `@Summary` | One-line title shown in the UI. |
| `@Description` | Longer text shown when the route is expanded. |
| `@Tags` | Groups routes in the UI (`health`, `posts`, ...). |
| `@Accept json` | Request body content type. Only for routes with a body. |
| `@Produce json` | Response content type. |
| `@Security BearerAuth` | The route requires a JWT. Add it to every route behind `middleware.Auth`. |
| `@Param name in type required "description"` | A parameter. `in` is `path`, `query`, `header`, or `body`. For `body`, `type` is a struct such as `dto.CreatePostRequest`. |
| `@Success code {object} Type` | Successful response and its body type. |
| `@Failure code {object} Type` | Error response. Use `dto.ErrorResponse` for errors written as `{"error": "..."}`. |
| `@Header code {string} Name "description"` | A response header, such as `Location`. |
| `@Router /path [method]` | Path and method. It must match the route in `internal/router`: swag does not read the router. |

Formatting rules:

- Keep the `// FuncName godoc` line and the empty `//` line, and indent the annotations with a tab
  after `//`, so `gofmt` does not reflow them.
- Types are referenced as `package.Type`, and the handler file must import that package.
- Errors from `middleware.Auth` are plain text (`unauthorized`), hence `{string} string` on `401`.

The full syntax is in the
[swag documentation](https://github.com/swaggo/swag#declarative-comments-format).

## Document request and response types

swag builds schemas from the Go structs in `internal/dto` and `internal/model`. These struct tags
only affect the documentation, not runtime behavior:

| Tag | Effect | Example |
|---|---|---|
| `example:"..."` | Sample value, used to prefill **Try it out**. | `example:"Hello"` |
| `validate:"required"` | Marks the field as required. It does not validate anything: checks stay in the service. | `validate:"required"` |
| `maxLength:"N"`, `minLength:"N"` | String length limits. | `maxLength:"255"` |
| `format:"..."` | String format. | `format:"uuid"`, `format:"date-time"` |

Always add `format:"date-time"` to `time.Time` fields: swag does not add it automatically.

Handlers should encode the same types the annotations reference (for example `dto.StatusResponse`
and `dto.ErrorResponse`), so the documentation cannot drift from the real responses.

## Regenerate the spec

```sh
make swagger
```

It runs `go tool swag init -g cmd/api/main.go -o docs --parseInternal`:

- `-g cmd/api/main.go`: the file holding the general API info (title, version, `BearerAuth`
  scheme), in the comments above `main()`.
- `-o docs`: output folder.
- `--parseInternal`: required, because handlers and types live under `internal/`.

`swag` is pinned as a Go tool in `go.mod` (`tool github.com/swaggo/swag/cmd/swag`), so no global
install is needed. To upgrade it: `go get -tool github.com/swaggo/swag/cmd/swag@latest`.

With `make docker-dev`, air regenerates the spec before each rebuild (see `.air.toml`). `docs/` is
excluded from air's watched folders, otherwise writing `docs.go` would trigger an endless rebuild
loop. If an annotation is invalid, air shows the swag error and keeps the previous build running.

## Troubleshooting

| Symptom | Cause / fix |
|---|---|
| A route is missing from the UI, or shows old data | The spec was not regenerated: run `make swagger`. `TestSwaggerRoutes` in `internal/router` also fails when a documented route is missing. |
| `401` from **Try it out** | Missing `Bearer ` prefix, expired token, or `JWT_SECRET_KEY` different from Auth's. |
| `cannot find type definition: dto.X` | The handler file does not import the `dto` package, or `--parseInternal` is missing. |
| `warning: failed to get package name in dir: ./` | Harmless: the repository root has no Go files. |
| The API exits with `ping database: context deadline exceeded` | Not Swagger-related: the Neon compute was suspended and took too long to wake up. Start the API again. |
