package service

import (
	"strings"

	gocaptcha "github.com/go-zoox/captcha"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/crypto/aes"
	"github.com/go-zoox/zoox"
)

var captchKey = "gz_cap"

var encryptor, _ = aes.NewCFB(256, &aes.HexEncoding{}, nil)

func GenerateCaptcha(cfg *config.Config, ctx *zoox.Context) {
	secret := []byte(strings.Repeat(cfg.SecretKey, 16)[:32])
	cap := gocaptcha.New()

	encrypted, err := encryptor.Encrypt([]byte(cap.Text()), secret)
	if err != nil {
		panic(err)
	}

	ctx.Session().Set(captchKey, string(encrypted))

	ctx.Status(200)

	cap.Write(ctx.Writer)
}

func ValidateCaptcha(cfg *config.Config, ctx *zoox.Context, input string) bool {
	secret := []byte(strings.Repeat(cfg.SecretKey, 16)[:32])
	cap, err := encryptor.Decrypt([]byte(ctx.Session().Get(captchKey)), secret)
	if err != nil {
		panic(err)
	}

	ctx.Session().Del(captchKey)

	if len(cap) != 0 {
		return strings.EqualFold(string(cap), input)
	}

	return false
}
