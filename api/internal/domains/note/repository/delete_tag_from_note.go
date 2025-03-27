package repository

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/google/uuid"
)

func (r *NoteRepository) DeleteTagFromNote(ctx context.Context, tag *entities.Tag, noteID uuid.UUID) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	if _, err := tx.Exec("DELETE FROM note_tag WHERE note_id = $1 AND tag_text = $2", noteID, tag.Text); err != nil {
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
