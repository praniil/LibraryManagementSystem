package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"libraryManagementSystem/database"
	"libraryManagementSystem/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var student models.Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Fatalf("failed to decode json format to the original format. %v", err)
	}
	bookId := student.BookId
	studentId := insertStudent(student, bookId)
	var msg string
	if studentId == 0 {
		msg = "Student not created"
	} else {
		msg = "A new Student crested"
	}
	res := Response{
		ID:      studentId,
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func insertStudent(student models.Student, bookID int64) int64 {
	db := database.Database_connection()
	if exists := db.Migrator().HasTable(&models.Book{}); exists {
		fmt.Println("Table books exists")
	} else {
		db.AutoMigrate(&models.Book{}, &models.Student{})
	}
	var book models.Book
	findBook := db.Find(&book, bookID)
	if findBook.Error != nil {
		fmt.Println("couldnt find a book with given bookid")
	}
	if book.TotalBooks > 0 {
		book.TotalBooks = book.TotalBooks - 1
		tx := db.Begin()
		tx.Model(&models.Book{}).Where("id = ?", bookID).Update("total_books", book.TotalBooks) //updates the book.TotalBooks-- in the table
		tx.Commit()
		result := db.Create(&student)
		if result.Error != nil {
			log.Fatalf("unable to create a table for students. %v", result.Error)
		}
		return int64(student.ID)
	} else {

		return 0
	}
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"]) //{id} extract
	if err != nil {
		log.Fatalf("couldnot extract id from the url, %v", err)
	}
	var student models.Student
	json.NewDecoder(r.Body).Decode(&student)

	rowsUpdated := updateStudent(student, int64(id))
	msg := fmt.Sprintf("number of rows updated: %d", rowsUpdated)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

func updateStudent(student models.Student, id int64) int64 {
	db := database.Database_connection()
	result := db.Model(&models.Student{}).Where("id = ?", id).Updates(student)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Printf("Record not found with id: %d", id)
	}
	if result.Error != nil {
		log.Fatalf("unable to update the record . %v", result.Error)
	}
	rowsUpdated := result.RowsAffected
	return rowsUpdated

}
