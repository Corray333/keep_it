package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

type tagDeleter interface {
	DeleteTag(ctx context.Context, tag *entities.Tag) error
}

func (c *NoteService) DeleteTag(ctx context.Context, tag *entities.Tag) error {
	return c.tagDeleter.DeleteTag(ctx, tag)
}
