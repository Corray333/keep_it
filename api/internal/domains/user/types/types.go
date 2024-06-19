package types

type User struct {
	ID            int    `json:"id,omitempty" db:"user_id"`
	Username      string `json:"username" db:"username"`
	Email         string `json:"email" db:"email"`
	Avatar        string `json:"avatar" db:"avatar"`
	Password      string `json:"password,omitempty" db:"password"`
	EmailVerified bool   `json:"email_verified" db:"email_verified"`
}
