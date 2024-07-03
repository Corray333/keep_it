package types

type User struct {
	ID               int    `json:"id,omitempty" db:"user_id"`
	Username         string `json:"username" db:"username"`
	TelegramUsername string `json:"tg_username" db:"tg_username"`
	Email            string `json:"email" db:"email"`
	Avatar           string `json:"avatar" db:"avatar"`
	Password         string `json:"password,omitempty" db:"password"`
	RefCode          string `json:"-" db:"ref_code"`
}

type CodeQuery struct {
	Username string `json:"username"`
	TG       string `json:"tg"`
	TypeID   int    `json:"type_id"`
	Code     string `json:"code"`
}
