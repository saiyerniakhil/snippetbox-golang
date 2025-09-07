package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "snippetbox.saiyerniakhil.in/internal/db"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	Queries *db.Queries
}

const (
	LIMIT = 10
)

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// used placeholders instead of string interpolation to avoid SQL injection.
	result, err := m.Queries.AddSnippet(context.TODO(), db.AddSnippetParams{
		Title:   title,
		Content: content,
		Expires: time.Now().AddDate(0, 0, expires),
	})
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	row, err := m.Queries.GetSnippetById(context.TODO(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return Snippet{
		ID:      int(row.ID),
		Title:   row.Title,
		Content: row.Content,
		Created: row.Created,
		Expires: row.Expires,
	}, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {

	rows, err := m.Queries.GetLatestSnippets(context.TODO(), LIMIT)
	if err != nil {
		return nil, err
	}

	var snippets []Snippet

	for _, row := range rows {

		snippets = append(snippets, Snippet{
			ID:      int(row.ID),
			Title:   row.Title,
			Content: row.Content,
			Created: row.Created,
			Expires: row.Expires,
		})
	}

	return snippets, nil
}
