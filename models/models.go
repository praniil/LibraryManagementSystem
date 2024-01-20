package models

import (
	"gorm.io/gorm"
)

type BookInfo struct {
	gorm.Model
	Title            string
	AuthorId         int
	Description      string
	PublishedAt      string
	StudentsId       int
	StudentsFullName string
}

type StudentInfo struct {
	gorm.Model
	FullName     string
	CampusRollNo string
	DateBorrowed string
	BooksTitle   string
	Book         BookInfo `gorm:"foreignKey:BooksTitle;references:Title"`
}
