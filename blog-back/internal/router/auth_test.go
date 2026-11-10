package router

import (
	"blog-back/internal/config"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestManagementRoutesRequireSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("ADMIN_TOKEN", "router-test-password")
	r := Setup(config.Config{UploadPath: t.TempDir()}, nil)
	routes := []struct{ method, path string }{
		{"POST", "/article"}, {"POST", "/article/save"}, {"PUT", "/article/update"},
		{"DELETE", "/article/delete/1"}, {"POST", "/tag"}, {"PUT", "/tag"}, {"DELETE", "/tag/1"},
		{"POST", "/user/create"}, {"PUT", "/user/update"}, {"DELETE", "/user/delete/1"},
		{"POST", "/hollow"}, {"POST", "/upload/image"},
	}
	for _, route := range routes {
		req := httptest.NewRequest(route.method, route.path+"?admin_token=router-test-password", nil)
		req.Header.Set("X-Admin-Token", "router-test-password")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != 401 {
			t.Fatalf("unprotected %s %s: %d", route.method, route.path, w.Code)
		}
	}
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"password":"router-test-password"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var result struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil || w.Code != 200 {
		t.Fatalf("login: %s", w.Body)
	}
	// Invalid payloads should reach each protected handler, without a database.
	for _, path := range []string{"/article", "/hollow", "/upload/image"} {
		req = httptest.NewRequest("POST", path, strings.NewReader("invalid"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+result.Data.Token)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != 400 || !strings.Contains(w.Body.String(), "code") {
			t.Fatalf("session not accepted at %s: %s", path, w.Body)
		}
	}
	// Server restart (and hence password rotation with restart) invalidates old sessions.
	t.Setenv("ADMIN_TOKEN", "rotated-password")
	restarted := Setup(config.Config{UploadPath: t.TempDir()}, nil)
	req = httptest.NewRequest("GET", "/auth/session", nil)
	req.Header.Set("Authorization", "Bearer "+result.Data.Token)
	w = httptest.NewRecorder()
	restarted.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Fatal("old session survived restart")
	}
	req = httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"password":"router-test-password"}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	restarted.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Fatal("old password accepted after rotation")
	}
}
