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
	// db.AutoMigrate(&models.Book{}, &models.Student{})
	var book models.Book
	findBook := db.Find(&book, bookID)
	if findBook.Error != nil {
		fmt.Println("couldnt find a book with given bookid")
	}
	if book.TotalBooks > 0 {
		book.TotalBooks = book.TotalBooks - 1
		tx := db.Begin()
		tx.Model(&models.Book{}).Where("id = ?", bookID).Update("total_books", book.TotalBooks) //updates the book.TotalBooks-- in the table
		result := db.Create(&student)
		if result.Error != nil {
			log.Fatalf("unable to create a table for students. %v", result.Error)
		}
		book.StudentIds = append(book.StudentIds, int64(student.ID))
		tx.Model(&models.Book{}).Where("id = ?", bookID).Update("student_ids", book.StudentIds)
		tx.Commit()
		return int64(student.ID)
	} else {

		return 0
	}
}

func GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to extract id from the url. %v", err)
	}

	student, err := getStudent(int64(id))
	if err != nil {
		fmt.Println("unable to get the student")
	}
	json.NewEncoder(w).Encode(student)
}

func getStudent(id int64) (models.Student, error) {
	db := database.Database_connection()
	var student models.Student
	result := db.First(&student, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Printf("couldnt find the record with %d, %v", id, result.Error)
		return student, result.Error
	}
	if result.Error != nil {
		fmt.Printf("couldnt get the record with id : %d, error: %v", id, result.Error)
		return student, result.Error
	}
	return student, nil
}
func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	//r *http.Request http.Request struct contains information about an incomming HTTP request from a client, using a pointer to this struct allows the function to access and modify the req data
	//its like func insert(node *Node)
	w.Header().Set("Content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	students, err := getAllStudents()
	if err != nil {
		log.Fatalf("unable to find students. %v", err)
	}

	json.NewEncoder(w).Encode(students)
}

func getAllStudents() ([]models.Student, error) {
	db := database.Database_connection()
	var students []models.Student //{arrays of information of different students}
	result := db.Find(&students)  //retrieves all the information from models.Student table

	if result.Error != nil {
		log.Fatalf("unable to find students. %v", result.Error)
	}
	return students, nil
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

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("couldnt extract id from the url. %v", err)
	}

	deletedId := deleteStudent(int64(id))
	res := Response{
		ID:      deletedId,
		Message: "this row is deleted",
	}
	json.NewEncoder(w).Encode(res)
}

func deleteStudent(id int64) int64 {
	db := database.Database_connection()
	result := db.Delete(&models.Book{}, id)
	if result.Error != nil {
		log.Fatalf("couldnt delete the row %v", result.Error)
	}
	rowsDeleted := result.RowsAffected
	return rowsDeleted
}
