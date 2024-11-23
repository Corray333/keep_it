package user

import (
	"github.com/Corray333/keep_it/internal/domains/user/repository"
	"github.com/Corray333/keep_it/internal/domains/user/service"
	"github.com/Corray333/keep_it/internal/domains/user/transport"
	"github.com/Corray333/keep_it/internal/storage"
	"github.com/go-chi/chi/v5"
)

type UserController struct {
	repo      repository.UserRepository
	service   service.UserService
	transport transport.UserTransport
}

func NewUserController(router *chi.Mux, store *storage.Storage) *UserController {
	repo := repository.New(store)
	service := service.New(repo)
	transport := transport.New(router, service)

	return &UserController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *UserController) Build() {
	c.transport.RegisterRoutes()
}

func (c *UserController) Run() {
	c.service.Run()
}

func (c *UserController) GetService() *service.UserService {
	return &c.service
}
