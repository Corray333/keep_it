package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

type notesGetter interface {
	fileGetter
	
	GetNotes(ctx context.Context, userID int64, offset int, filters entities.NoteFilter) ([]entities.Note, error)
}

func (c *NoteService) GetNotes(ctx context.Context, userID int64, offset int, filters entities.NoteFilter) ([]entities.Note, error) {
	notes, err := c.notesGetter.GetNotes(ctx, userID, offset, filters)
	if err != nil {
		return nil, err
	}

	for note
}
