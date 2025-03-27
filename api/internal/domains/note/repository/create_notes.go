package repository

import (
	"context"
	"fmt"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/google/uuid"
)

func (r *NoteRepository) CreateNote(ctx context.Context, note *entities.Note) (noteID uuid.UUID, err error) {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	if isNew {
		defer tx.Rollback()
	}

	fmt.Println("Created at: ", note.CreatedAt)

	if err := tx.QueryRow("INSERT INTO notes (creator_id, title, source, original, created_at, type, category_id, content, icon, cover) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING note_id", note.CreatorID, note.Title, note.Source, note.Original, note.CreatedAt, note.Type, note.CategoryId, note.Content, note.Icon, note.Cover).Scan(&noteID); err != nil {
		return uuid.Nil, err
	}

	// TODO: insert tags and access

	if isNew {
		if err := tx.Commit(); err != nil {
			return uuid.Nil, err
		}
	}

	return noteID, nil
}
