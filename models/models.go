package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string
	AuthorId    int
	Description string
	PublishedAt string
}

type Student struct {
	gorm.Model
	FullName     string
	CampusRollNo string
	DateBorrowed string
	BookId       int64
	Book         Book
}
