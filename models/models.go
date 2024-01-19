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
	StudentId        int
	StudentsFullName string
}

type StudentInfo struct {
	gorm.Model
	FullName     string
	CampusRollNo string
	DateBorrowed string
	BookTitle    string
	Book         BookInfo `gorm:"foreignKey:BookTitle"`
}
