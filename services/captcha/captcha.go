package captcha

import (
	"strings"

	gocaptcha "github.com/go-zoox/captcha"
	"github.com/go-zoox/connect/config"
	"github.com/go-zoox/crypto/aes"
	"github.com/go-zoox/zoox"
)

var captchCookieKey = "gt_cap"

type CaptchaService struct {
	aes *aes.CFB
}

func New() *CaptchaService {
	a, err := aes.NewCFB(256, &aes.HexEncoding{}, nil)
	if err != nil {
		panic(err)
	}

	return &CaptchaService{
		aes: a,
	}
}

func (c *CaptchaService) Generate(cfg *config.Config, ctx *zoox.Context) {
	secret := []byte(strings.Repeat(cfg.SecretKey, 16)[:32])
	cap := gocaptcha.New()

	encrypted, err := c.aes.Encrypt([]byte(cap.Text()), secret)
	if err != nil {
		panic(err)
	}

	ctx.Session.Set(captchCookieKey, string(encrypted))

	ctx.Status(200)

	cap.Write(ctx.Writer)
}

func (c *CaptchaService) Validate(cfg *config.Config, ctx *zoox.Context, input string) bool {
	secret := []byte(strings.Repeat(cfg.SecretKey, 16)[:32])
	cap, err := c.aes.Decrypt([]byte(ctx.Session.Get(captchCookieKey)), secret)
	if err != nil {
		panic(err)
	}

	// ctx.Cookie.Del(captchCookieKey)
	ctx.Session.Del(captchCookieKey)

	if len(cap) != 0 {
		return strings.EqualFold(string(cap), input)
	}

	return false
}
