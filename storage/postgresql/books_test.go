package postgresql

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"boilgopher/storage/models"
)

func TestStorage_CreateBook(t *testing.T) {
	envTest(t)
	s, tearDownFn := newTestStorage(t)
	t.Cleanup(tearDownFn)
	ctx := context.Background()
	testData := &models.Book{
		Title:   "John Doe",
		Year:    "2000",
		Tags:    models.StringArray{"Horror"},
		Details: `{"writer":"Norman D"}`,
	}

	got, err := s.CreateBook(ctx, testData)
	if err != nil {
		t.Errorf("Storage.CreateBook() error = %v", err)
		return
	}
	uuidBook, err := uuid.Parse(got)
	if err != nil {
		t.Errorf("Storage.CreateBook() parse id = error %q", err)
	}
	filters := []models.BookFilterOption{}

	filters = append(filters, models.BookFilterByID(uuidBook))
	books, err := s.GetBook(ctx, filters...)
	if err != nil {
		t.Errorf("Storage.GetBook() error = %v", err)
	}
	if len(books) > 1 {
		t.Errorf("id should be unique")
	}
	b := books[0]
	assert.Equal(t, got, b.ID.String(), "Unexpected ID for person 1")

	// clean and append year filter
	filters = []models.BookFilterOption{}
	filters = append(filters, models.BookFilterByYear("2000"))

	books, err = s.GetBook(ctx, filters...)
	if err != nil {
		t.Errorf("Storage.GetBook() error = %v", err)
	}
	assert.Len(t, books, 1, "Expected a book")
	assert.Equal(t, b.Year, testData.Year, "Expected year is not match")

	// clean and append title filter
	filters = []models.BookFilterOption{}
	filters = append(filters, models.BookFilterByTitle("Doe"))
	books, err = s.GetBook(ctx, filters...)
	if err != nil {
		t.Errorf("Storage.GetBook() error = %v", err)
	}
	assert.Len(t, books, 1, "Expected a book")
	assert.Equal(t, b.Title, testData.Title, "Expected title is not match")

	// clean and append title filter
	filters = []models.BookFilterOption{}
	filters = append(filters, models.BookFilterByTags([]string{"Horror"}))
	books, err = s.GetBook(ctx, filters...)
	if err != nil {
		t.Errorf("Storage.GetBook() error = %v", err)
	}
	assert.Len(t, books, 1, "Expected a book")
	assert.Equal(t, b.Tags, testData.Tags, "Expected tags is not match")
}
