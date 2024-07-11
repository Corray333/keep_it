package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/note/types"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

type Storage interface {
	GetNote(note_id string) (*types.Note, error)
	CreateNote(note *types.Note) (string, error)
}

func GetNote(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		note, err := store.GetNote(chi.URLParam(r, "note_id"))
		if err != nil {
			slog.Error("error while getting note: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(note); err != nil {
			slog.Error("error while encoding response: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

type NewNoteResponse struct {
	NoteID string `json:"note_id"`
}

func NewNote(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		note := &types.Note{}
		if err := json.NewDecoder(r.Body).Decode(note); err != nil {
			slog.Error("error while decoding request body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		creds, err := auth.ExtractCredentials(r.Header.Get("Authorization"))
		if err != nil {
			slog.Error("error while extracting credentials: " + err.Error())
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
