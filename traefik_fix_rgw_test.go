package traefik_fix_rgw

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		name       string
		config     *Config
		url        string
		next       func(t *testing.T) http.Handler
		assertFunc func(t *testing.T, rw *httptest.ResponseRecorder)
	}{
		{
			name:   "base case",
			config: &Config{},
			url:    "http://localhost/",
			assertFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				t.Helper()
				if rr.Result().StatusCode != http.StatusOK {
					t.Fatalf("expected OK, got %d", rr.Result().StatusCode)
				}
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					if req.URL.String() != "http://localhost/" {
						t.Fatalf("wanted path %s, got req.URL: %+v", "http://localhost/", req.URL.String())
					}
				})
			},
		},
		{
			name:   "tilde path fix",
			config: &Config{},
			url:    "http://localhost/%7e/%7Epath",
			assertFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				t.Helper()
				if rr.Result().StatusCode != http.StatusOK {
					t.Fatalf("expected %d, got %d", http.StatusOK, rr.Result().StatusCode)
				}
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					if req.URL.String() != "http://localhost/~/~path" {
						t.Fatalf("wanted path %s, got req.URL: %+v", "http://localhost/~/~path", req.URL.String())
					}
				})
			},
		},
		{
			name:   "tilde query fix",
			config: &Config{},
			url:    "http://localhost/%7epath?request=1%7E27",
			assertFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				t.Helper()
				if rr.Result().StatusCode != http.StatusOK {
					t.Fatalf("expected %d, got %d", http.StatusOK, rr.Result().StatusCode)
				}
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					if req.URL.String() != "http://localhost/~path?request=1~27" {
						t.Fatalf("wanted path %s, got req.URL: %+v", "http://localhost/~path?request=1~27", req.URL.String())
					}
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			handler, err := New(ctx, tt.next(t), tt.config, "traefik_fix_rgw")
			if err != nil {
				t.Fatalf("error with new rgw fix: %+v", err)
			}
			recorder := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, tt.url, nil)
			if err != nil {
				t.Fatalf("error with new request: %+v", err)
			}

			handler.ServeHTTP(recorder, req)
			tt.assertFunc(t, recorder)
		})
	}
}
