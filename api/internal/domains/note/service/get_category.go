package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

type categoryGetter interface {
	GetCategories(ctx context.Context, userID int64) ([]entities.Category, error)
}

func (s *NoteService) GetCategories(ctx context.Context, userID int64) ([]entities.Category, error) {
	// Получаем категории из репозитория
	categories, err := s.categoryGetter.GetCategories(ctx, userID)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string][]entities.Category)

	var result []entities.Category

	for i := range categories {
		if categories[i].ParentCategoryID == nil {
			result = append(result, categories[i])
		} else {
			parentCategoryChildren := categoryMap[*categories[i].ParentCategoryID]

			parentCategoryChildren = append(parentCategoryChildren, categories[i])
			categoryMap[*categories[i].ParentCategoryID] = parentCategoryChildren
		}
	}

	s.assignChildren(result, categoryMap)

	return result, nil
}

func (s *NoteService) assignChildren(categories []entities.Category, categoryMap map[string][]entities.Category) {
	for i := range categories {
		children := categoryMap[categories[i].ID]
		categories[i].ChildrenCategories = children
		s.assignChildren(children, categoryMap)
	}
}
