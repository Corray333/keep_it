package service

import (
	"context"
	"errors"
	"fmt"
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

	fmt.Printf("Root categories: %+v\n", result)

	return result, nil
}

func (s *CategoryService) assignChildren(categories []entities.Category, categoryMap map[string][]entities.Category) {
	for i := range categories {
		children := categoryMap[categories[i].ID]
		categories[i].ChildrenCategories = children
		s.assignChildren(children, categoryMap)
	}
}
