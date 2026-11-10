package store

import (
	"database/sql"
	"time"

	"blog-back/internal/model"
)

type HollowStore struct {
	db *sql.DB
}

func NewHollowStore(db *sql.DB) *HollowStore {
	return &HollowStore{db: db}
}

func (s *HollowStore) List() ([]model.Hollow, error) {
	rows, err := s.db.Query(`SELECT id, word, date FROM hollow ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]model.Hollow, 0)
	for rows.Next() {
		var h model.Hollow
		var dt time.Time
		if err := rows.Scan(&h.ID, &h.Word, &dt); err != nil {
			return nil, err
		}
		h.Date = &dt
		list = append(list, h)
	}
	return list, rows.Err()
}

func (s *HollowStore) Save(h *model.Hollow) error {
	dt := time.Now()
	if h.Date != nil {
		dt = *h.Date
	}
	res, err := s.db.Exec(`INSERT INTO hollow (word, date) VALUES (?, ?)`, h.Word, dt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	h.ID = int(id)
	h.Date = &dt
	return nil
}
