package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/Corray333/keep_it/internal/helpers"
)

var (
	ErrNoAccess = helpers.NewError(http.StatusForbidden, "no access")
)

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	CreateNote(ctx context.Context, note *entities.Note) (noteID string, err error)
	GetNote(ctx context.Context, noteID string) (*entities.Note, error)
}

type NoteService struct {
	repo repository
}

func New(repo repository) *NoteService {
	s := &NoteService{
		repo: repo,
	}
	return s
}

func (s *NoteService) Run() {
}

func (c *NoteService) GetNoteByID(ctx context.Context, userID int64, noteID string) (*entities.Note, error) {
	note, err := c.repo.GetNote(ctx, noteID)
	if err != nil {
		return nil, err
	}

	if note.CreatorID != userID {
		return nil, ErrNoAccess
	}

	return note, nil
}

func (c *NoteService) CreateNote(ctx context.Context, userID int64, note entities.Note) (noteID string, err error) {
	ctx, err = c.repo.Begin(ctx)
	if err != nil {
		return "", err
	}

	note.CreatorID = userID

	contentRaw, err := json.Marshal(note.Content)
	if err != nil {
		return "", err
	}

	note.ContentRaw = string(contentRaw)

	iconRaw, err := json.Marshal(note.Icon)
	if err != nil {
		return "", err
	}

	note.IconRaw = iconRaw

	noteID, err = c.repo.CreateNote(ctx, &note)
	if err != nil {
		_ = c.repo.Rollback(ctx)
		return "", err
	}

	if err := c.repo.Commit(ctx); err != nil {
		return "", err
	}

	return noteID, nil
}
