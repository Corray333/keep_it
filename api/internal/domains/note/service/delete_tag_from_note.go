package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/google/uuid"
)

type tagFromNoteDeleter interface {
	DeleteTagFromNote(ctx context.Context, tag *entities.Tag, noteID uuid.UUID) error
}

func (c *NoteService) DeleteTagFromNote(ctx context.Context, tag *entities.Tag, noteID uuid.UUID) error {
	return c.tagFromNoteDeleter.DeleteTagFromNote(ctx, tag, noteID)
}
