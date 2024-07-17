package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/note/storage"
	"github.com/Corray333/keep_it/internal/domains/note/types"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

type Storage interface {
	GetNote(note_id string) (*types.Note, error)
	CreateNote(note *types.Note) (string, error)
	CheckNoteAccess(note_id string, user_id int) (bool, error)
	CreateTag(tag *types.Tag) (*types.Tag, error)
	UpdateNote(note_id string, data map[string]interface{}) error
}

// GetNote fetches a specific note.
// @Summary Fetch a note
// @Description Retrieves a note by its ID.
// @Tags notes
// @Accept json
// @Produce json
// @Param note_id path string true "Note ID"
// @Security ApiKeyAuth
// @Success 200 {object} types.Note
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Note not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes/{note_id} [get]
func GetNote(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		creds := r.Context().Value("creds").(auth.Credentials)

		if creds.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		note, err := store.GetNote(chi.URLParam(r, "note_id"))
		if err != nil {
			slog.Error("error while getting note: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if note.Creator != creds.ID {
			allowed, err := store.CheckNoteAccess(note.ID, creds.ID)
			if err != nil {
				slog.Error("error while checking note access: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !allowed {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(note); err != nil {
			slog.Error("error while encoding response: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// NewNoteResponse represents the response structure for creating a new note.
type NewNoteResponse struct {
	NoteID string `json:"note_id"`
}

// CreateNote creates a new note.
// @Summary Create a new note
// @Description Creates a new note with the provided details.
// @Tags notes
// @Accept json
// @Produce json
// @Param note body types.Note true "Note details"
// @Security ApiKeyAuth
// @Success 201 {object} NewNoteResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes [post]
func CreateNote(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		note := &types.Note{}
		if err := json.NewDecoder(r.Body).Decode(note); err != nil {
			slog.Error("error while decoding request body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		creds := r.Context().Value("creds").(auth.Credentials)
		if creds.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		note.Creator = creds.ID

		note_id, err := store.CreateNote(note)
		if err != nil {
			slog.Error("error while creating note: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(NewNoteResponse{NoteID: note_id}); err != nil {
			slog.Error("error while encoding response: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

// CreateTagRequest represents the request structure for creating a new tag.
type CreateTagRequest struct {
	Tag types.Tag `json:"tag"`
}

// CreateTagResponse represents the response structure for creating a new tag.
type CreateTagResponse struct {
	Tag types.Tag `json:"tag"`
}

// CreateTag creates a new tag.
// @Summary Create a new tag
// @Description Creates a new tag with the provided details.
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body CreateTagRequest true "Tag details"
// @Security ApiKeyAuth
// @Success 200 {object} CreateTagResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/tags [post]
func CreateTag(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := &types.Tag{}
		if err := json.NewDecoder(r.Body).Decode(tag); err != nil {
			slog.Error("error while decoding request body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		creds := r.Context().Value("creds").(auth.Credentials)
		if creds.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tag.Owner = creds.ID

		tag, err := store.CreateTag(tag)
		if err != nil {
			slog.Error("error while creating tag: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(CreateTagResponse{Tag: *tag}); err != nil {
			slog.Error("error while encoding response: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

type UpdateNoteRequest map[string]interface{}

// UpdateNote updates an existing note.
// @Summary Update a note
// @Description Updates an existing note with the provided fields.
// @Tags notes
// @Accept json
// @Produce json
// @Param note_id path string true "Note ID"
// @Security ApiKeyAuth
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Note not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes/{note_id} [patch]
func UpdateNote(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		note := UpdateNoteRequest{}
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			slog.Error("error while decoding request body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		creds := r.Context().Value("creds").(auth.Credentials)
		if creds.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		note_id := chi.URLParam(r, "note_id")
		if note_id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := store.UpdateNote(note_id, note); err != nil {
			slog.Error("error while updating note: " + err.Error())
			if err == storage.ErrorNoteDoesNotExist {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
