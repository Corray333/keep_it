package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

type categoryCreater interface {
	CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error)
}

func (s *NoteService) CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	return s.categoryCreater.CreateCategory(ctx, category)
}
