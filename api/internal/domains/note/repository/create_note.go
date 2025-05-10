package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/google/uuid"
)

func (r *NoteRepository) CreateNote(ctx context.Context, note entities.Note) (noteID uuid.UUID, err error) {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	if isNew {
		defer tx.Rollback()
	}

	fmt.Println("Source: ", note.Source)
	if err := tx.QueryRow("INSERT INTO notes (creator_id, title, source, original, created_at, type, content, icon, cover) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING note_id", note.CreatorID, note.Title, note.Source, note.Original, note.CreatedAt, note.Type, note.Content, note.Icon, note.Cover).Scan(&noteID); err != nil {
		slog.Error("Error creating note", "error", err)
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
