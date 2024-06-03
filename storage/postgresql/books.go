package postgresql

import (
	"boilgopher/storage/models"
	"database/sql"
	"fmt"
)

const (
	createBook = `
		INSERT INTO books (
			title,
			year,
			tags,
			details
		) VALUES (
			:title,
			:year,
			:tags,
			:details	
		) RETURNING id
	`
)

func (s *Storage) CreateBook(svt *models.Book) (string, error) {
	var id string
	stmt, err := s.db.PrepareNamed(createBook)
	if err != nil {
		return "", fmt.Errorf("preparing named query for createBook: %w", err)
	}
	defer stmt.Close()

	if err := stmt.Get(&id, svt); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("failed adding book: %w", err)
		}
		return "", err
	}
	return id, nil
}
