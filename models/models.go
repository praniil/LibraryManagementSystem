package models

type Book struct {
	ID          int64
	Title       string
	AuthorId    int
	Description string
	PublishedAt string
	Students    []Student
}

type Student struct {
	ID           int64
	FullName     string
	CampusRollNo int64
	DateBorrowed string
	BookId       int64
	Book         Book
}
