package service

import (
	"context"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
)

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
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

func (c *NoteService) CreateNote(ctx context.Context, note entities.Note) {

}
