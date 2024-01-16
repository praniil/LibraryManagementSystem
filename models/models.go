package models

import (
	pq "github.com/lib/pq"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string
	AuthorId    int
	Description string
	PublishedAt string
	TotalBooks  int
	StudentIds  pq.Int64Array `gorm:"type:bigInt[]"`
}

type Student struct {
	gorm.Model
	FullName     string
	CampusRollNo string
	DateBorrowed string
	BookId       int64
	Book         Book
}
