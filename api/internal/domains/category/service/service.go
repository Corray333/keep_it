package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/category/entities"
	"github.com/Corray333/keep_it/internal/helpers"
)

var (
	ErrNoAccess = helpers.NewError(http.StatusForbidden, errors.New("no access"))
)

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error)
	GetCategories(ctx context.Context, userID int64) ([]entities.Category, error)
}

type userService interface {
}

type CategoryService struct {
	repo repository
	userService
}

func New(repo repository, userService userService) *CategoryService {
	s := &CategoryService{
		repo:        repo,
		userService: userService,
	}
	return s
}

func (s *CategoryService) Run() {
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	return s.repo.CreateCategory(ctx, category)
}

func (s *CategoryService) GetCategories(ctx context.Context, userID int64) ([]entities.Category, error) {
	// Получаем категории из репозитория
	categories, err := s.repo.GetCategories(ctx, userID)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]*entities.Category)
	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	// Создаем список корневых категорий
	var rootCategories []*entities.Category

	// Строим дерево
	for i := range categories {
		category := &categories[i]

		if category.ParentCategoryID == nil {
			// Если у категории нет родителя, добавляем ее в корневые
			rootCategories = append(rootCategories, category)
		} else {
			// Если есть родитель, добавляем в его дочерние категории
			parent, exists := categoryMap[*category.ParentCategoryID]
			if exists {
				parent.ChildrenCategories = append(parent.ChildrenCategories, *category)
			}
		}
	}

	result := []entities.Category{}
	for i := range rootCategories {
		result = append(result, *rootCategories[i])
	}

	return result, nil
}
