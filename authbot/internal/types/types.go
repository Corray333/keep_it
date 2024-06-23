package types

type CodeQuery struct {
	Username string `json:"username"`
	TG       string `json:"tg"`
	TypeID   int    `json:"type_id"`
	Code     string `json:"code"`
}
