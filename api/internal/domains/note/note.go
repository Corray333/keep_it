package note

import (
	storage "github.com/Corray333/keep_it/internal/domains/note/repository"
	"github.com/Corray333/keep_it/internal/domains/note/service"
	"github.com/Corray333/keep_it/internal/domains/note/transport"
	"github.com/Corray333/keep_it/internal/global_storage"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	store global_storage.Storage
}

func NewController() *Controller {
	return &Controller{}
}

func (c Controller) Init(router *chi.Mux, storeGlobal global_storage.Storage) {

	store := storage.NewStorage(storeGlobal.GetDB(), storeGlobal.GetRedis())

	service := service.NewService(store)
	server := transport.NewServer(service)

	router.With(auth.NewAuthMiddleware()).Post("/api/tags", server.CreateTag())
	router.With(auth.NewAuthMiddleware()).Post("/api/notes", server.CreateNote())
	router.With(auth.NewAuthMiddleware()).Get("/api/notes", server.ListNotes())
	router.With(auth.NewAuthMiddleware()).Patch("/api/notes/{note_id}", server.UpdateNote())
	router.With(auth.NewAuthMiddleware()).Get("/api/notes/{note_id}", server.GetNote())
}
