package types

import (
	"time"
)

type CodeQuery struct {
	Username string `json:"username"`
	TG       string `json:"tg"`
	TypeID   int    `json:"type_id"`
	Code     string `json:"code"`
}

type Note struct {
	ID            string     `json:"id" db:"note_id"`
	Creator       int        `json:"creator" db:"creator"`
	Title         string     `json:"title" db:"title"` // Using null.String for nullable fields
	Icon          string     `json:"icon" db:"icon"`
	Source        string     `json:"source" db:"source"`
	Original      string     `json:"original" db:"original"`
	Font          string     `json:"font" db:"font"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	CopiedAt      time.Time  `json:"copied_at" db:"copied_at"`
	Type          int16      `json:"type" db:"type"`
	Checked       bool       `json:"checked" db:"checked"`
	Content       string     `json:"content" db:"content"`
	Cover         string     `json:"cover" db:"cover"`
	CategoryOwner *int64     `json:"category_owner" db:"category_owner"`
	CategoryId    *int       `json:"category_id" db:"category_id"`
}

type Original struct {
	Text string `json:"text"`
	Link string `json:"link"`
}
