package static

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/zoox"
)

//go:embed dist
var dist embed.FS

// DistFS returns the embedded admin UI asset tree (contents of dist/).
func DistFS() (http.FileSystem, error) {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		return nil, err
	}
	return http.FS(sub), nil
}

// Mount registers embedded admin static files on app at cfg.Admin.Entry when cfg.Admin.Enabled.
// It is a no-op when admin is disabled. The exact entry path redirects to entry + "/" for consistent file serving.
func Mount(app *zoox.Application, cfg *config.Config) error {
	if cfg == nil || !cfg.Admin.Enabled {
		return nil
	}

	prefix := config.NormalizedAdminEntry(cfg.Admin.Entry)
	fsys, err := DistFS()
	if err != nil {
		return fmt.Errorf("admin static fs: %w", err)
	}

	indexHTML, err := dist.ReadFile("dist/index.html")
	if err != nil {
		return fmt.Errorf("admin static index: %w", err)
	}

	// Do not use two app.Get routes for prefix vs prefix+"/": zoox normalizes them to the same key ("GET /admin").
	app.Use(func(ctx *zoox.Context) {
		if ctx.Method != http.MethodGet && ctx.Method != http.MethodHead {
			ctx.Next()
			return
		}
		p := ctx.Path
		switch {
		case p == prefix:
			ctx.Redirect(prefix+"/", 302)
			return
		case p == prefix+"/":
			ctx.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
			return
		case strings.HasPrefix(p, prefix+"/"):
			ctx.Next()
			return
		default:
			ctx.Next()
		}
	})
	app.StaticFS(prefix, fsys)
	return nil
}
