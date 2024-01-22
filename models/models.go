package models

import (
	"time"

	"gorm.io/gorm"
)

type BookInformation struct {
	gorm.Model
	Title            string
	AuthorId         int
	Description      string
	PublishedAt      string
	StudentsId       int
	StudentsFullName string
}

type StudentInformation struct {
	gorm.Model
	FullName     string
	CampusRollNo string
	DateBorrowed string
	BooksId      int
	Book         BookInformation `gorm:"foreignKey:BooksId"`
}

type LoanInformation struct {
	gorm.Model
	BookID        int64
	StudentsID    int64
	DueDate       time.Time
	RemainingTime time.Duration
	Fine          int64
	Returned      bool
	Book          BookInformation    `gorm:"foreignKey:BookID"`
	Student       StudentInformation `gorm:"foreignKey:StudentsID"`
}
