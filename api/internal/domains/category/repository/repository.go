package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Corray333/keep_it/internal/domains/category/entities"
	"github.com/Corray333/keep_it/internal/storage"
)

type CategoryRepository struct {
	*storage.Storage
}

func New(store *storage.Storage) *CategoryRepository {
	return &CategoryRepository{
		Storage: store,
	}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	tx, isNew, err := r.GetTx(ctx)
	if err != nil {
		return nil, err
	}
	if isNew {
		defer tx.Rollback()
	}
	fmt.Println("Parent: ", *category.ParentCategoryID)
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

func (r *CategoryRepository) GetCategories(ctx context.Context, userID int64) ([]entities.Category, error) {
	categories := []entities.Category{}
	if err := r.DB.Select(&categories, "SELECT * FROM categories WHERE owner_id = $1", userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return categories, nil
}
