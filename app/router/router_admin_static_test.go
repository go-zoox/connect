package router_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/router"
	"github.com/go-zoox/zoox/defaults"
)

const adminUIMarker = `data-connect-builtin-admin="1"`

func newTestAppAdminStatic(t *testing.T, adminEnabled bool, entry string) http.Handler {
	t.Helper()
	cfg := &config.Config{
		Auth: config.Auth{Mode: "none"},
		Admin: config.Admin{
			Enabled: adminEnabled,
			Entry:   entry,
		},
	}
	cfg.ApplyDefault()
	app := defaults.Application()
	router.New(app, cfg)
	return app
}

func TestAdminStatic_Enabled_ServesAtDefaultEntry(t *testing.T) {
	h := newTestAppAdminStatic(t, true, "")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /admin/: want status %d, got %d", http.StatusOK, rec.Code)
	}
	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if !strings.Contains(string(body), adminUIMarker) {
		t.Fatalf("expected admin UI marker in body, got %q", string(body))
	}
}

func TestAdminStatic_Enabled_CustomEntry(t *testing.T) {
	h := newTestAppAdminStatic(t, true, "/console")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/console/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /console/: want status %d, got %d", http.StatusOK, rec.Code)
	}
	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if !strings.Contains(string(body), adminUIMarker) {
		t.Fatalf("expected admin UI marker for custom entry, got %q", string(body))
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/admin/", nil)
	h.ServeHTTP(rec2, req2)
	if rec2.Code == http.StatusOK {
		b2, _ := io.ReadAll(rec2.Body)
		if strings.Contains(string(b2), adminUIMarker) {
			t.Fatalf("did not expect admin UI at default /admin/ when entry is /console")
		}
	}
}

func TestAdminStatic_Enabled_RedirectsBareEntryToSlash(t *testing.T) {
	h := newTestAppAdminStatic(t, true, "/console")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/console", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound && rec.Code != http.StatusMovedPermanently {
		t.Fatalf("GET /console: want redirect, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if loc != "/console/" {
		t.Fatalf("Location: want %q, got %q", "/console/", loc)
	}
}

func TestAdminStatic_Disabled_NoBuiltinUIAtDefaultEntry(t *testing.T) {
	h := newTestAppAdminStatic(t, false, "")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/", nil)
	h.ServeHTTP(rec, req)
	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if strings.Contains(string(body), adminUIMarker) {
		snippet := string(body)
		if len(snippet) > 200 {
			snippet = snippet[:200]
		}
		t.Fatalf("did not expect admin UI marker when admin disabled, got status %d body snippet %q", rec.Code, snippet)
	}
}
