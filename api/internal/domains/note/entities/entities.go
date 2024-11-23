package entities

import "encoding/json"

type Note struct {
	ID             string          `json:"id" db:"note_id"`
	CreatorID      int64           `json:"creator" db:"creator_id"`
	Tags           []Tag           `json:"tags" db:"tags"`
	Title          string          `json:"title" db:"title"`
	Source         string          `json:"source" db:"source"`
	Original       string          `json:"original" db:"original"`
	ContentDecoded any             `json:"-" db:"-"`
	Content        json.RawMessage `json:"content" db:"content"`
	Cover          string          `json:"cover" db:"cover"`

	CreatedAt int64 `json:"createdAt" db:"created_at"`
	CopiedAt  int64 `json:"copiedAt" db:"copied_at"`

	Type int16 `json:"type" db:"type"`

	Checked bool `json:"checked" db:"checked"`

	CategoryId *string `json:"categoryID" db:"category_id"`

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

type NewNoteMessage struct {
	Note   Note   `json:"note"`
	Source string `json:"source"`
	UserID string `json:"userID"`
}
