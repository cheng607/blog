package store

import (
	"blog-back/internal/model"
	"context"
	"database/sql"
	"html"
	"regexp"
	"strings"
	"time"
)

type ArticleSummary struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Excerpt     string      `json:"excerpt"`
	CoverImage  *string     `json:"coverImage,omitempty"`
	PublishTime *time.Time  `json:"publishTime,omitempty"`
	ViewsCount  *int        `json:"viewsCount,omitempty"`
	Tags        []model.Tag `json:"tags"`
}
type ArticlePage struct {
	Items    []ArticleSummary `json:"items"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

var previewRules = []*regexp.Regexp{
	regexp.MustCompile("(?s)```.*?(?:```|$)"),
	regexp.MustCompile("(?s)~~~.*?(?:~~~|$)"),
	regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`),
	regexp.MustCompile(`<[^>]*>`),
}
var previewLinks = regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
var previewMarkers = regexp.MustCompile("(?m)^\\s*[#>+-]+\\s*|[*_~`]")

func articleExcerpt(raw string) string {
	for _, rule := range previewRules {
		raw = rule.ReplaceAllString(raw, " ")
	}
	raw = previewLinks.ReplaceAllString(raw, "$1")
	raw = previewMarkers.ReplaceAllString(raw, "")
	plain := strings.Join(strings.Fields(html.UnescapeString(raw)), " ")
	runes := []rune(plain)
	if len(runes) > 120 {
		return string(runes[:120]) + "…"
	}
	return plain
}

// Page the articles before joining tags, so multi-tag articles are never split
// across pages. Only a bounded prefix of the body is read for the summary.
func (s *ArticleStore) Page(ctx context.Context, page, size, tagID, userID int) (ArticlePage, error) {
	result := ArticlePage{Items: []ArticleSummary{}, Page: page, PageSize: size}
	where := " WHERE 1=1"
	args := []interface{}{}
	if tagID > 0 {
		where += " AND EXISTS (SELECT 1 FROM article_tag f WHERE f.article_id = a.id AND f.tag_id = ?)"
		args = append(args, tagID)
	}
	if userID > 0 {
		where += " AND a.user_id = ?"
		args = append(args, userID)
	}
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM article a"+where, args...).Scan(&result.Total); err != nil {
		return result, err
	}
	if result.Total == 0 {
		result.Page = 1
		return result, nil
	}
	maxPage := (result.Total + size - 1) / size
	if result.Page > maxPage {
		result.Page = maxPage
	}
	args = append(args, size, (result.Page-1)*size)
	query := `SELECT a.id,a.title,a.content,a.cover_image,a.publish_time,a.update_time,a.user_id,a.views_count,
 t.id,t.name,t.description FROM (
 SELECT a.id,a.title,LEFT(a.content,1000) AS content,a.cover_image,a.publish_time,a.update_time,a.user_id,a.views_count
 FROM article a` + where + ` ORDER BY a.publish_time DESC,a.id DESC LIMIT ? OFFSET ?
 ) a LEFT JOIN article_tag at ON at.article_id=a.id LEFT JOIN tag t ON t.id=at.tag_id
 ORDER BY a.publish_time DESC,a.id DESC,t.id`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	articles, err := scanArticleRows(rows)
	if err != nil {
		return result, err
	}
	for _, a := range articles {
		result.Items = append(result.Items, ArticleSummary{ID: a.ID, Title: a.Title, Excerpt: articleExcerpt(a.Content), CoverImage: a.CoverImage, PublishTime: a.PublishTime, ViewsCount: a.ViewsCount, Tags: a.Tags})
	}
	return result, nil
}

type SiteSummary struct {
	Profile      *model.User `json:"profile"`
	ArticleCount int         `json:"articleCount"`
	TagCount     int         `json:"tagCount"`
	HollowCount  int         `json:"hollowCount"`
}
type SummaryStore struct{ db *sql.DB }

func NewSummaryStore(db *sql.DB) *SummaryStore { return &SummaryStore{db: db} }
func (s *SummaryStore) Get(ctx context.Context) (SiteSummary, error) {
	var result SiteSummary
	err := s.db.QueryRowContext(ctx, `SELECT (SELECT COUNT(*) FROM article),(SELECT COUNT(*) FROM tag),(SELECT COUNT(*) FROM hollow)`).Scan(&result.ArticleCount, &result.TagCount, &result.HollowCount)
	if err != nil {
		return result, err
	}
	var u model.User
	var avatar, signature sql.NullString
	err = s.db.QueryRowContext(ctx, "SELECT id,username,avatar,signature FROM `user` ORDER BY id LIMIT 1").Scan(&u.ID, &u.Username, &avatar, &signature)
	if err == sql.ErrNoRows {
		return result, nil
	}
	if err != nil {
		return result, err
	}
	u.Avatar = avatar.String
	u.Signature = signature.String
	result.Profile = &u
	return result, nil
}
