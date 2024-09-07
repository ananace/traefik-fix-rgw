package fix_rgw

import (
  "context"
	"net/http"
	"strings"
)

// Config the plugin configuration.
type Config struct {
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
	}
}

// WebsitePathConverter a WebsitePathConverter plugin.
type FixRGW struct {
	next      http.Handler
	name      string
}

// New created a new WebsitePathConverter plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &FixRGW{
		next: next,
		name: name,
	}, nil
}

func (a *FixRGW) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, "%7E") || strings.Contains(req.URL.Path, "%7e") {
		req.URL.Path = strings.ReplaceAll(req.URL.Path, "%7E", "~")
		req.URL.Path = strings.ReplaceAll(req.URL.Path, "%7e", "~")
	}
	a.next.ServeHTTP(rw, req)
}
