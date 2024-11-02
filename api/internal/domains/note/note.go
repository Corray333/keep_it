package note

import (
	"github.com/Corray333/keep_it/internal/domains/note/repository"
	"github.com/Corray333/keep_it/internal/domains/note/service"
	"github.com/Corray333/keep_it/internal/domains/note/transport"
	"github.com/Corray333/keep_it/internal/storage"
	"github.com/go-chi/chi/v5"
)

type NoteController struct {
	repo      repository.NoteRepository
	service   service.NoteService
	transport transport.NoteTransport
}

func NewNoteController(router *chi.Mux, store *storage.Storage) *NoteController {
	repo := repository.New(store)
	service := service.New(repo)
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
	c.service.Run()
}
