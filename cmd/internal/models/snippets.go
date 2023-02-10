package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO SNIPPETS (title, content, created, expires)
VALUES(?, ?, DATETIME('now'), DATETIME('now','+' || ? || ' days'))`

	result, err := m.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > datetime('now') AND id = ?
`
	s := &Snippet{}

	var created, expires string

	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &created, &expires)
	s.Created, _ = time.Parse(time.DateTime, created)
	s.Expires, _ = time.Parse(time.DateTime, expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// Return the 10 most recently created snippet
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > datetime('now') ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		var created, expires string
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &created, &expires)
		s.Created, _ = time.Parse(time.DateTime, created)
		s.Expires, _ = time.Parse(time.DateTime, expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
