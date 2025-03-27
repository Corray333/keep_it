package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

type notesGetter interface {
	GetNotes(ctx context.Context, userID int64, offset int, filters entities.NoteFilter) ([]entities.Note, error)
}

func (c *NoteService) GetNotes(ctx context.Context, userID int64, offset int, filters entities.NoteFilter) ([]entities.Note, error) {
	return c.notesGetter.GetNotes(ctx, userID, offset, filters)
}
