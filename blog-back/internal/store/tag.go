package store

import (
	"database/sql"

	"blog-back/internal/model"
)

type TagStore struct {
	db *sql.DB
}

func NewTagStore(db *sql.DB) *TagStore {
	return &TagStore{db: db}
}

func (s *TagStore) Exists(id int) (bool, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(1) FROM tag WHERE id = ?`, id).Scan(&n)
	return n > 0, err
}

const tagWithCountSQL = `
SELECT t.id, t.name, t.description, t.created_at, COUNT(at.article_id) AS cnt
FROM tag t
LEFT JOIN article_tag at ON t.id = at.tag_id
`

func (s *TagStore) ListWithCount() ([]model.Tag, error) {
	rows, err := s.db.Query(tagWithCountSQL + ` GROUP BY t.id ORDER BY t.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTags(rows)
}

func (s *TagStore) ByIDWithCount(id int) (*model.Tag, error) {
	rows, err := s.db.Query(tagWithCountSQL+` WHERE t.id = ? GROUP BY t.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list, err := scanTags(rows)
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return &list[0], nil
}

func scanTags(rows *sql.Rows) ([]model.Tag, error) {
	list := make([]model.Tag, 0)
	for rows.Next() {
		var t model.Tag
		var created sql.NullTime
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &created, &t.Count); err != nil {
			return nil, err
		}
		if created.Valid {
			ct := created.Time
			t.CreatedAt = &ct
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (s *TagStore) Save(t *model.Tag) (int, error) {
	res, err := s.db.Exec(`INSERT INTO tag (name, description, created_at) VALUES (?, ?, NOW())`, t.Name, t.Description)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *TagStore) Update(t *model.Tag) error {
	result, err := s.db.Exec(`UPDATE tag SET name=?, description=? WHERE id=?`, t.Name, t.Description, t.ID)
	if err != nil {
		return err
	}
	return checkUpdated(s.db, "tag", t.ID, result)
}

func (s *TagStore) Delete(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.Exec("DELETE FROM tag WHERE id=?", id)
	if err != nil {
		return err
	}
	if err = checkDeleted(result); err != nil {
		return err
	}
	if _, err = tx.Exec("DELETE FROM article_tag WHERE tag_id=?", id); err != nil {
		return err
	}
	return tx.Commit()
}
