package models

import "time"

type Book struct {
	ID          int
	Title       string
	AuthorId    int
	Description string
	PublishedAt time.Time
}
