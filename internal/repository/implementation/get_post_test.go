package implementation

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"social/internal/model"
	"social/internal/repository"
)

type fakeRow struct {
	post model.Post
	err  error
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	values := []any{r.post.ID, r.post.AuthorID, r.post.Content, r.post.MediaURL, r.post.CreatedAt, r.post.UpdatedAt}
	for i, value := range values {
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(value))
	}
	return nil
}

type fakeQueryRower struct {
	row   fakeRow
	query string
	args  []any
	ctx   context.Context
}

func (q *fakeQueryRower) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	q.ctx, q.query, q.args = ctx, query, args
	return q.row
}

func TestGetByID(t *testing.T) {
	storageErr := errors.New("database unavailable")
	url := "https://example.com/image.jpg"
	for _, tc := range []struct {
		name        string
		media       *string
		input, want error
	}{
		{"without media", nil, nil, nil}, {"with media", &url, nil, nil},
		{"not found", nil, pgx.ErrNoRows, repository.ErrPostNotFound},
		{"storage error", nil, storageErr, storageErr},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want := model.Post{ID: "post-id", AuthorID: "author-id", Content: "hello", MediaURL: tc.media, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			db := &fakeQueryRower{row: fakeRow{post: want, err: tc.input}}
			ctx := context.Background()
			post, err := NewPostRepository(db).GetByID(ctx, want.ID)
			if !errors.Is(err, tc.want) || (err == nil && !reflect.DeepEqual(post, want)) {
				t.Fatalf("post=%+v, err=%v", post, err)
			}
			if db.ctx != ctx || !strings.Contains(db.query, "WHERE id = $1") || len(db.args) != 1 || db.args[0] != want.ID {
				t.Fatalf("unexpected query: %q, args=%v", db.query, db.args)
			}
		})
	}
}
