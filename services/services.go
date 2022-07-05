package services

import (
	"github.com/go-zoox/connect/services/app"
	"github.com/go-zoox/connect/services/captcha"
	"github.com/go-zoox/connect/services/menu"
	"github.com/go-zoox/connect/services/token"
	"github.com/go-zoox/connect/services/user"
)

var Captcha = captcha.New()
var Token = token.New()

var App = app.New()
var User = user.New()
var Menu = menu.New()
