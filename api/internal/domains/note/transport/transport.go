package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/Corray333/keep_it/internal/helpers"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

type NoteTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
	CreateNote(ctx context.Context, userID int64, note entities.Note) (noteID string, err error)
	GetNoteByID(ctx context.Context, userID int64, noteID string) (*entities.Note, error)
}

func New(router *chi.Mux, service service) *NoteTransport {
	return &NoteTransport{
		router:  router,
		service: service,
	}
}

func (t *NoteTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())

		r.Post("/api/notes", t.createNote)
		r.Get("/api/notes/{note_id}", t.getNote)
	})
}

type CreateNoteRequest struct {
	entities.Note
}

type CreateNoteResponse struct {
	NoteID string `json:"note_id"`
}

func (t *NoteTransport) createNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		slog.Error("user id not found in context")
		helpers.SendError(w, helpers.ErrInternal)
		return
	}

	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request: " + err.Error())
		helpers.SendError(w, err)
		return
	}

	noteID, err := t.service.CreateNote(ctx, userID, req.Note)
	if err != nil {
		slog.Error("failed to create note: " + err.Error())
		helpers.SendError(w, err)
		return
	}

	resp := CreateNoteResponse{
		NoteID: noteID,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode response: " + err.Error())
		helpers.SendError(w, err)
		return
	}
}

func (t *NoteTransport) getNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	noteID := chi.URLParam(r, "note_id")

	note, err := t.service.GetNoteByID(ctx, userID, noteID)
	if err != nil {
		slog.Error("failed to get note: " + err.Error())
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		slog.Error("failed to encode response: " + err.Error())
		helpers.SendError(w, err)
		return
	}
}
