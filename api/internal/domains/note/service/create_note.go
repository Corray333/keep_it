package service

import (
	"context"
	"strconv"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	user_entities "github.com/Corray333/keep_it/internal/domains/user/entities"
	"github.com/Corray333/keep_it/internal/storage"
	"github.com/google/uuid"
)

type noteCreater interface {
	storage.Transactioner
	CreateNote(ctx context.Context, note *entities.Note) (noteID uuid.UUID, err error)
}

func (c *NoteService) CreateNote(ctx context.Context, note *entities.NewNoteMessage) (noteID uuid.UUID, err error) {

	if note.Source == "tg" {
		userID, err := strconv.Atoi(note.UserID)
		if err != nil {
			return uuid.Nil, err
		}

		user, err := c.GetUser(ctx, &user_entities.User{
			TelegramID: int64(userID),
		})
		if err != nil {
			return uuid.Nil, err
		}

		note.Note.CreatorID = user.ID
	}

	noteID, err = c.noteCreater.CreateNote(ctx, &note.Note)
	if err != nil {
		return uuid.Nil, err
	}

	return noteID, nil
}
