package repository

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/Corray333/keep_it/internal/storage"
)

type NoteRepository struct {
	*storage.Storage
}

func New(store *storage.Storage) *NoteRepository {
	return &NoteRepository{
		Storage: store,
	}
}

func (r *NoteRepository) CreateNote(ctx context.Context, note *entities.Note) (noteID string, err error) {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return "", err
	}
	if isNew {
		defer tx.Rollback()
	}

	if err := tx.QueryRow("INSERT INTO notes (creator, title, source, original, font, created_at, type, category_owner, category_id, content, icon) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING note_id", note.Creator, note.Title, note.Source, note.OriginalRaw, note.Font, note.CreatedAt, note.Type, note.CategoryOwner, note.CategoryId, note.ContentRaw, note.IconRaw).Scan(&noteID); err != nil {
		return "", err
	}

	if _, err := tx.Exec("INSERT INTO user_note_access VALUES($1, $2)", note.Creator, noteID); err != nil {
		return "", err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return "", err
		}
	}

	return noteID, nil
}
