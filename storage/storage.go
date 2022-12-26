package storage

import (
	"fmt"
	"time"
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
	ID      string    `db:"id"`
	Title   string    `db:"title"`
	Created time.Time `db:"created"`
	Update  time.Time `db:"updated"`
}

func (b Book) String() string {
	return fmt.Sprintf(" %s - %d", b.Title, b.Created.Year())
}
