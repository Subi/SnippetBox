package mysql

import (
	"database/sql"
	"github.com/subi/snippetbox/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

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

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	newSnippet := &models.Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&newSnippet.ID, &newSnippet.Title, &newSnippet.Content, &newSnippet.Created, &newSnippet.Expires)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return newSnippet, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		newSnippet := &models.Snippet{}

		err := rows.Scan(&newSnippet.ID, &newSnippet.Title, &newSnippet.Content, &newSnippet.Created, &newSnippet.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, newSnippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
