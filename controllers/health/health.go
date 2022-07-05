package health

import "github.com/go-zoox/zoox"

func New() func(*zoox.Context) {
	return func(ctx *zoox.Context) {
		ctx.Status(200)
	}
}
