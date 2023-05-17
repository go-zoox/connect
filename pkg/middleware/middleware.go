package middleware

type User struct {
	ID       string `json:"user_id"`
	Nickname string `json:"user_nickname"`
	Avatar   string `json:"user_avatar"`
	Email    string `json:"user_email"`
}
