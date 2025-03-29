package repository

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (r *NoteRepository) SetCategory(ctx context.Context, userID int64, noteID uuid.UUID, categoryID string) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	if _, err := tx.Exec("UPDATE notes SET category_id = $1 WHERE note_id = $2 AND creator_id = $3", categoryID, noteID, userID); err != nil {
		slog.Error("Error setting category", "error", err)
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
