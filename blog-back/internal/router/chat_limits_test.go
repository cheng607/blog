package router

import (
	"blog-back/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestChatAliasesShareQuotaAndCannotSpoofIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var calls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer upstream.Close()
	cfg := config.Config{UploadPath: t.TempDir(), DashScopeAPIKey: "test", DashScopeBaseURL: upstream.URL, AIRatePerMinute: 1, AIDailyLimit: 2}
	r := Setup(cfg, nil)
	request := func(engine http.Handler, path, ip, forwarded string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", path, strings.NewReader(`{"message":"hi"}`))
		req.RemoteAddr = ip + ":1234"
		req.Header.Set("X-Forwarded-For", forwarded)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		return w
	}
	if w := request(r, "/chat", "203.0.113.1", "1.1.1.1"); w.Code != 200 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	if w := request(r, "/api/chat", "203.0.113.1", "2.2.2.2"); w.Code != 429 || w.Header().Get("Retry-After") == "" {
		t.Fatal("spoofed IP or alias bypassed limit", w.Code)
	}
	if w := request(r, "/api/chat", "203.0.113.2", ""); w.Code != 200 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	restarted := Setup(cfg, nil)
	if w := request(restarted, "/chat", "203.0.113.3", ""); w.Code != 429 {
		t.Fatal("quota bypassed on restart", w.Code)
	}
	if calls.Load() != 2 {
		t.Fatal("rejected calls reached provider", calls.Load())
	}
}
