package types

import "encoding/json"

type Note struct {
	ID            string           `json:"id" db:"note_id"`
	Creator       int              `json:"creator" db:"creator"`
	Tags          []Tag            `json:"tags" db:"tags"`
	Title         string           `json:"title" db:"title"`
	Source        string           `json:"source" db:"source"`
	Original      any              `json:"original" db:"-"`
	OriginalRaw   *json.RawMessage `json:"-" db:"original"`
	Font          string           `json:"font" db:"font"`
	CreatedAt     *int64           `json:"created_at" db:"created_at"`
	CopiedAt      int64            `json:"copied_at" db:"copied_at"`
	Type          int16            `json:"type" db:"type"`
	Checked       bool             `json:"checked" db:"checked"`
	Content       any              `json:"content" db:"-"`
	ContentRaw    string           `json:"-" db:"content"`
	Cover         string           `json:"cover" db:"cover"`
	CategoryOwner *int64           `json:"category_owner" db:"category_owner"`
	CategoryId    *int             `json:"category_id" db:"category_id"`
}

type Tag struct {
	ID    int    `json:"id" db:"tag_id"`
	Text  string `json:"text" db:"text"`
	Color string `json:"color" db:"color"`
	Owner int    `json:"owner" db:"owner"`
}
