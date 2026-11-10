package store

import (
	"database/sql"
	"strings"
	"time"

	"blog-back/internal/model"
)

type ArticleStore struct {
	db       *sql.DB
	tagStore *TagStore
}

func NewArticleStore(db *sql.DB, tagStore *TagStore) *ArticleStore {
	return &ArticleStore{db: db, tagStore: tagStore}
}

const articleJoinSQL = `
SELECT a.id, a.title, a.content, a.cover_image, a.publish_time, a.update_time, a.user_id, a.views_count,
       t.id AS tag_id, t.name AS tag_name, t.description AS tag_description
FROM article a
LEFT JOIN article_tag at ON a.id = at.article_id
LEFT JOIN tag t ON at.tag_id = t.id
`

func scanArticleRows(rows *sql.Rows) ([]model.Article, error) {
	order := []int{}
	byID := map[int]*model.Article{}

	for rows.Next() {
		var (
			id, userID, views int
			title, content    string
			cover             sql.NullString
			publish, update   sql.NullTime
			tagID             sql.NullInt64
			tagName, tagDesc  sql.NullString
		)
		if err := rows.Scan(&id, &title, &content, &cover, &publish, &update, &userID, &views,
			&tagID, &tagName, &tagDesc); err != nil {
			return nil, err
		}
		art, ok := byID[id]
		if !ok {
			a := model.Article{
				ID: id, Title: title, Content: content, UserID: userID,
				ViewsCount: &views, Tags: []model.Tag{},
			}
			if cover.Valid {
				s := cover.String
				a.CoverImage = &s
			}
			if publish.Valid {
				t := publish.Time
				a.PublishTime = &t
			}
			if update.Valid {
				t := update.Time
				a.UpdateTime = &t
			}
			byID[id] = &a
			order = append(order, id)
			art = &a
		}
		if tagID.Valid {
			art.Tags = append(art.Tags, model.Tag{
				ID: int(tagID.Int64), Name: tagName.String, Description: tagDesc.String,
			})
		}
	}
	list := make([]model.Article, 0, len(order))
	for _, id := range order {
		list = append(list, *byID[id])
	}
	return list, rows.Err()
}

func (s *ArticleStore) ListWithTags() ([]model.Article, error) {
	rows, err := s.db.Query(articleJoinSQL + " ORDER BY a.publish_time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanArticleRows(rows)
}

func (s *ArticleStore) ByTagID(tagID int) ([]model.Article, error) {
	rows, err := s.db.Query(articleJoinSQL+` WHERE at.tag_id = ? ORDER BY a.publish_time DESC`, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanArticleRows(rows)
}

func (s *ArticleStore) ByIDWithTags(id int) (*model.Article, error) {
	rows, err := s.db.Query(articleJoinSQL+` WHERE a.id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list, err := scanArticleRows(rows)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

func (s *ArticleStore) ByUserID(userID int) ([]model.Article, error) {
	rows, err := s.db.Query(`SELECT id, title, content, cover_image, publish_time, update_time, user_id, views_count
		FROM article WHERE user_id = ? ORDER BY publish_time DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]model.Article, 0)
	for rows.Next() {
		var a model.Article
		var cover sql.NullString
		var publish, update sql.NullTime
		var views int
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &cover, &publish, &update, &a.UserID, &views); err != nil {
			return nil, err
		}
		a.ViewsCount = &views
		if cover.Valid {
			s := cover.String
			a.CoverImage = &s
		}
		if publish.Valid {
			t := publish.Time
			a.PublishTime = &t
		}
		if update.Valid {
			t := update.Time
			a.UpdateTime = &t
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (s *ArticleStore) Save(a *model.Article, tagIDs []int) (int, error) {
	var authorCount int
	if err := s.db.QueryRow("SELECT COUNT(1) FROM `user` WHERE id = ?", a.UserID).Scan(&authorCount); err != nil {
		return 0, err
	}
	if authorCount == 0 {
		return 0, ErrInvalidAuthor
	}

	if err := s.validateTagIDs(tagIDs); err != nil {
		return 0, err
	}
	now := time.Now()
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO article (title, content, cover_image, publish_time, update_time, user_id, views_count)
		VALUES (?, ?, ?, ?, ?, ?, 0)`, a.Title, a.Content, nullStr(a.CoverImage), now, now, a.UserID)
	if err != nil {
		return 0, err
	}
	id64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	articleID := int(id64)

	if err := execInsertTags(tx, articleID, tagIDs); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return articleID, nil
}

func (s *ArticleStore) Update(a *model.Article, tagIDs []int) error {
	if err := s.validateTagIDs(tagIDs); err != nil {
		return err
	}
	now := time.Now()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var existingID int
	if err := tx.QueryRow("SELECT id FROM article WHERE id = ? FOR UPDATE", a.ID).Scan(&existingID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	_, err = tx.Exec(`UPDATE article SET title=?, content=?, cover_image=?, update_time=? WHERE id=?`,
		a.Title, a.Content, nullStr(a.CoverImage), now, a.ID)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM article_tag WHERE article_id = ?`, a.ID); err != nil {
		return err
	}
	if err := execInsertTags(tx, a.ID, tagIDs); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *ArticleStore) validateTagIDs(tagIDs []int) error {
	for _, tid := range tagIDs {
		ok, err := s.tagStore.Exists(tid)
		if err != nil {
			return err
		}
		if !ok {
			return ErrInvalidTag
		}
	}
	return nil
}

func execInsertTags(tx *sql.Tx, articleID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}
	var ph []string
	var vals []interface{}
	for _, tid := range tagIDs {
		ph = append(ph, "(?, ?)")
		vals = append(vals, articleID, tid)
	}
	_, err := tx.Exec(`INSERT INTO article_tag (article_id, tag_id) VALUES `+strings.Join(ph, ","), vals...)
	return err
}

func (s *ArticleStore) Delete(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM article_tag WHERE article_id = ?`, id); err != nil {
		return err
	}
	result, err := tx.Exec(`DELETE FROM article WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if err = checkDeleted(result); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *ArticleStore) IncrementViews(id int) (int, error) {
	// 使用原子操作避免并发竞争
	res, err := s.db.Exec(`UPDATE article SET views_count = views_count + 1 WHERE id = ?`, id)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rows == 0 {
		return 0, ErrNotFound
	}
	// 重新查询最新值返回
	var newCount int
	err = s.db.QueryRow(`SELECT views_count FROM article WHERE id = ?`, id).Scan(&newCount)
	if err != nil {
		return 0, err
	}
	return newCount, nil
}

func nullStr(p *string) interface{} {
	if p == nil || *p == "" {
		return nil
	}
	return *p
}
