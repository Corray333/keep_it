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

func (r *NoteRepository) GetNote(ctx context.Context, noteID string) (*entities.Note, error) {
	note := &entities.Note{}
	if err := r.DB.Get(note, "SELECT * FROM notes WHERE note_id = $1", noteID); err != nil {
		return nil, err
	}

	// TODO: get category and tags

	return note, nil
}

func (r *NoteRepository) CreateNote(ctx context.Context, note *entities.Note) (noteID string, err error) {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return "", err
	}
	if isNew {
		defer tx.Rollback()
	}

	if err := tx.QueryRow("INSERT INTO notes (creator_id, title, source, original, created_at, type, category_id, content, icon, cover) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING note_id", note.CreatorID, note.Title, note.Source, note.Original, note.CreatedAt, note.Type, note.CategoryId, note.ContentRaw, note.IconRaw, note.Cover).Scan(&noteID); err != nil {
		return "", err
	}

	// TODO: insert tags and access

	if isNew {
		if err := tx.Commit(); err != nil {
			return "", err
		}
	}

	return noteID, nil
}
