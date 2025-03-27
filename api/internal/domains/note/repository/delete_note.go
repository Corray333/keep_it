package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *NoteRepository) DeleteNotes(ctx context.Context, userID int64, noteIDs []uuid.UUID) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}
	for _, noteID := range noteIDs {
		if _, err := tx.Exec("DELETE FROM notes WHERE creator_id = $1 AND note_id = $2", userID, noteID); err != nil {
			return err
		}
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
