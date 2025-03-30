package entities

type CodeQuery struct {
	Username   string `json:"username"`
	TelegramID int64  `json:"tgID"`
	Type       int    `json:"type"`
	Code       string `json:"code"`
	Syn        int64  `json:"syn"`
}
