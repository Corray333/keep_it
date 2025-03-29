package entities

import (
	"encoding/json"
	"time"
)

type Source string

type NoteType int8

const (
	NoteTypeDefault = iota
)

type Note struct {
	CreatorID      int64           `json:"creator" db:"creator_id"`
	Tags           []Tag           `json:"tags" db:"tags"`
	Title          string          `json:"title" db:"title"`
	Source         Source          `json:"source" db:"source"`
	Original       string          `json:"original" db:"original"`
	ContentDecoded []any           `json:"contentDecoded" db:"-"`
	Content        json.RawMessage `json:"content" db:"content"`
	Cover          string          `json:"cover" db:"cover"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	CopiedAt  time.Time `json:"copiedAt" db:"copied_at"`

	Type NoteType `json:"type" db:"type"`

	Checked bool `json:"checked" db:"checked"`

	IconDecoded Icon            `json:"-" db:"-"`
	Icon        json.RawMessage `json:"icon" db:"icon"`
}

type Tag struct {
	ID    int    `json:"id" db:"tag_id"`
	Text  string `json:"text" db:"tag_text"`
	Color string `json:"color" db:"tag_color"`
	Owner int64  `json:"owner" db:"owner_id"`
}

type Icon struct {
	Type string `json:"type"`
	Data string `json:"data"`
	Name string `json:"name"`
}

type Meta struct {
	Offset int `json:"offset"`
	Length int `json:"length"`

	Link          string `json:"link,omitempty"`
	Color         string `json:"color,omitempty"`
	Weight        string `json:"weight,omitempty"`
	Italic        bool   `json:"italic,omitempty"`
	Underline     bool   `json:"underline,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty"`
}

type RichText struct {
	PlainText string `json:"plain_text"`
	Meta      []Meta `json:"meta"`
}

type TextElement struct {
	Type     ElementType `json:"type"`
	RichText RichText    `json:"rich_text"`
}

type ImgElement struct {
	Type  ElementType `json:"type"`
	Src   string      `json:"src"`
	Width int         `json:"width"`
	Align string      `json:"align"`
}

type ElementType string

const (
	ElementTypeH1 ElementType = "h1"
	ElementTypeH2 ElementType = "h2"
	ElementTypeH3 ElementType = "h3"

	ElementTypeParagraph ElementType = "p"
	ElementTypeImage     ElementType = "img"
	ElementTypeCheckbox  ElementType = "checkbox"
)

type NewNoteMessage struct {
	Note   Note   `json:"note"`
	Source string `json:"source"`
	UserID string `json:"userID"`
}
