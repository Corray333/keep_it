package service

import (
	"context"

	"github.com/google/uuid"
)

type categorySetter interface {
	SetCategory(ctx context.Context, userID int64, noteID uuid.UUID, categoryID string) error
}

func (c *NoteService) SetCategory(ctx context.Context, userID int64, noteID uuid.UUID, categoryID string) error {
	return c.categorySetter.SetCategory(ctx, userID, noteID, categoryID)
}
