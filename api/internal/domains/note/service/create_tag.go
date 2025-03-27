package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

type tagCreater interface {
	CreateTag(ctx context.Context, tag *entities.Tag) error
}

func (c *NoteService) CreateTag(ctx context.Context, tag *entities.Tag) error {
	// TODO: add limit in 128 tags
	return c.tagCreater.CreateTag(ctx, tag)
}
