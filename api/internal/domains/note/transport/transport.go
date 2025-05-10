package transport

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/Corray333/keep_it/internal/helpers"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NoteTransport struct {
	router  *chi.Mux
	kafka   sarama.Consumer
	service service
}

type service interface {
	CreateNote(ctx context.Context, note *entities.NewNoteMessage) (noteID uuid.UUID, err error)
	GetNoteByID(ctx context.Context, userID int64, noteID uuid.UUID) (*entities.Note, error)
	DeleteNotes(ctx context.Context, useID int64, noteIDs []uuid.UUID) error

	GetNotes(ctx context.Context, userID int64, offset int, filters entities.NoteFilter) ([]entities.Note, error)
	// GetNewNotes(ctx context.Context, userID int64, offset int) ([]entities.Note, error)
	SetCategory(ctx context.Context, userID int64, noteID uuid.UUID, categoryID string) error

	GetTags(ctx context.Context, userID int64) ([]entities.Tag, error)
	CreateTag(ctx context.Context, tag *entities.Tag) error
	DeleteTag(ctx context.Context, tag *entities.Tag) error
	DeleteTagFromNote(ctx context.Context, tag *entities.Tag, noteID uuid.UUID) error
	AddTagToNote(ctx context.Context, tag *entities.Tag, noteID uuid.UUID, isNew bool) error

	CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error)
	GetCategories(ctx context.Context, userID int64) ([]entities.Category, error)
}

func New(router *chi.Mux, service service) *NoteTransport {

	// Create new consumer
	brokers := []string{"kafka:9092"}

	// Set up Sarama configuration
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create a new consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	return &NoteTransport{
		router:  router,
		service: service,
		kafka:   consumer,
	}
}

func (t *NoteTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {
		r.Use(auth.NewAuthMiddleware())

		r.Post("/api/notes", t.createNote)
		r.Get("/api/notes", t.getNotes)
		// r.Get("/api/notes/new", t.getNewNotes)
		r.Get("/api/notes/{note_id}", t.getNote)
		r.Delete("/api/notes", t.deleteNote)

		r.Put("/api/notes/{note_id}/categories/{category_id}", t.setCategory)

		r.Post("/api/tags", t.createTag)
		r.Delete("/api/tags/{tagText}", t.deleteTag)
		r.Get("/api/tags", t.getTags)
		r.Post("/api/notes/{noteID}/tags", t.addTagToNote)
		r.Delete("/api/notes/{noteID}/tags/{tagText}", t.deleteTagFromNote)

		r.Post("/api/categories", t.createCategory)
		r.Get("/api/categories", t.getCategories)
	})
}

func (t *NoteTransport) Run() {

	// Topic to consume messages from
	topic := "newNotes"

	// Get partitions for the topic
	partitions, err := t.kafka.Partitions(topic)
	if err != nil {
		log.Fatalf("Failed to get partitions for topic %s: %v", topic, err)
	}

	// Consume messages from each partition continuously
	for _, partition := range partitions {
		go func(partition int32) {
			pc, err := t.kafka.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Fatalf("Failed to start consumer for partition %d: %v", partition, err)
			}
			defer pc.Close()

			for {
				select {
				case msg := <-pc.Messages():
					if msg != nil {
						note := &entities.NewNoteMessage{}
						if err := json.Unmarshal(msg.Value, note); err != nil {
							slog.Error("Failed to unmarshal message", "error", err)
							continue
						}

						if _, err := t.service.CreateNote(context.Background(), note); err != nil {
							slog.Error("Failed to create note", "error", err)
							continue
						}

					}
				case err := <-pc.Errors():
					slog.Error("Failed to consume message", "error", err)
				}
			}
		}(partition)
	}
}

type CreateNoteRequest struct {
	entities.Note
}

type CreateNoteResponse struct {
	NoteID uuid.UUID `json:"note_id"`
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
		slog.Error("failed to decode request: ", "error", err)
		helpers.SendError(w, err)
		return
	}

	req.Note.CreatorID = userID

	newNoteMsg := &entities.NewNoteMessage{
		Note:   req.Note,
		Source: "api",
		UserID: strconv.Itoa(int(userID)),
	}

	noteID, err := t.service.CreateNote(ctx, newNoteMsg)
	if err != nil {
		slog.Error("failed to create note: ", "error", err)
		helpers.SendError(w, err)
		return
	}

	resp := CreateNoteResponse{
		NoteID: noteID,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode response: ", "error", err)
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

	noteIDStr := chi.URLParam(r, "note_id")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		slog.Error("Failed to parse note id", "error", err)
		helpers.SendError(w, err)
		return
	}

	note, err := t.service.GetNoteByID(ctx, userID, noteID)
	if err != nil {
		slog.Error("failed to get note: ", "error", err)
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		slog.Error("failed to encode response: ", "error", err)
		helpers.SendError(w, err)
		return
	}
}

func (t *NoteTransport) getNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	offset, err := helpers.GetIntQueryParam(r, "offset")
	if err != nil {
		helpers.SendError(w, err)
		return
	}

	tags := r.URL.Query()["tags"]
	category := r.URL.Query().Get("category")

	tagsIDs := []int64{}
	for _, tag := range tags {
		tagID, err := strconv.ParseInt(tag, 10, 64)
		if err != nil {
			helpers.SendError(w, err)
			return
		}
		tagsIDs = append(tagsIDs, tagID)
	}

	filters := entities.NoteFilter{
		Category: category,
		Tags:     tagsIDs,
	}

	notes, err := t.service.GetNotes(ctx, userID, offset, filters)
	if err != nil {
		slog.Error("failed to get notes: ", "error", err)
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		slog.Error("failed to encode response: ", "error", err)
		helpers.SendError(w, err)
		return
	}
}

