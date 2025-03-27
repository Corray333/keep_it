package repository

import (
	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/google/uuid"
)

type tagDB struct {
	NoteID uuid.UUID `db:"note_id"`
	Text   string    `db:"tag_text"`
	Color  string    `db:"tag_color"`
}

func (t *tagDB) ToTag() *entities.Tag {
	return &entities.Tag{
		Text:  t.Text,
		Color: t.Color,
	}
}
