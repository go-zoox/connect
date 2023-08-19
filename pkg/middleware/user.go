package middleware

import (
	"github.com/go-zoox/crypto/jwt"
)

// User ...
type User struct {
	ID       string `json:"user_id"`
	Nickname string `json:"user_nickname"`
	Avatar   string `json:"user_avatar"`
	Email    string `json:"user_email"`
	//
	FeishuOpenID string `json:"user_feishu_open_id"`
}

// Encode ...
func (u *User) Encode(signer jwt.Jwt) (string, error) {
	return signer.Sign(map[string]interface{}{
		"user_id":             u.ID,
		"user_nickname":       u.Nickname,
		"user_avatar":         u.Avatar,
		"user_email":          u.Email,
		"user_feishu_open_id": u.FeishuOpenID,
	})
}

// Decode ...
func (u *User) Decode(signer jwt.Jwt, token string) error {
	jwtValue, err := signer.Verify(token)
	if err != nil {
		return err
	}

	u.ID = jwtValue.Get("user_id").String()
	u.Nickname = jwtValue.Get("user_nickname").String()
	u.Avatar = jwtValue.Get("user_avatar").String()
	u.Email = jwtValue.Get("user_email").String()
	u.FeishuOpenID = jwtValue.Get("user_feishu_open_id").String()

	return nil
}
