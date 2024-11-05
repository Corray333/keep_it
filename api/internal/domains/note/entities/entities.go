package entities

import "encoding/json"

type Note struct {
	ID         string `json:"id" db:"note_id"`
	CreatorID  int64  `json:"creator" db:"creator_id"`
	Tags       []Tag  `json:"tags" db:"tags"`
	Title      string `json:"title" db:"title"`
	Source     string `json:"source" db:"source"`
	Original   string `json:"original" db:"original"`
	Content    any    `json:"content" db:"-"`
	ContentRaw string `json:"-" db:"content"`
	Cover      string `json:"cover" db:"cover"`

	CreatedAt int64 `json:"created_at" db:"created_at"`
	CopiedAt  int64 `json:"copied_at" db:"copied_at"`

	Type int16 `json:"type" db:"type"`

	Checked bool `json:"checked" db:"checked"`

	CategoryId *string `json:"category_id" db:"category_id"`

	Icon    Icon            `json:"icon" db:"-"`
	IconRaw json.RawMessage `json:"-" db:"icon"`
}

type Tag struct {
	ID    int    `json:"id" db:"tag_id"`
	Text  string `json:"text" db:"text"`
	Color string `json:"color" db:"color"`
	Owner int    `json:"owner" db:"owner"`
}

type Icon struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
