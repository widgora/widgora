package server

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

var testUI = fstest.MapFS{
	"index.html":                  {Data: []byte("<html>index</html>")},
	"robots.txt":                  {Data: []byte("User-agent: *\n")},
	"_app/immutable/entry/app.js": {Data: []byte("console.log(1)")},
}

func newTestHandler(ui fstest.MapFS) http.Handler {
	return New(slog.New(slog.DiscardHandler), ui)
}

func do(t *testing.T, h http.Handler, method, target string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec
}

func TestRoutes(t *testing.T) {
	h := newTestHandler(testUI)
	tests := []struct {
		name, method, target string
		status               int
		body                 string
		cache                string
	}{
		{"healthz", http.MethodGet, "/healthz", http.StatusOK, "ok\n", "no-store"},
		{"root serves index", http.MethodGet, "/", http.StatusOK, "<html>index</html>", "no-cache"},
		{"static file", http.MethodGet, "/robots.txt", http.StatusOK, "User-agent: *\n", "no-cache"},
		{"immutable asset", http.MethodGet, "/_app/immutable/entry/app.js", http.StatusOK, "console.log(1)", "public, max-age=31536000, immutable"},
		{"client route falls back", http.MethodGet, "/hosts/example", http.StatusOK, "<html>index</html>", "no-cache"},
		{"dot segments are redirected", http.MethodGet, "/../../etc/passwd", http.StatusTemporaryRedirect, "", ""},
		{"path outside ui falls back", http.MethodGet, "/etc/passwd", http.StatusOK, "<html>index</html>", "no-cache"},
		{"missing asset is 404", http.MethodGet, "/_app/immutable/missing.js", http.StatusNotFound, "", ""},
		{"api is not the ui", http.MethodGet, "/api/v1/unknown", http.StatusNotFound, "", ""},
		{"post to ui rejected", http.MethodPost, "/", http.StatusMethodNotAllowed, "", ""},
		{"post to healthz rejected", http.MethodPost, "/healthz", http.StatusMethodNotAllowed, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(t, h, tt.method, tt.target)
			if rec.Code != tt.status {
				t.Fatalf("status = %d, want %d", rec.Code, tt.status)
			}
			if tt.body != "" && rec.Body.String() != tt.body {
				t.Errorf("body = %q, want %q", rec.Body.String(), tt.body)
			}
			if tt.cache != "" && rec.Header().Get("Cache-Control") != tt.cache {
				t.Errorf("Cache-Control = %q, want %q", rec.Header().Get("Cache-Control"), tt.cache)
			}
			if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
				t.Error("missing X-Content-Type-Options header")
			}
		})
	}
}

func TestUINotBuilt(t *testing.T) {
	h := newTestHandler(fstest.MapFS{".gitkeep": {}})
	rec := do(t, h, http.MethodGet, "/")
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
	if rec := do(t, h, http.MethodGet, "/healthz"); rec.Code != http.StatusOK {
		t.Fatalf("healthz status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestServeShutsDownOnCancel(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() { done <- Serve(ctx, slog.New(slog.DiscardHandler), ln, newTestHandler(testUI)) }()

	resp, err := http.Get("http://" + ln.Addr().String() + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if !strings.HasPrefix(string(body), "ok") {
		t.Fatalf("healthz body = %q", body)
	}

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Serve returned %v, want nil", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Serve did not return after cancel")
	}
}
