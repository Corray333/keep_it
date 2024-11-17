package transport

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/Corray333/keep_it/internal/helpers"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
)

type NoteTransport struct {
	router  *chi.Mux
	kafka   sarama.Consumer
	service service
}

type service interface {
	CreateNote(ctx context.Context, userID int64, note entities.Note) (noteID string, err error)
	GetNoteByID(ctx context.Context, userID int64, noteID string) (*entities.Note, error)

	GetNotes(ctx context.Context, userID int64, offset int) ([]entities.Note, error)
	GetNewNotes(ctx context.Context, userID int64, offset int) ([]entities.Note, error)

	GetTags(ctx context.Context, userID int64) ([]entities.Tag, error)
	CreateTag(ctx context.Context, tag *entities.Tag) error
	DeleteTag(ctx context.Context, tag *entities.Tag) error
	RemoveTagFromNote(ctx context.Context, tag *entities.Tag, noteID string) error
	AddTagToNote(ctx context.Context, tag *entities.Tag, noteID string, isNew bool) error
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
		r.Get("/api/notes/new", t.getNewNotes)
		r.Get("/api/notes/{note_id}", t.getNote)

		r.Post("/api/tags", t.createTag)
		r.Delete("/api/tags/{tagText}", t.deleteTag)
		r.Get("/api/tags", t.getTags)
		r.Post("/api/notes/{noteID}/tags", t.addTagToNote)
		r.Delete("/api/notes/{noteID}/tags/{tagText}", t.removeTagFromNote)
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
						var note entities.Note
						if err := json.Unmarshal(msg.Value, &note); err != nil {
							log.Printf("Failed to unmarshal message: %v", err)
							continue
						}

						if _, err := t.service.CreateNote(context.Background(), note.CreatorID, note); err != nil {
							slog.Error("Failed to create note: " + err.Error())
							continue
						}

					}
				case err := <-pc.Errors():
					slog.Error("Failed to consume message: " + err.Error())
				}
			}
		}(partition)
	}
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

	notes, err := t.service.GetNotes(ctx, userID, offset)
	if err != nil {
		slog.Error("failed to get notes: " + err.Error())
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		slog.Error("failed to encode response: " + err.Error())
		helpers.SendError(w, err)
		return
	}
}

func (t *NoteTransport) getNewNotes(w http.ResponseWriter, r *http.Request) {
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

	notes, err := t.service.GetNewNotes(ctx, userID, offset)
	if err != nil {
		slog.Error("failed to get notes: " + err.Error())
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		slog.Error("failed to encode response: " + err.Error())
		helpers.SendError(w, err)
		return
	}
}

type CreateTagRequest struct {
	entities.Tag
}

func (t *NoteTransport) createTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request: " + err.Error())
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
		slog.Error("failed to create tag: " + err.Error())
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
		slog.Error("failed to delete tag: " + err.Error())
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
		slog.Error("failed to get tags: " + err.Error())
		helpers.SendError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(tags); err != nil {
		slog.Error("failed to encode response: " + err.Error())
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
		slog.Error("failed to decode request: " + err.Error())
		helpers.SendError(w, err)
		return
	}

	req.Owner = userID

	noteID := chi.URLParam(r, "noteID")

	err := t.service.AddTagToNote(ctx, &req.Tag, noteID, req.IsNew)
	if err != nil {
		slog.Error("failed to add tag to note: " + err.Error())
		helpers.SendError(w, err)
		return
	}
}

func (t *NoteTransport) removeTagFromNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(helpers.CtxUserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusInternalServerError)
		slog.Error("user id not found in context")
		return
	}

	noteID := chi.URLParam(r, "noteID")
	tagText := chi.URLParam(r, "tagText")

	err := t.service.RemoveTagFromNote(ctx, &entities.Tag{
		Text:  tagText,
		Owner: userID,
	}, noteID)
	if err != nil {
		slog.Error("failed to remove tag from note: " + err.Error())
		helpers.SendError(w, err)
		return
	}
}
