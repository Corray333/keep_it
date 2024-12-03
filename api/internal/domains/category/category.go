package category

import (
	"github.com/Corray333/keep_it/internal/domains/category/repository"
	"github.com/Corray333/keep_it/internal/domains/category/service"
	"github.com/Corray333/keep_it/internal/domains/category/transport"
	user_service "github.com/Corray333/keep_it/internal/domains/user/service"
	"github.com/Corray333/keep_it/internal/storage"
	"github.com/go-chi/chi/v5"
)

type CategoryController struct {
	repo      repository.CategoryRepository
	service   service.CategoryService
	transport transport.CategoryTransport
}

func NewCategoryController(router *chi.Mux, store *storage.Storage, userService *user_service.UserService) *CategoryController {
	repo := repository.New(store)
	service := service.New(repo, userService)
	transport := transport.New(router, service)

	return &CategoryController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *CategoryController) Build() {
	c.transport.RegisterRoutes()
}

func (c *CategoryController) Run() {
	go c.service.Run()
	go c.transport.Run()
}
