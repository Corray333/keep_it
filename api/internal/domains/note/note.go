package note

import (
	"github.com/Corray333/keep_it/internal/domains/note/storage"
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

	router.With(auth.NewAuthMiddleware()).Post("/api/tags", transport.CreateTag(store))
	router.With(auth.NewAuthMiddleware()).Post("/api/notes", transport.CreateNote(store))
	router.With(auth.NewAuthMiddleware()).Get("/api/notes", transport.ListNotes(store))
	router.With(auth.NewAuthMiddleware()).Patch("/api/notes/{note_id}", transport.UpdateNote(store))
	router.With(auth.NewAuthMiddleware()).Get("/api/notes/{note_id}", transport.GetNote(store))
}
