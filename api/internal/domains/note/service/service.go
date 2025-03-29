package service

import (
	"context"
	"errors"
	"net/http"

	user_entities "github.com/Corray333/keep_it/internal/domains/user/entities"
	"github.com/Corray333/keep_it/internal/helpers"
	"github.com/Corray333/keep_it/internal/storage"
)

var (
	ErrNoAccess = helpers.NewError(http.StatusForbidden, errors.New("no access"))
)

type repository interface {
	storage.Transactioner

	noteGetter
	noteCreater
	tagToNoteAdder
	tagCreater
	notesDeleter
	tagDeleter
	notesGetter
	tagFromNoteDeleter
	categorySetter
	tagsGetter
	categoryCreater
	categoryGetter
}

type userService interface {
	GetUser(ctx context.Context, searchUser *user_entities.User) (*user_entities.User, error)
}

type fileGetter interface {
	GetFileURL(ctx context.Context, name string) (string, error)
}

// Note service
// Requires UserService
type NoteService struct {
	userService

	fileGetter fileGetter

	transactioner      storage.Transactioner
	noteGetter         noteGetter
	noteCreater        noteCreater
	tagToNoteAdder     tagToNoteAdder
	tagCreater         tagCreater
	notesDeleter       notesDeleter
	tagDeleter         tagDeleter
	notesGetter        notesGetter
	tagFromNoteDeleter tagFromNoteDeleter
	categorySetter     categorySetter
	tagsGetter         tagsGetter
	categoryCreater    categoryCreater
	categoryGetter     categoryGetter
}

func New(options ...option) *NoteService {
	s := &NoteService{}

	for _, option := range options {
		option(s)
	}

	return s
}

type option func(*NoteService)

func WithFileGetter(fileGetter fileGetter) option {
	return func(s *NoteService) {
		s.fileGetter = fileGetter
	}
}

func WithUserService(userService userService) option {
	return func(s *NoteService) {
		s.userService = userService
	}
}

func WithTransactioner(transactioner storage.Transactioner) option {
	return func(s *NoteService) {
		s.transactioner = transactioner
	}
}

func WithNoteGetter(noteGetter noteGetter) option {
	return func(s *NoteService) {
		s.noteGetter = noteGetter
	}
}
func WithNoteCreater(noteCreater noteCreater) option {
	return func(s *NoteService) {
		s.noteCreater = noteCreater
	}
}

func WithTagToNoteAdder(tagToNoteAdder tagToNoteAdder) option {
	return func(s *NoteService) {
		s.tagToNoteAdder = tagToNoteAdder
	}
}

func WithTagCreater(tagCreater tagCreater) option {
	return func(s *NoteService) {
		s.tagCreater = tagCreater
	}
}

func WithNotesDeleter(notesDeleter notesDeleter) option {
	return func(s *NoteService) {
		s.notesDeleter = notesDeleter
	}
}

func WithTagDeleter(tagDeleter tagDeleter) option {
	return func(s *NoteService) {
		s.tagDeleter = tagDeleter
	}
}

func WithNotesGetter(notesGetter notesGetter) option {
	return func(s *NoteService) {
		s.notesGetter = notesGetter
	}
}

func WithTagFromNoteDeleter(tagFromNoteDeleter tagFromNoteDeleter) option {
	return func(s *NoteService) {
		s.tagFromNoteDeleter = tagFromNoteDeleter
	}
}

func WithCategorySetter(categorySetter categorySetter) option {
	return func(s *NoteService) {
		s.categorySetter = categorySetter
	}
}

func WithTagsGetter(tagsGetter tagsGetter) option {
	return func(s *NoteService) {
		s.tagsGetter = tagsGetter
	}
}

func WithCategoryCreater(categoryCreater categoryCreater) option {
	return func(s *NoteService) {
		s.categoryCreater = categoryCreater
	}
}

func WithCategoryGetter(categoryGetter categoryGetter) option {
	return func(s *NoteService) {
		s.categoryGetter = categoryGetter
	}
}

func WithRepository(repo repository) option {
	return func(s *NoteService) {
		s.transactioner = repo
		s.noteGetter = repo
		s.noteCreater = repo
		s.tagToNoteAdder = repo
		s.tagCreater = repo
		s.notesDeleter = repo
		s.tagDeleter = repo
		s.notesGetter = repo
		s.tagFromNoteDeleter = repo
		s.categorySetter = repo
		s.tagsGetter = repo
		s.categoryCreater = repo
		s.categoryGetter = repo
	}
}

func (s *NoteService) Run() {}

// func (c *NoteService) GetNewNotes(ctx context.Context, userID int64, offset int) ([]entities.Note, error) {
// 	return c.repo.GetNotes(ctx, userID, offset, nil)
// }
