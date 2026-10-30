package implementation

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"social/internal/repository"
)

type deleteRow struct{ err error }

func (r deleteRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	*dest[0].(*string) = "post-id"
	return nil
}

type deleteQueryRower struct {
	row   deleteRow
	query string
	args  []any
	ctx   context.Context
}

func (q *deleteQueryRower) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	q.ctx, q.query, q.args = ctx, query, args
	return q.row
}

func TestDelete(t *testing.T) {
	storageErr := errors.New("database unavailable")
	for _, tc := range []struct {
		name        string
		input, want error
	}{
		{"deleted", nil, nil}, {"no matching row", pgx.ErrNoRows, repository.ErrPostNotFound}, {"storage failure", storageErr, storageErr},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := &deleteQueryRower{row: deleteRow{err: tc.input}}
			ctx := context.Background()
			err := NewPostRepository(db).Delete(ctx, "post-id", "author-id")
			if !errors.Is(err, tc.want) {
				t.Fatalf("err=%v", err)
			}
			if db.ctx != ctx || !strings.Contains(db.query, "WHERE id = $1 AND author_id = $2") || len(db.args) != 2 || db.args[0] != "post-id" || db.args[1] != "author-id" {
				t.Fatalf("unexpected query: %q, args=%v", db.query, db.args)
			}
		})
	}
}
