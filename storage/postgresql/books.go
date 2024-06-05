package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"boilgopher/storage/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	getBook = `
		SELECT 
			id,
			title,
			year,
			details,
			tags,
			created,
			updated
 		FROM books 
	`
	IDClause    = ` id = :id`
	YearClause  = ` year = :year`
	TitleClause = ` to_tsvector('simple', lower(title)) @@ to_tsquery('simple', lower(:title))`
)

func (s *Storage) CreateBook(ctx context.Context, svt *models.Book) (string, error) {
	var id string
	stmt, err := s.db.PrepareNamedContext(ctx, createBook)
	if err != nil {
		return "", fmt.Errorf("preparing named query for createBook: %w", err)
	}
	defer stmt.Close()

	if err := stmt.GetContext(ctx, &id, svt); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("failed adding book: %w", err)
		}
		return "", err
	}
	return id, nil
}

func (s *Storage) GetBook(ctx context.Context, opts ...models.BookFilterOption) ([]*models.Book, error) {
	f := &models.BookFilter{}
	for _, o := range opts {
		o(f)
	}

	query, args := buildBookFilter(f)
	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("preparing named query for GetBookByID: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.QueryContext(ctx, args)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the results
	var book []*models.Book
	for rows.Next() {
		var b models.Book
		var tags pq.StringArray
		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Year,
			&b.Details,
			&tags,
			&b.Created,
			&b.Update,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		if tags != nil {
			b.Tags = []string(tags)
		}
		book = append(book, &b)
	}
	// Check for errors in rows.Next()
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
		return nil, err
	}
	return book, nil
}

func buildBookFilter(f *models.BookFilter) (string, map[string]interface{}) {
	query := getBook
	params := make(map[string]interface{})

	// id filter
	if f.ID != uuid.Nil {
		query = addQueryString(query, IDClause)
		params["id"] = f.ID
	}

	// year filter
	if f.Year != "" {
		query = addQueryString(query, YearClause)
		params["year"] = f.Year
	}

	// title filter
	if f.Title != "" {
		query = addQueryString(query, TitleClause)
		params["title"] = f.Title
	}

	// Build the WHERE clause dynamically
	if len(f.Tags) > 0 {
		var subquery string
		var conditions []string
		for _, tag := range f.Tags {
			conditions = append(conditions, fmt.Sprintf("'%s' = ANY(tags)", tag))
		}
		subquery += strings.Join(conditions, " AND ")
		query = addQueryString(query, subquery)
	}

	return query, params
}

func addQueryString(query, clause string) string {
	if !strings.Contains(query, "where") {
		query += " where"
	} else {
		query += " and"
	}
	query += clause
	return query
}
