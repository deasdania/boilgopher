package postgresql

import (
	"testing"

	"boilgopher/storage/models"
)

func TestStorage_CreateBook(t *testing.T) {
	envTest(t)
	s, tearDownFn := newTestStorage(t)
	t.Cleanup(tearDownFn)
	testData := &models.Book{
		Title:   "John Doe",
		Year:    "2000",
		Tags:    models.StringArray{"Horror"},
		Details: `{"writer":"Norman D"}`,
	}

	got, err := s.CreateBook(testData)
	if err != nil {
		t.Errorf("Storage.CreateBook() error = %v", err)
		return
	}
	if got == "" {
		t.Errorf("Storage.CreateBook() = got empty id %q", got)
	}
}
