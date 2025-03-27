package repository

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

func (r *NoteRepository) CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return nil, err
	}
	if isNew {
		defer tx.Rollback()
	}
	if err := tx.QueryRow("INSERT INTO categories (owner_id, name, parent_category_id) VALUES ($1, $2, $3) RETURNING category_id", category.OwnerID, category.Name, category.ParentCategoryID).Scan(&category.ID); err != nil {
		return nil, err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
	}

	return category, nil
}
