# API Routes

All routes except health, readiness and Swagger require a valid JWT in the
`Authorization: Bearer <token>` header, validated through `authkit-go`.

## Public

| Method | Path | Description |
| --- | --- | --- |
| GET | `/health` | Health check — always returns 200 |
| GET | `/ready` | Readiness check — pings PostgreSQL, returns 200 or 503 |
| GET | `/swagger/` | Swagger UI and OpenAPI spec |

## Posts

| Method | Path | Description |
| --- | --- | --- |
| POST | `/api/v1/posts` | Create a new post (authenticated user becomes author) |
| GET | `/api/v1/posts/{id}` | Get a post by ID |
| GET | `/api/v1/users/{userId}/posts` | Get all posts by a user |
| DELETE | `/api/v1/posts/{id}` | Delete a post (author only) |

## Likes

| Method | Path | Description |
| --- | --- | --- |
| POST | `/api/v1/posts/{postId}/likes` | Like a post |
| GET | `/api/v1/posts/{postId}/likes` | Get all likes on a post |
| DELETE | `/api/v1/posts/{postId}/likes` | Unlike a post (removes authenticated user's like) |
