package repository

import (
	"encoding/json"
	"time"

	"github.com/Corray333/keep_it/internal/domains/note/entities"
	"github.com/google/uuid"
)

type noteDb struct {
	NoteID     uuid.UUID       `db:"note_id"`
	CreatorID  int64           `db:"creator_id"`
	Title      string          `db:"title"`
	Source     string          `db:"source"`
	Original   string          `db:"original"`
	Icon       json.RawMessage `db:"icon"`
	CreatedAt  time.Time       `db:"created_at"`
	CopiedAt   time.Time       `db:"copied_at"`
	Type       int8            `db:"type"`
	Content    json.RawMessage `db:"content"`
	Cover      string          `db:"cover"`
	Checked    bool            `db:"checked"`
	CategoryID string          `db:"category_id"`
}

func (n *noteDb) ToNote() *entities.Note {
	return &entities.Note{
		ID:         n.NoteID,
		CreatorID:  n.CreatorID,
		Title:      n.Title,
		Source:     entities.Source(n.Source),
		Original:   n.Original,
		Icon:       n.Icon,
		CreatedAt:  n.CreatedAt,
		CopiedAt:   n.CopiedAt,
		Type:       entities.NoteType(n.Type),
		Content:    n.Content,
		Cover:      n.Cover,
		Checked:    n.Checked,
		CategoryId: n.CategoryID,
		Tags:       []entities.Tag{},
	}
}
