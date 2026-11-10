package handler

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestInvalidPaginationRejectedBeforeDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/", NewArticleHandler(nil).List)
	for _, query := range []string{"page=0", "page=-1", "page=x", "page=1000001", "pageSize=0", "pageSize=51", "pageSize=abc"} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", "/?"+query, nil))
		if w.Code != 400 {
			t.Fatalf("%s: %d", query, w.Code)
		}
	}
}
