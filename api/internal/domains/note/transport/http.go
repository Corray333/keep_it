package transport

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Corray333/keep_it/internal/domains/note/storage"
	"github.com/Corray333/keep_it/internal/domains/note/types"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

type service interface{}
type server struct {
	service service
}

type Storage interface {
	GetNote(note_id string) (*types.Note, error)
	CreateNote(note *types.Note) (string, error)
	CheckNoteAccess(note_id string, user_id int) (bool, error)
	CreateTag(tag *types.Tag) (*types.Tag, error)
	UpdateNote(note_id string, data map[string]interface{}) error
	GetNotes(user_id int, offset int, filter map[string]interface{}) ([]*types.Note, bool, error)
	DeleteNote(note_id string, uid int) error
}

type GetNoteResponse struct {
	ID            string      `json:"id" db:"note_id"`
	Creator       int         `json:"creator" db:"creator"`
	Tags          []types.Tag `json:"tags" db:"tags"`
	Title         string      `json:"title" db:"title"`
	Source        string      `json:"source" db:"source"`
	Original      any         `json:"original" db:"original"`
	Font          string      `json:"font" db:"font"`
	CreatedAt     int64       `json:"created_at" db:"created_at"`
	CopiedAt      int64       `json:"copied_at" db:"copied_at"`
	Type          int16       `json:"type" db:"type"`
	Checked       bool        `json:"checked" db:"checked"`
	Content       any         `json:"content" db:"content"`
	Cover         string      `json:"cover" db:"cover"`
	CategoryOwner int64       `json:"category_owner" db:"category_owner"`
	CategoryId    int         `json:"category_id" db:"category_id"`
}

// @Summary Get Note
// @Description Retrieve a specific note by its ID
// @Tags notes
// @Accept json
// @Produce json
// @Param Authorization header string true "Access JWT"
// @Param note_id path string true "Note ID"
// @Success 200 {object} GetNoteResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
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

type ListNotesResponse struct {
	Notes   []*types.Note `json:"notes"`
	HasMore bool          `json:"has_more"`
	Offset  int           `json:"offset"`
}

// ListNotes handles listing user notes
// @Summary List user notes
// @Description List notes for the authenticated user with optional filters
// @Tags notes
// @Produce json
// @Param Authorization header string true "Access JWT"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} ListNotesResponse "List of notes with new offset and has more flag"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/notes [get]
func ListNotes(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		creds := r.Context().Value("creds").(auth.Credentials)

		if creds.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		req := map[string]interface{}{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error("error while decoding request: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		offsetRaw := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetRaw)
		if err != nil {
			offset = 0
		}

		notes, hasMore, err := store.GetNotes(creds.ID, offset, req)
		if err != nil {
			slog.Error("error while getting notes: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ListNotesResponse{
			Notes:   notes,
			HasMore: hasMore,
			Offset:  offset + len(notes),
		}); err != nil {
			slog.Error("error while encoding response: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

type NewNoteRequest struct {
	Tags          []types.Tag `json:"tags" db:"tags"`
	Title         string      `json:"title" db:"title"`
	Source        string      `json:"source" db:"source"`
	Icon          any         `json:"icon" db:"icon"`
	Original      any         `json:"original" db:"original"`
	Font          string      `json:"font" db:"font"`
	CreatedAt     int64       `json:"created_at" db:"created_at"`
	CopiedAt      int64       `json:"copied_at" db:"copied_at"`
	Type          int16       `json:"type" db:"type"`
	Checked       bool        `json:"checked" db:"checked"`
	Content       any         `json:"content" db:"content"`
	Cover         string      `json:"cover" db:"cover"`
	CategoryOwner int64       `json:"category_owner" db:"category_owner"`
	CategoryId    int         `json:"category_id" db:"category_id"`
}

type NewNoteResponse struct {
	NoteID string `json:"note_id"`
}

// @Summary Create Note
// @Description Create a new note
// @Tags notes
// @Accept json
// @Produce json
// @Param Authorization header string true "Access JWT"
// @Param NewNoteRequest body NewNoteRequest true "New Note Request"
// @Success 201 {object} NewNoteResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes [post]
func CreateNote(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &NewNoteRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			slog.Error("error while decoding request body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		originalBytes, err := json.Marshal(req.Original)
		if err != nil {
			slog.Error("error while marshaling original: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		contentBytes, err := json.Marshal(req.Content)
		if err != nil {
			slog.Error("error while marshaling content: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		iconBytes, err := json.Marshal(req.Icon)
		if err != nil {
			slog.Error("error while marshaling icon: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		noteContentRaw := string(contentBytes)
		noteOriginalRaw := json.RawMessage(originalBytes)

		note := &types.Note{
			Tags:          req.Tags,
			Title:         req.Title,
			Source:        req.Source,
			Icon:          req.Icon,
			Original:      req.Original,
			Font:          req.Font,
			CreatedAt:     &req.CreatedAt,
			CopiedAt:      req.CopiedAt,
			Type:          req.Type,
			Checked:       req.Checked,
			Content:       req.Content,
			Cover:         req.Cover,
			CategoryOwner: &req.CategoryOwner,
			ContentRaw:    noteContentRaw,
			OriginalRaw:   &noteOriginalRaw,
			IconRaw:       json.RawMessage(iconBytes),
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

// @Summary Delete note
// @Description Delete a specific note by its ID
// @Tags notes
// @Accept json
// @Produce json
// @Param Authorization header string true "Access JWT"
// @Param note_id path string true "Note ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes/{note_id} [delete]
func DeleteNote(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds := r.Context().Value("creds").(auth.Credentials)
		if creds.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		note_id := chi.URLParam(r, "note_id")
		if err := store.DeleteNote(note_id, creds.ID); err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			slog.Error("error while deleting note: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

type CreateTagRequest struct {
	Tag types.Tag `json:"tag"`
}

// @Summary Create Tag
// @Description Create a new tag
// @Tags tags
// @Accept json
// @Produce json
// @Param Authorization header string true "Access JWT"
// @Param CreateTagRequest body CreateTagRequest true "Create Tag Request"
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

		if err := json.NewEncoder(w).Encode(tag); err != nil {
			slog.Error("error while encoding response: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

type UpdateNoteRequest map[string]interface{}

// @Summary Update Note
// @Description Update a specific note by its ID
// @Tags notes
// @Accept json
// @Produce json
// @Param Authorization header string true "Access JWT"
// @Param note_id path string true "Note ID"
// @Param UpdateNoteRequest body UpdateNoteRequest true "Update Note Request"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
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
