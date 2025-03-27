package repository

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

func (r *NoteRepository) CreateTag(ctx context.Context, tag *entities.Tag) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	if _, err := tx.Exec("INSERT INTO tags (owner_id, tag_text, tag_color) VALUES ($1, $2, $3)", tag.Owner, tag.Text, tag.Color); err != nil {
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
