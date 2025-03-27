package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

type tagsGetter interface {
	GetTags(ctx context.Context, userID int64) ([]entities.Tag, error)
}

func (c *NoteService) GetTags(ctx context.Context, userID int64) ([]entities.Tag, error) {
	return c.tagsGetter.GetTags(ctx, userID)
}
