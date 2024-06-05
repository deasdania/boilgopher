package models

import (
	"fmt"
	"strings"
	"time"

	"database/sql/driver"
	"github.com/google/uuid"
)

//go:generate stringer -type=BookType -linecomment
type BookType int

// comment and -linecomment for customizing
const (
	AdventureStories  BookType = iota + 1 // adventure-stories
	Classics                              // classics
	Crime                                 // crime
	FairyTales                            // fairy-tales
	Fantasy                               // fantasy
	HistoricalFiction                     // historical-fiction
	Horror                                // horror
	HumourAndSatire                       // humour-and-satire
)

type Book struct {
	ID      uuid.UUID   `db:"id"`
	Title   string      `db:"title"`
	Year    string      `db:"year"`
	Details string      `db:"details"`
	Tags    StringArray `db:"tags"`
	Created time.Time   `db:"created"`
	Update  time.Time   `db:"updated"`
}

func (b Book) String() string {
	return fmt.Sprintf(" %s - %d", b.Title, b.Created.Year())
}

type BookFilter struct {
	ID    uuid.UUID
	Title string
	Year  string
	Tags  StringArray
}

type BookFilterOption func(*BookFilter)

func BookFilterByID(ID uuid.UUID) BookFilterOption {
	return func(a *BookFilter) {
		a.ID = ID
	}
}

func BookFilterByYear(Year string) BookFilterOption {
	return func(a *BookFilter) {
		a.Year = Year
	}
}

func BookFilterByTitle(Title string) BookFilterOption {
	return func(a *BookFilter) {
		a.Title = Title
	}
}

func BookFilterByTags(Tags StringArray) BookFilterOption {
	return func(a *BookFilter) {
		a.Tags = Tags
	}
}

// StringArray is a custom type representing a slice of strings.
type StringArray []string

// Value implements the driver.Valuer interface for StringArray.
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return nil, nil
	}
	// Convert StringArray to PostgreSQL array string representation.
	arrayString := "{" + strings.Join(sa, ",") + "}"
	return []byte(arrayString), nil
}
