package repository

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

func (r *NoteRepository) GetTags(ctx context.Context, userID int64) ([]entities.Tag, error) {
	tags := []entities.Tag{}
	if err := r.DB.Select(&tags, "SELECT * FROM tags WHERE owner_id = $1", userID); err != nil {
		return nil, err
	}

	return tags, nil
}
