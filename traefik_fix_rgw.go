package traefik_fix_rgw

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
	return &Config{}
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
	req.URL.Path = fixTildes(req.URL.Path)
	req.URL.RawPath = fixTildes(req.URL.RawPath)
	req.URL.RawQuery = fixTildes(req.URL.RawQuery)
	req.RequestURI = req.URL.RequestURI()

	a.next.ServeHTTP(rw, req)
}

func fixTildes(str string) (string) {
	if strings.Contains(str, "%7E") || strings.Contains(str, "%7e") {
		str = strings.ReplaceAll(str, "%7E", "~")
		str = strings.ReplaceAll(str, "%7e", "~")
	}
	return str
}