// func (t *NoteTransport) getNewNotes(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
// 	if !ok {
// 		http.Error(w, "user id not found in context", http.StatusInternalServerError)
// 		slog.Error("user id not found in context")
// 		return
// 	}

// 	offset, err := helpers.GetIntQueryParam(r, "offset")
// 	if err != nil {
// 		helpers.SendError(w, err)
// 		return
// 	}

// 	notes, err := t.service.GetNewNotes(ctx, userID, offset)
// 	if err != nil {
// 		slog.Error("failed to get notes: ", "error", err)
// 		helpers.SendError(w, err)
// 		return
// 	}

// 	if err := json.NewEncoder(w).Encode(notes); err != nil {
// 		slog.Error("failed to encode response: ", "error", err)
// 		helpers.SendError(w, err)
// 		return
// 	}
// }

type CreateTagRequest struct {
	entities.Tag
}

func (t *NoteTransport) createTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request: ", "error", err)
		helpers.SendError(w, err)
		return
	}

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	req.Owner = userID

	err := t.service.CreateTag(ctx, &req.Tag)
	if err != nil {
		slog.Error("failed to create tag: ", "error", err)
		helpers.SendError(w, err)
		return
	}
}

func (t *NoteTransport) deleteTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	tagText := chi.URLParam(r, "tagText")

	err := t.service.DeleteTag(ctx, &entities.Tag{
		Text:  tagText,
		Owner: userID,
	})
	if err != nil {
		slog.Error("failed to delete tag: ", "error", err)
		helpers.SendError(w, err)
		return
	}
}

func (t *NoteTransport) getTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	// offset, err := helpers.GetIntQueryParam(r, "offset")
	// if err != nil {
	// 	helpers.SendError(w, err)
	// 	return
	// }

	tags, err := t.service.GetTags(ctx, userID)
	if err != nil {
		slog.Error("failed to get tags: ", "error", err)
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(tags); err != nil {
		slog.Error("failed to encode response: ", "error", err)
		helpers.SendError(w, err)
		return
	}
}

type AddTagToNoteRequest struct {
	entities.Tag
	IsNew bool `json:"isNew"`
}

func (t *NoteTransport) addTagToNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	var req AddTagToNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request: ", "error", err)
		helpers.SendError(w, err)
		return
	}

	req.Owner = userID

	noteIDStr := chi.URLParam(r, "noteID")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		slog.Error("Failed to parse note id", "error", err)
		helpers.SendError(w, err)
		return
	}

	if err := t.service.AddTagToNote(ctx, &req.Tag, noteID, req.IsNew); err != nil {
		slog.Error("failed to add tag to note: ", "error", err)
		helpers.SendError(w, err)
		return
	}
}

func (t *NoteTransport) deleteTagFromNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	noteIDStr := chi.URLParam(r, "noteID")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		slog.Error("Failed to parse note id", "error", err)
		helpers.SendError(w, err)
		return
	}
	tagText := chi.URLParam(r, "tagText")

	if err := t.service.DeleteTagFromNote(ctx, &entities.Tag{
		Text:  tagText,
		Owner: userID,
	}, noteID); err != nil {
		slog.Error("failed to remove tag from note: ", "error", err)
		helpers.SendError(w, err)
		return
	}
}

func (t *NoteTransport) deleteNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	noteIDsStr := r.URL.Query()["note_id"]
	if len(noteIDsStr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	noteIDs := []uuid.UUID{}
	for _, noteIDStr := range noteIDsStr {
		noteID, err := uuid.Parse(noteIDStr)
		if err != nil {
			slog.Error("Failed to parse note id", "error", err)
			helpers.SendError(w, err)
			return
		}
		noteIDs = append(noteIDs, noteID)
	}

	if err := t.service.DeleteNotes(ctx, userID, noteIDs); err != nil {
		slog.Error("failed to remove note: ", "error", err)
		helpers.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *NoteTransport) setCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	noteIDStr := chi.URLParam(r, "note_id")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		slog.Error("Failed to parse note id", "error", err)
		helpers.SendError(w, err)
		return
	}
	categoryID := chi.URLParam(r, "category_id")

	if err := t.service.SetCategory(ctx, userID, noteID, categoryID); err != nil {
		slog.Error("failed to set note category: ", "error", err)
		helpers.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *NoteTransport) createCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		slog.Error("User id not found in context")
		helpers.SendError(w, helpers.ErrInternal)
		return
	}

	category := &entities.Category{}
	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		helpers.SendError(w, err)
		return
	}

	category.OwnerID = userID

	category, err := t.service.CreateCategory(ctx, category)
	if err != nil {
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(category); err != nil {
		helpers.SendError(w, err)
		return
	}

}

func (t *NoteTransport) getCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		slog.Error("User id not found in context")
		helpers.SendError(w, helpers.ErrInternal)
		return
	}

	categories, err := t.service.GetCategories(ctx, userID)
	if err != nil {
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		helpers.SendError(w, err)
		return
	}
}
