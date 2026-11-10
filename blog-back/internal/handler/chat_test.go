package handler

import (
	"blog-back/internal/ailimit"
	"blog-back/internal/config"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func chatFixture(t *testing.T, upstream http.HandlerFunc) (*ChatHandler, *gin.Engine) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(strings.NewReader(string(raw)))
		upstream(w, r)
	}))
	t.Cleanup(server.Close)
	h := NewChatHandler(config.Config{DashScopeAPIKey: "private-test-key", DashScopeModel: "primary,backup", DashScopeBaseURL: server.URL, UploadPath: t.TempDir(), AIRatePerMinute: 60})
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.POST("/chat", h.Chat)
	r.POST("/api/chat", h.Chat)
	return h, r
}
func chatRequest(r http.Handler, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func TestChatSuccessAndFallback(t *testing.T) {
	var calls atomic.Int32
	_, r := chatFixture(t, func(w http.ResponseWriter, r *http.Request) {
		var body dashScopeRequest
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			t.Error("invalid request")
		}
		if body.MaxTokens != 1024 || r.Header.Get("Authorization") != "Bearer private-test-key" {
			t.Error("missing upstream constraints")
		}
		calls.Add(1)
		if body.Model == "primary" {
			w.WriteHeader(503)
			return
		}
		w.Write([]byte(`{"choices":[{"message":{"content":"answer"}}]}`))
	})
	w := chatRequest(r, "/chat", `{"message":"question"}`)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "answer") || calls.Load() != 2 {
		t.Fatalf("%d %s calls=%d", w.Code, w.Body, calls.Load())
	}
}
func TestChatRejectsInputBeforeCallingProvider(t *testing.T) {
	var calls atomic.Int32
	_, r := chatFixture(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) })
	for _, body := range []string{`{}`, `{"message":"  "}`, `{"message":"` + strings.Repeat("问", 2001) + `"}`, `{} {}`, `not json`} {
		w := chatRequest(r, "/chat", body)
		if w.Code != 400 {
			t.Fatalf("expected 400: %d %s", w.Code, w.Body)
		}
	}
	w := chatRequest(r, "/chat", `{"message":"`+strings.Repeat("x", 20000)+`"}`)
	if w.Code != 413 || calls.Load() != 0 {
		t.Fatalf("oversize: %d calls=%d", w.Code, calls.Load())
	}
}
func TestChatProviderErrorsAreRedactedAndNotBlindlyRetried(t *testing.T) {
	var calls atomic.Int32
	_, r := chatFixture(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(401)
		w.Write([]byte(`{"error":{"message":"private-test-key secret.internal"}}`))
	})
	w := chatRequest(r, "/chat", `{"message":"hi"}`)
	if w.Code != 502 || calls.Load() != 1 || strings.Contains(w.Body.String(), "private") || strings.Contains(w.Body.String(), "secret") {
		t.Fatalf("%d %s calls=%d", w.Code, w.Body, calls.Load())
	}
}
func TestChatAttemptTimeoutAllowsFallback(t *testing.T) {
	var calls atomic.Int32
	h, r := chatFixture(t, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			<-r.Context().Done()
			return
		}
		w.Write([]byte(`{"choices":[{"message":{"content":"fallback"}}]}`))
	})
	h.attemptTimeout = 30 * time.Millisecond
	h.totalTimeout = time.Second
	w := chatRequest(r, "/chat", `{"message":"hi"}`)
	if w.Code != 200 || calls.Load() != 2 {
		t.Fatalf("%d %s calls=%d", w.Code, w.Body, calls.Load())
	}
}
func TestChatTotalTimeoutStopsFallback(t *testing.T) {
	var calls atomic.Int32
	h, r := chatFixture(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1); <-r.Context().Done() })
	h.attemptTimeout = time.Second
	h.totalTimeout = 30 * time.Millisecond
	w := chatRequest(r, "/chat", `{"message":"hi"}`)
	if w.Code != 504 || calls.Load() != 1 {
		t.Fatalf("%d %s calls=%d", w.Code, w.Body, calls.Load())
	}
}
func TestChatClientCancellationStopsWork(t *testing.T) {
	started := make(chan struct{})
	var calls atomic.Int32
	_, r := chatFixture(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1); close(started); <-r.Context().Done() })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest("POST", "/chat", strings.NewReader(`{"message":"hi"}`)).WithContext(ctx)
	done := make(chan struct{})
	go func() { r.ServeHTTP(httptest.NewRecorder(), req); close(done) }()
	<-started
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("request not canceled")
	}
	if calls.Load() != 1 {
		t.Fatal("canceled request retried")
	}
}

func TestChatConcurrencyRejectsBeforeUpstreamAndReleasesSlot(t *testing.T) {
	started := make(chan struct{})
	unblock := make(chan struct{})
	var calls atomic.Int32
	h, r := chatFixture(t, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			close(started)
			<-unblock
		}
		w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	})
	h.limiter = ailimit.New(60, 1, 10, filepath.Join(t.TempDir(), "quota.json"))
	done := make(chan *httptest.ResponseRecorder, 1)
	go func() { done <- chatRequest(r, "/chat", `{"message":"hi"}`) }()
	<-started
	w := chatRequest(r, "/api/chat", `{"message":"hi"}`)
	if w.Code != 429 || calls.Load() != 1 {
		t.Error("concurrent request reached provider")
	}
	close(unblock)
	if w = <-done; w.Code != 200 {
		t.Fatal(w.Code)
	}
	if w = chatRequest(r, "/chat", `{"message":"hi"}`); w.Code != 200 {
		t.Fatal("slot not released", w.Code)
	}
}

func TestChatQuotaFailureDoesNotCallProvider(t *testing.T) {
	var calls atomic.Int32
	h, r := chatFixture(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) })
	// A directory cannot be read as the quota JSON file.
	h.limiter = ailimit.New(60, 1, 10, t.TempDir())
	w := chatRequest(r, "/chat", `{"message":"hi"}`)
	if w.Code != 503 || calls.Load() != 0 {
		t.Fatalf("quota error allowed provider call: %d %d", w.Code, calls.Load())
	}
}
