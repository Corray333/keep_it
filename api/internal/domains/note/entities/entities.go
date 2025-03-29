package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NoteFilter struct {
	Category string
	Tags     []int64
	Text     string
	Source   string
}

type Source string

const (
	SourceTelegram = "tg"
	SourceWeb      = "web"
	SourceVK       = "vk"
)

type NoteType int8

const (
	NoteTypeDefault = iota
)

type Note struct {
	ID             uuid.UUID       `json:"id" db:"note_id"`
	CreatorID      int64           `json:"creator" db:"creator_id"`
	Tags           []Tag           `json:"tags" db:"tags"`
	Title          string          `json:"title" db:"title"`
	Source         Source          `json:"source" db:"source"`
	Original       string          `json:"original" db:"original"`
	ContentDecoded any             `json:"-" db:"-"`
	Content        json.RawMessage `json:"content" db:"content"`
	Cover          string          `json:"cover" db:"cover"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	CopiedAt  time.Time `json:"copiedAt" db:"copied_at"`

	Type NoteType `json:"type" db:"type"`

	Checked bool `json:"checked" db:"checked"`

	CategoryId uuid.UUID `json:"categoryID" db:"category_id"`

	IconDecoded Icon            `json:"-" db:"-"`
	Icon        json.RawMessage `json:"icon" db:"icon"`
}

type Tag struct {
	ID    int64  `json:"id" db:"tag_id"`
	Text  string `json:"text" db:"tag_text"`
	Color string `json:"color" db:"tag_color"`
	Owner int64  `json:"owner" db:"owner_id"`
}

type Icon struct {
	Type string `json:"type"`
	Data string `json:"data"`
	Name string `json:"name"`
}

type NewNoteMessage struct {
	Note   Note   `json:"note"`
	Source string `json:"source"`
	UserID string `json:"userID"`
}
