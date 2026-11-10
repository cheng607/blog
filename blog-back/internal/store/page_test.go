package store

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

type pageDriver struct{}
type pageConn struct{ empty bool }
type pageRows struct {
	columns []string
	values  [][]driver.Value
	index   int
}

func (pageDriver) Open(name string) (driver.Conn, error) {
	return pageConn{empty: name == "empty"}, nil
}
func (pageConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (pageConn) Close() error                        { return nil }
func (pageConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }
func (r *pageRows) Columns() []string                { return r.columns }
func (r *pageRows) Close() error                     { return nil }
func (r *pageRows) Next(dest []driver.Value) error {
	if r.index >= len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.index])
	r.index++
	return nil
}
func (c pageConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if strings.Contains(query, "SELECT (SELECT COUNT(*)") {
		return &pageRows{columns: []string{"articles", "tags", "hollows"}, values: [][]driver.Value{{int64(7), int64(2), int64(3)}}}, nil
	}
	if strings.Contains(query, "FROM `user` ORDER BY id LIMIT 1") {
		return &pageRows{columns: []string{"id", "username", "avatar", "signature"}, values: nil}, nil
	}
	if strings.HasPrefix(query, "SELECT COUNT(*)") {
		total := int64(7)
		if c.empty {
			total = 0
		}
		return &pageRows{columns: []string{"total"}, values: [][]driver.Value{{total}}}, nil
	}
	if !strings.Contains(query, "LEFT(a.content,1000)") || !strings.Contains(query, "LIMIT ? OFFSET ?\n ) a LEFT JOIN") || !strings.Contains(query, "a.publish_time DESC,a.id DESC") {
		return nil, fmt.Errorf("pagination must precede the tag join with bounded content: %s", query)
	}
	if len(args) < 2 || args[len(args)-2].Value != int64(2) {
		return nil, fmt.Errorf("page size not bound: %v", args)
	}
	offset := args[len(args)-1].Value.(int64)
	if offset != 2 && offset != 6 {
		return nil, fmt.Errorf("incorrect offset: %d", offset)
	}
	if len(args) > 2 && (!strings.Contains(query, "EXISTS") || args[0].Value != int64(9)) {
		return nil, fmt.Errorf("incorrect tag filter")
	}
	now := time.Now()
	cols := []string{"id", "title", "content", "cover", "published", "updated", "user", "views", "tag", "tagName", "tagDescription"}
	rows := [][]driver.Value{
		{int64(5), "A", "# Heading\n\n**摘要** [link](https://example.com)", nil, now, now, int64(1), int64(0), int64(9), "tag A", ""},
		{int64(5), "A", "# Heading\n\n**摘要** [link](https://example.com)", nil, now, now, int64(1), int64(0), int64(10), "tag B", ""},
		{int64(4), "B", strings.Repeat("文", 1000), nil, now, now, int64(1), int64(1), int64(9), "tag A", ""},
	}
	if offset == 6 {
		rows = rows[:2]
	}
	return &pageRows{columns: cols, values: rows}, nil
}
func init() { sql.Register("blog-page-test", pageDriver{}) }
func TestArticlePageKeepsAllTagsAndOmitsBody(t *testing.T) {
	db, err := sql.Open("blog-page-test", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := NewArticleStore(db, NewTagStore(db))
	page, err := s.Page(context.Background(), 2, 2, 9, 0)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 7 || page.Page != 2 || len(page.Items) != 2 || len(page.Items[0].Tags) != 2 {
		t.Fatalf("unexpected page: %+v", page)
	}
	raw, _ := json.Marshal(page)
	if strings.Contains(string(raw), `"content"`) || len([]rune(page.Items[1].Excerpt)) > 121 {
		t.Fatal("body leaked or excerpt too long")
	}
	page, err = s.Page(context.Background(), 100, 2, 0, 0)
	if err != nil || page.Page != 4 || len(page.Items) != 1 {
		t.Fatalf("last page: %+v %v", page, err)
	}
	summary, err := NewSummaryStore(db).Get(context.Background())
	if err != nil || summary.ArticleCount != 7 || summary.TagCount != 2 || summary.HollowCount != 3 || summary.Profile != nil {
		t.Fatalf("summary: %+v %v", summary, err)
	}
}
func TestEmptyArticlePage(t *testing.T) {
	db, _ := sql.Open("blog-page-test", "empty")
	defer db.Close()
	page, err := NewArticleStore(db, NewTagStore(db)).Page(context.Background(), 3, 5, 0, 0)
	raw, _ := json.Marshal(page)
	if err != nil || page.Page != 1 || page.Total != 0 || !strings.Contains(string(raw), `"items":[]`) {
		t.Fatalf("%s %v", raw, err)
	}
}
func TestExcerptRemovesMarkdownAndHandlesChinese(t *testing.T) {
	got := articleExcerpt("# 标题\n```go\nsecret code\n```\n**你好** [世界](https://example.com) ![img](x.png)")
	if strings.Contains(got, "secret") || strings.Contains(got, "https") || !strings.Contains(got, "你好 世界") {
		t.Fatal(got)
	}
}
