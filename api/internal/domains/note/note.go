package note

import (
	"github.com/Corray333/keep_it/internal/domains/note/repository"
	"github.com/Corray333/keep_it/internal/domains/note/service"
	"github.com/Corray333/keep_it/internal/domains/note/transport"
	user_service "github.com/Corray333/keep_it/internal/domains/user/service"
	"github.com/Corray333/keep_it/internal/storage"
	"github.com/go-chi/chi/v5"
)

type NoteController struct {
	repo      repository.NoteRepository
	service   service.NoteService
	transport transport.NoteTransport
}

func NewNoteController(router *chi.Mux, store *storage.Storage, userService *user_service.UserService) *NoteController {
	repo := repository.New(store)
	service := service.New(repo, userService)
	transport := transport.New(router, service)

	return &NoteController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *NoteController) Build() {
	c.transport.RegisterRoutes()
}

func (c *NoteController) Run() {
	go c.service.Run()
	go c.transport.Run()
}
