package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/Corray333/keep_it/internal/storage"
	"github.com/google/uuid"
)

type tagToNoteAdder interface {
	storage.Transactioner
	tagCreater
	AddTagToNote(ctx context.Context, tag *entities.Tag, noteID uuid.UUID) error
}

func (s *NoteService) AddTagToNote(ctx context.Context, tag *entities.Tag, noteID uuid.UUID, isNew bool) error {
	ctx, err := s.tagToNoteAdder.Begin(ctx)
	if err != nil {
		return err
	}
	defer s.tagToNoteAdder.Rollback(ctx)

	if isNew {
		if err := s.tagToNoteAdder.CreateTag(ctx, tag); err != nil {
			return err
		}
	}

	// TODO: add limit in 5 tags

	if err := s.tagToNoteAdder.AddTagToNote(ctx, tag, noteID); err != nil {
		return err
	}

	return s.tagToNoteAdder.Commit(ctx)
}
