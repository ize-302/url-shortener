// Package util
package util

import (
	"database/sql"
	"time"

	"ize-302/url-shortener/model"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func scanURLs(rows *sql.Rows) ([]model.URL, error) {
	var urls []model.URL
	for rows.Next() {
		var u model.URL
		err := rows.Scan(&u.ID, &u.URL, &u.Code, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

func (s *Store) FetchURLByID(id int) (model.URL, error) {
	rows, err := s.db.Query(`SELECT id, url, code, createdAt FROM urls WHERE id = ?`, id)
	if err != nil {
		return model.URL{}, err
	}
	defer rows.Close()

	urls, err := scanURLs(rows)
	if err != nil {
		return model.URL{}, err
	}

	if len(urls) == 0 {
		return model.URL{}, sql.ErrNoRows
	}
	return urls[0], nil
}

func (s *Store) FetchURLByCode(code string) (model.URL, error) {
	rows, err := s.db.Query(`SELECT id, url, code, createdAt FROM urls WHERE code = ?`, code)
	if err != nil {
		return model.URL{}, err
	}
	defer rows.Close()

	urls, err := scanURLs(rows)
	if err != nil {
		return model.URL{}, err
	}
	if len(urls) == 0 {
		return model.URL{}, sql.ErrNoRows
	}
	return urls[0], nil
}

func (s *Store) FetchURLs() ([]model.URL, error) {
	rows, err := s.db.Query("SELECT id, url, code, createdAt FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls, err := scanURLs(rows)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *Store) SaveURL(url *model.URL) (sql.Result, error) {
	code, err := GenerateRandomCode()
	if err != nil {
		return nil, err
	}
	createdAt := time.Now().UTC().Format(time.RFC3339)
	result, err := s.db.Exec(`INSERT into urls(url, code, createdAt) VALUES(?, ?, ?)`, url.URL, code, createdAt)
	if err != nil {
		return nil, err
	}
	return result, nil
}
