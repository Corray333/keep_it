package repository

import (
	"context"
	"database/sql"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

func (r *NoteRepository) GetCategories(ctx context.Context, userID int64) ([]entities.Category, error) {
	categories := []entities.Category{}
	if err := r.DB.Select(&categories, "SELECT * FROM categories WHERE owner_id = $1", userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return categories, nil
}
