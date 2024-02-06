package utils

import (
	"regexp"
	"time"

	"github.com/go-zoox/cache"
)

// DefaultIgnoreAuthoriedPaths is the default ignore authorized paths
var DefaultIgnoreAuthoriedPaths = []string{
	// static
	"^/css/",
	"^/js/",
	"^/static/",
	"^/public/",
	"^/assets/",
	"\\.(css|js|ico|jpg|png|jpeg|webp|gif|socket|ws|map|webmanifest)$",

	// robots.txt
	"^/robots.txt$",

	// manifest
	"^/manifest.json$",

	// favicon
	"^/favicon.ico$",

	// captcha
	"^/captcha$",

	// umi
	// umi dev
	"^/__umi_ping$",
	"^/__umiDev/routes$",
	// umi prod
	"^/asset-manifest.json$",

	// webpack dev server
	"^/sockjs-node",
	"\\.hot-update.json$",

	// open page
	"^/open/(.*)",

	// open api
	"^/api/open/(.*)",

	// built-in apis
	"^/api/_/built_in_apis$",
	"^/api/login$",
	// fmt.Sprintf("^/api%s$", cfg.BuiltInAPIs.Login),
	"^/api/app$",
	// fmt.Sprintf("^/api%s$", cfg.BuiltInAPIs.App),
	"^/api/qrcode/",
	// fmt.Sprintf("^/api%s/", cfg.BuiltInAPIs.QRCode),
}

// CreateIsPathIgnoreAuthoriedMatcherOption ...
type CreateIsPathIgnoreAuthoriedMatcherOption struct {
	Excludes []string
}

// IsPathIgnoreAuthoriedMatcher ...
type IsPathIgnoreAuthoriedMatcher interface {
	Match(path string) bool
}

type isPathIgnoreAuthoriedMatcher struct {
	match func(path string) bool
}

func (m *isPathIgnoreAuthoriedMatcher) Match(path string) bool {
	return m.match(path)
}

// CreateIsPathIgnoreAuthoriedMatcher creates the IsPathIgnoreAuthoried matcher
func CreateIsPathIgnoreAuthoriedMatcher(opts ...func(opt *CreateIsPathIgnoreAuthoriedMatcherOption)) IsPathIgnoreAuthoriedMatcher {
	opt := &CreateIsPathIgnoreAuthoriedMatcherOption{}
	for _, o := range opts {
		o(opt)
	}

	excludes := append(DefaultIgnoreAuthoriedPaths, opt.Excludes...)

	excludesRe := []*regexp.Regexp{}
	for _, exclude := range excludes {
		excludesRe = append(excludesRe, regexp.MustCompile(exclude))
	}

	matchedCache := cache.New()
	match := func(path string) bool {
		if matchedCache.Has(path) {
			var matched bool
			if err := matchedCache.Get(path, &matched); err == nil {
				return matched
			}
		}

		matched := false
		for _, exclude := range excludesRe {
			matched = exclude.MatchString(path)
			if matched {
				break
			}
		}

		matchedCache.Set(path, &matched, 30*time.Second)

		return matched
	}

	return &isPathIgnoreAuthoriedMatcher{
		match: match,
	}
}
