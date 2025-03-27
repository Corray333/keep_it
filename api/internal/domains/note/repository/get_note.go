package repository

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/google/uuid"
)

func (r *NoteRepository) GetNote(ctx context.Context, noteID uuid.UUID) (*entities.Note, error) {
	note := &entities.Note{}
	if err := r.DB.Get(note, "SELECT * FROM notes WHERE note_id = $1", noteID); err != nil {
		return nil, err
	}

	// TODO: get category and tags

	return note, nil
}
