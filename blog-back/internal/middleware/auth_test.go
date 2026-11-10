package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAdminSessionLifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewAdminSessions("test-admin-password")
	now := time.Now()
	s.now = func() time.Time { return now }
	r := gin.New()
	r.POST("/login", s.Login)
	r.POST("/logout", s.Logout)
	r.POST("/write", s.Require(), func(c *gin.Context) { c.Status(204) })
	request := func(path, body, token string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	if w := request("/login", `{"password":"wrong"}`, ""); w.Code != 401 {
		t.Fatalf("wrong password: %d", w.Code)
	}
	if w := request("/write", "", "test-admin-password"); w.Code != 401 {
		t.Fatal("password accepted as session")
	}
	login := func() string {
		w := request("/login", `{"password":"test-admin-password"}`, "")
		if w.Code != 200 {
			t.Fatalf("login: %s", w.Body)
		}
		var result struct {
			Data struct {
				Token string `json:"token"`
			} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if len(result.Data.Token) != 64 {
			t.Fatal("invalid random session")
		}
		return result.Data.Token
	}
	token := login()
	if w := request("/write", "", token); w.Code != 204 {
		t.Fatal("valid session rejected")
	}
	if w := request("/logout", "", token); w.Code != 200 {
		t.Fatal("logout failed")
	}
	if w := request("/write", "", token); w.Code != 401 {
		t.Fatal("revoked session accepted")
	}
	token = login()
	now = now.Add(8 * time.Hour)
	if w := request("/write", "", token); w.Code != 401 {
		t.Fatal("expired session accepted")
	}
}

func TestMissingAdminConfigurationFailsClosed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewAdminSessions("")
	r := gin.New()
	r.POST("/login", s.Login)
	r.POST("/write", s.Require(), func(c *gin.Context) { t.Fatal("write reached") })
	for path, status := range map[string]int{"/login": 503, "/write": 401} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("POST", path, strings.NewReader(`{"password":""}`)))
		if w.Code != status {
			t.Fatalf("%s: got %d", path, w.Code)
		}
	}
}
