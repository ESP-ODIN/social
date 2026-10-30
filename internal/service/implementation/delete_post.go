package implementation

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"social/internal/service"
)

func (s *postService) Delete(ctx context.Context, id string, authorID uuid.UUID) error {
	parsed, err := uuid.Parse(id)
	if err != nil || parsed == uuid.Nil {
		return service.ErrInvalidPostID
	}
	if authorID == uuid.Nil {
		return service.ErrPostUnauthorized
	}
	err = s.repo.Delete(ctx, parsed.String(), authorID.String())
	if errors.Is(err, service.ErrPostNotFound) {
		// Only diagnose a failed deletion; the database enforces ownership atomically.
		if _, err := s.repo.GetByID(ctx, parsed.String()); err != nil {
			return err
		}
		return service.ErrPostForbidden
	}
	return err
}
