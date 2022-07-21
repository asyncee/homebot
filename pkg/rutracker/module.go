package rutracker

import (
	"net/http"

	"go.uber.org/fx"

	cookiejar "github.com/juju/persistent-cookiejar"
)

var Module = fx.Options(
	fx.Provide(func() (http.CookieJar, error) {
		return cookiejar.New(nil)
	}),
	fx.Provide(NewClient),
)
