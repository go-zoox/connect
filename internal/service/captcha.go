package service

import (
	"strings"

	gocaptcha "github.com/go-zoox/captcha"
	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/crypto/aes"
	"github.com/go-zoox/zoox"
)

var captchCookieKey = "gt_cap"

var encryptor, _ = aes.NewCFB(256, &aes.HexEncoding{}, nil)

func GenerateCaptcha(cfg *config.Config, ctx *zoox.Context) {
	secret := []byte(strings.Repeat(cfg.SecretKey, 16)[:32])
	cap := gocaptcha.New()

	encrypted, err := encryptor.Encrypt([]byte(cap.Text()), secret)
	if err != nil {
		panic(err)
	}

	ctx.Session().Set(captchCookieKey, string(encrypted))

	ctx.Status(200)

	cap.Write(ctx.Writer)
}

func ValidateCaptcha(cfg *config.Config, ctx *zoox.Context, input string) bool {
	secret := []byte(strings.Repeat(cfg.SecretKey, 16)[:32])
	cap, err := encryptor.Decrypt([]byte(ctx.Session().Get(captchCookieKey)), secret)
	if err != nil {
		panic(err)
	}

	ctx.Session().Del(captchCookieKey)

	if len(cap) != 0 {
		return strings.EqualFold(string(cap), input)
	}

	return false
}
