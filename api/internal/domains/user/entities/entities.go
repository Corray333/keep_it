package entities

type User struct {
	ID         int64  `json:"id,omitempty" db:"user_id"`
	Username   string `json:"username,omitempty" db:"username"`
	TelegramID int64  `json:"tgID,omitempty" db:"tg_id"`
	Email      string `json:"email,omitempty" db:"email"`
	Avatar     string `json:"avatar,omitempty" db:"avatar"`
	Password   string `json:"password,omitempty" db:"password"`
	RefCode    string `json:"-" db:"ref_code"`
}

type CodeQuery struct {
	Username   string `json:"username"`
	TelegramID int64  `json:"tgID"`
	Type       int    `json:"type"`
	Code       string `json:"code"`
}
