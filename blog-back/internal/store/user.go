package store

import (
	"database/sql"

	"blog-back/internal/model"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) List() ([]model.User, error) {
	rows, err := s.db.Query("SELECT id, username, avatar, signature, created_at FROM `user` ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]model.User, 0)
	for rows.Next() {
		var u model.User
		var created sql.NullTime
		if err := rows.Scan(&u.ID, &u.Username, &u.Avatar, &u.Signature, &created); err != nil {
			return nil, err
		}
		if created.Valid {
			t := created.Time
			u.CreatedAt = &t
		}
		list = append(list, u)
	}
	return list, rows.Err()
}

func (s *UserStore) ByID(id int) (*model.User, error) {
	var u model.User
	var created sql.NullTime
	err := s.db.QueryRow("SELECT id, username, avatar, signature, created_at FROM `user` WHERE id = ?", id).
		Scan(&u.ID, &u.Username, &u.Avatar, &u.Signature, &created)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if created.Valid {
		t := created.Time
		u.CreatedAt = &t
	}
	return &u, nil
}

func (s *UserStore) Save(u *model.User) (int, error) {
	res, err := s.db.Exec("INSERT INTO `user` (username, avatar, signature, created_at) VALUES (?, ?, ?, NOW())",
		u.Username, u.Avatar, u.Signature)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *UserStore) Update(u *model.User) (int, error) {
	result, err := s.db.Exec("UPDATE `user` SET username=?, avatar=?, signature=? WHERE id=?",
		u.Username, u.Avatar, u.Signature, u.ID)
	if err != nil {
		return 0, err
	}
	return u.ID, checkUpdated(s.db, "`user`", u.ID, result)
}

func (s *UserStore) Delete(id int) error {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(1) FROM article WHERE user_id=?", id).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return ErrConflict
	}
	result, err := s.db.Exec("DELETE FROM `user` WHERE id=?", id)
	if err != nil {
		return err
	}
	return checkDeleted(result)
}
