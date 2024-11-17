package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/Corray333/keep_it/internal/helpers"
)

var (
	ErrNoAccess = helpers.NewError(http.StatusForbidden, errors.New("no access"))
)

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	CreateNote(ctx context.Context, note *entities.Note) (noteID string, err error)
	GetNote(ctx context.Context, noteID string) (*entities.Note, error)
	GetNotes(ctx context.Context, userID int64, offset int, filters []helpers.Filter) ([]entities.Note, error)

	GetTags(ctx context.Context, userID int64) ([]entities.Tag, error)
	CreateTag(ctx context.Context, tag *entities.Tag) error
	DeleteTag(ctx context.Context, tag *entities.Tag) error
	RemoveTagFromNote(ctx context.Context, tag *entities.Tag, noteID string) error
	AddTagToNote(ctx context.Context, tag *entities.Tag, noteID string) error
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

	if note.Tags == nil {
		note.Tags = []entities.Tag{}
	}

	return note, nil
}

func (c *NoteService) CreateNote(ctx context.Context, userID int64, note entities.Note) (noteID string, err error) {
	ctx, err = c.repo.Begin(ctx)
	if err != nil {
		return "", err
	}

	note.CreatorID = userID

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

func (c *NoteService) GetNotes(ctx context.Context, userID int64, offset int, filters map[string][]string) ([]entities.Note, error) {

	newFilters := []helpers.Filter{}

	for key, values := range filters {
		switch key {
		case "tag":
			newFilters = append(newFilters, helpers.Filter{
				Field:     helpers.FilterKeyTag,
				Operation: "IN",
				Value:     values,
			})
		}
	}

	notes, err := c.repo.GetNotes(ctx, userID, offset, newFilters)
	if err != nil {
		return nil, err
	}

	for i := range notes {
		if notes[i].Tags == nil {
			notes[i].Tags = []entities.Tag{}
		}
	}

	return notes, nil
}

func (c *NoteService) GetNewNotes(ctx context.Context, userID int64, offset int) ([]entities.Note, error) {
	return c.repo.GetNotes(ctx, userID, offset, nil)
}

func (c *NoteService) CreateTag(ctx context.Context, tag *entities.Tag) error {
	// TODO: add limit in 128 tags
	return c.repo.CreateTag(ctx, tag)
}

func (c *NoteService) DeleteTag(ctx context.Context, tag *entities.Tag) error {
	return c.repo.DeleteTag(ctx, tag)
}

func (c *NoteService) RemoveTagFromNote(ctx context.Context, tag *entities.Tag, noteID string) error {
	return c.repo.RemoveTagFromNote(ctx, tag, noteID)
}

func (s *NoteService) AddTagToNote(ctx context.Context, tag *entities.Tag, noteID string, isNew bool) error {
	ctx, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer s.repo.Rollback(ctx)

	if isNew {
		if err := s.repo.CreateTag(ctx, tag); err != nil {
			return err
		}
	}

	// TODO: add limit in 5 tags

	if err := s.repo.AddTagToNote(ctx, tag, noteID); err != nil {
		return err
	}

	return s.repo.Commit(ctx)
}

func (c *NoteService) GetTags(ctx context.Context, userID int64) ([]entities.Tag, error) {
	return c.repo.GetTags(ctx, userID)
}
