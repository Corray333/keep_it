package repository

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

func (r *NoteRepository) DeleteTag(ctx context.Context, tag *entities.Tag) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	if _, err := tx.Exec("DELETE FROM tags WHERE creator_id = $1 AND tag_text = $2", tag.Owner, tag.Text); err != nil {
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
