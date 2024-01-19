package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title            string
	AuthorId         int
	Description      string
	PublishedAt      string
	TotalBooks       int
	StudentId        int
	StudentsFullName string
}

type Student struct {
	gorm.Model
	FullName     string
	CampusRollNo string
	DateBorrowed string
	BookTitle    string
	Book         Book
}
