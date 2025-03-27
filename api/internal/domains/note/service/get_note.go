package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/google/uuid"
)

type noteGetter interface {
	GetNote(ctx context.Context, noteID uuid.UUID) (*entities.Note, error)
}

func (c *NoteService) GetNoteByID(ctx context.Context, userID int64, noteID uuid.UUID) (*entities.Note, error) {
	note, err := c.noteGetter.GetNote(ctx, noteID)
	if err != nil {
		return nil, err
	}

	if note.CreatorID != userID {
		return nil, ErrNoAccess
	}

	if note.Tags == nil {
		note.Tags = []entities.Tag{}
	}

	return note, nil
}
