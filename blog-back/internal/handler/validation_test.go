package handler

import (
	"blog-back/internal/store"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCRUDValidationBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name, body string
		handler    gin.HandlerFunc
	}{
		{"article empty title", `{"title":" ","content":"x","userId":1,"tagIds":[1]}`, NewArticleHandler(nil).Save},
		{"article empty content", `{"title":"t","content":" ","userId":1,"tagIds":[1]}`, NewArticleHandler(nil).Save},
		{"article repeated tags", `{"title":"t","content":"x","userId":1,"tagIds":[1,1]}`, NewArticleHandler(nil).Save},
		{"article invalid author", `{"title":"t","content":"x","userId":-1,"tagIds":[1]}`, NewArticleHandler(nil).Save},
		{"article oversize content", `{"title":"t","content":"` + strings.Repeat("a", (1<<20)+1) + `","userId":1,"tagIds":[1]}`, NewArticleHandler(nil).Save},
		{"article update id", `{"id":-1}`, NewArticleHandler(nil).Update},
		{"article update blank", `{"id":1,"title":" ","content":"x","tagIds":[1]}`, NewArticleHandler(nil).Update},
		{"tag blank", `{"name":" "}`, NewTagHandler(nil).Save},
		{"tag update id", `{"id":0,"name":"valid"}`, NewTagHandler(nil).Update},
		{"tag long description", `{"id":1,"name":"valid","description":"` + strings.Repeat("字", 256) + `"}`, NewTagHandler(nil).Update},
		{"user invalid URL", `{"username":"valid","avatar":"javascript:alert(1)"}`, NewUserHandler(nil).Save},
		{"user negative id", `{"id":-1,"username":"valid"}`, NewUserHandler(nil).Update},
		{"hollow blank", `{"word":" "}`, NewHollowHandler(nil).Save},
		{"hollow long", `{"word":"` + strings.Repeat("字", 1001) + `"}`, NewHollowHandler(nil).Save},
		{"hollow invalid date", `{"word":"valid","date":"yesterday"}`, NewHollowHandler(nil).Save},
		{"extra JSON", `{"word":"valid"} {}`, NewHollowHandler(nil).Save},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.POST("/", tc.handler)
			w := chatRequest(r, "/", tc.body)
			if w.Code != 400 {
				t.Fatalf("%d %s", w.Code, w.Body)
			}
		})
	}
	for _, id := range []string{"-1", "0", "2147483648", "abc"} {
		r := gin.New()
		r.GET("/:id", NewArticleHandler(nil).GetByID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", "/"+id, nil))
		if w.Code != 400 {
			t.Fatal("invalid ID accepted", id)
		}
	}
}

type failingDriver struct{}
type failingConn struct{}

func (failingDriver) Open(string) (driver.Conn, error)  { return failingConn{}, nil }
func (failingConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (failingConn) Close() error                        { return nil }
func (failingConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }
func (failingConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return nil, errors.New("SQL SELECT password FROM secret_database; private-host")
}
func init() { sql.Register("handler-failing-test", failingDriver{}) }
func TestDatabaseErrorsAreNotReturnedToClient(t *testing.T) {
	db, err := sql.Open("handler-failing-test", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := gin.New()
	r.GET("/", NewUserHandler(store.NewUserStore(db)).List)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 500 || strings.Contains(w.Body.String(), "SQL") || strings.Contains(w.Body.String(), "private-host") {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}
func TestStoreErrorsHaveStablePublicStatus(t *testing.T) {
	for err, status := range map[error]int{store.ErrNotFound: 404, store.ErrInvalidTag: 400, store.ErrInvalidAuthor: 400, store.ErrConflict: 409} {
		r := gin.New()
		r.GET("/", func(c *gin.Context) { storeError(c, err) })
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
		if w.Code != status {
			t.Fatalf("%v: %d", err, w.Code)
		}
	}
}
