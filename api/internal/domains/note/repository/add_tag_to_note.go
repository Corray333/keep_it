package repository

import (
	"context"
	"fmt"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/google/uuid"
)

func (r *NoteRepository) AddTagToNote(ctx context.Context, tag *entities.Tag, noteID uuid.UUID) error {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	fmt.Println("noteID: ", noteID)
	fmt.Printf("tag: %v\n", tag)

	if _, err := tx.Exec("INSERT INTO note_tag (note_id, tag_text, owner_id) VALUES ($1, $2, $3)", noteID, tag.Text, tag.Owner); err != nil {
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
