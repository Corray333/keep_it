package service

import (
	"fmt"
	"log/slog"

	"github.com/Corray333/keep_it/internal/domains/note/types"
)

type repository interface {
	GetNote(note_id string) (*types.Note, error)
	CreateNote(note *types.Note) (string, error)
	CheckNoteAccess(note_id string, user_id int) (bool, error)
	CreateTag(tag *types.Tag) (*types.Tag, error)
	UpdateNote(note_id string, data map[string]interface{}) error
	GetNotes(user_id int, offset int, filter map[string]interface{}) ([]*types.Note, bool, error)
	DeleteNote(note_id string, uid int) error
}

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetNote(uid int, noteID string) (*types.Note, error) {
	note, err := s.repo.GetNote(noteID)
	if err != nil {
		slog.Error("error while getting note: " + err.Error())
		return nil, fmt.Errorf("error while getting note: " + err.Error())
	}

	if note.Creator != uid {
		allowed, err := s.repo.CheckNoteAccess(note.ID, uid)
		if err != nil {
			return nil, fmt.Errorf("error while checking note access: " + err.Error())
		}
		if !allowed {
			return nil, fmt.Errorf("no access to this note")
		}
	}
	return note, nil
}

func (s *service) ListNotes(uid int, offset int, filters map[string]interface{}) ([]*types.Note, bool, error) {
	return s.repo.GetNotes(uid, offset, filters)
}

func (s *service) CreateNote(note *types.Note) (string, error) {
	return s.repo.CreateNote(note)
}

func (s *service) DeleteNote(note_id string, uid int) error {
	return s.repo.DeleteNote(note_id, uid)
}

func (s *service) CreateTag(tag *types.Tag) (*types.Tag, error) {
	return s.repo.CreateTag(tag)
}

func (s *service) UpdateNote(note_id string, data map[string]interface{}) error {
	return s.repo.UpdateNote(note_id, data)
}
