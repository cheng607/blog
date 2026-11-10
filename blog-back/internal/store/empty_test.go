package store

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"io"
	"testing"
)

// A database returning zero rows verifies the actual store scanners and their JSON contract.
type emptyDriver struct{}
type emptyConn struct{}
type emptyRows struct{}

func (emptyDriver) Open(string) (driver.Conn, error)  { return emptyConn{}, nil }
func (emptyConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (emptyConn) Close() error                        { return nil }
func (emptyConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }
func (emptyConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return emptyRows{}, nil
}
func (emptyRows) Columns() []string         { return []string{"unused"} }
func (emptyRows) Close() error              { return nil }
func (emptyRows) Next([]driver.Value) error { return io.EOF }

func init() { sql.Register("blog-empty-test", emptyDriver{}) }

func TestEmptyListsEncodeAsArrays(t *testing.T) {
	db, err := sql.Open("blog-empty-test", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	tags := NewTagStore(db)
	articles := NewArticleStore(db, tags)
	check := func(name string, value interface{}, err error) {
		t.Helper()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		data, err := json.Marshal(value)
		if err != nil || string(data) != "[]" {
			t.Fatalf("%s: expected [], got %s (%v)", name, data, err)
		}
	}
	u, err := NewUserStore(db).List()
	check("users", u, err)
	h, err := NewHollowStore(db).List()
	check("hollows", h, err)
	ts, err := tags.ListWithCount()
	check("tags", ts, err)
	a, err := articles.ListWithTags()
	check("articles", a, err)
	a, err = articles.ByTagID(1)
	check("tag articles", a, err)
	a, err = articles.ByUserID(1)
	check("user articles", a, err)
}
