package user

import (
	"github.com/go-zoox/crypto/jwt"
)

// User ...
type User struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	//
	FeishuOpenID string `json:"feishu_open_id"`
	//
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
}

// Encode ...
func (u *User) Encode(signer jwt.Jwt) (string, error) {
	return signer.Sign(map[string]interface{}{
		"id":             u.ID,
		"nickname":       u.Nickname,
		"avatar":         u.Avatar,
		"email":          u.Email,
		"feishu_open_id": u.FeishuOpenID,
	})
}

// Decode ...
func (u *User) Decode(signer jwt.Jwt, token string) error {
	jwtValue, err := signer.Verify(token)
	if err != nil {
		return err
	}

	u.ID = jwtValue.Get("id").String()
	u.Nickname = jwtValue.Get("nickname").String()
	u.Avatar = jwtValue.Get("avatar").String()
	u.Email = jwtValue.Get("email").String()
	u.FeishuOpenID = jwtValue.Get("feishu_open_id").String()
	// u.Username = jwtValue.Get("username").String()
	// u.Permissions = jwtValue.Get("permissions").Array()

	return nil
}
