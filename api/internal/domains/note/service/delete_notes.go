package service

import (
	"context"

	"github.com/google/uuid"
)

type notesDeleter interface {
	DeleteNotes(ctx context.Context, userID int64, noteIDs []uuid.UUID) error
}

func (c *NoteService) DeleteNotes(ctx context.Context, useID int64, noteIDs []uuid.UUID) error {
	return c.notesDeleter.DeleteNotes(ctx, useID, noteIDs)
}
