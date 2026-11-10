package store

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrInvalidTag    = errors.New("invalid tag")
	ErrInvalidAuthor = errors.New("invalid author")
	ErrConflict      = errors.New("record still referenced")
)

// MySQL reports zero changed rows for an unchanged UPDATE. Check existence
// before interpreting that as a missing record. Table names are internal constants.
func checkUpdated(db *sql.DB, table string, id int, result sql.Result) error {
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	var count int
	if err = db.QueryRow("SELECT COUNT(1) FROM "+table+" WHERE id = ?", id).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}
func checkDeleted(result sql.Result) error {
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
