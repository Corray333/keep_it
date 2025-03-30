package entities

type User struct {
	ID       int64  `json:"id,omitempty" db:"user_id"`
	Username string `json:"username,omitempty" db:"username"`
	Email    string `json:"email,omitempty" db:"email"`
	Avatar   string `json:"avatar,omitempty" db:"avatar"`
	Password string `json:"password,omitempty" db:"password"`
	RefCode  string `json:"-" db:"ref_code"`

	TelegramID int64 `json:"tgID,omitempty" db:"tg_id"`
	VKID       int64 `json:"vkID,omitempty" db:"vk_id"`
}
