package implementation

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"social/internal/repository"
	"social/internal/service"
)

func (s *postService) Delete(ctx context.Context, id string, authorID uuid.UUID) error {
	if authorID == uuid.Nil {
		return service.ErrPostForbidden
	}
	post, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if post.AuthorID != authorID.String() {
		return service.ErrPostForbidden
	}
	// Include the author in the DELETE predicate as well as the permission check.
	err = s.repo.Delete(ctx, post.ID, authorID.String())
	if errors.Is(err, repository.ErrPostNotFound) {
		return service.ErrPostNotFound
	}
	return err
}
