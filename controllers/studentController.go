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

	var student models.StudentInfo
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Fatalf("failed to decode json format to the original format. %v", err)
	}
	bookTitle := student.BooksTitle
	studentId := insertStudent(student, bookTitle)
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

func insertStudent(student models.StudentInfo, bookTitle string) int64 {
	db := database.Database_connection()
	// if exists := db.Migrator().HasTable(&models.BookInfo{}); exists {
	// 	fmt.Println("Table books exists")
	// } else {
	// 	db.AutoMigrate(&models.BookInfo{}, &models.StudentInfo{})
	// }
	db.AutoMigrate(&models.BookInfo{}, &models.StudentInfo{})
	var book models.BookInfo
	tx := db.Begin()

	// var book models.BookInfo

	// Find the book by title and student_id
	tx.Model(&models.BookInfo{}).Where("title = ? AND students_id = ?", student.BooksTitle, 0).Find(&book)

	if book.ID == 0 {
		tx.Rollback()
		log.Fatalf("Book not found for title: %s and students_id: %d", student.BooksTitle, 0)
	}

	result := tx.Create(&student)

	if result.Error != nil {
		tx.Rollback()
		log.Fatalf("Unable to create a record for students. %v", result.Error)
	}
	// Update the book with student information
	book.StudentsId = int(student.ID)
	book.StudentsFullName = student.FullName

	// Update the BookInfo table with the modified book information
	tx.Exec("UPDATE book_infos SET students_id = ?, students_full_name = ? WHERE id = ?", book.StudentsId, book.StudentsFullName, book.ID)

	// Create the student record

	tx.Commit()

	return int64(student.ID)

	// return int64(student.ID)

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

func getStudent(id int64) (models.StudentInfo, error) {
	db := database.Database_connection()
	var student models.StudentInfo
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

func getAllStudents() ([]models.StudentInfo, error) {
	db := database.Database_connection()
	var students []models.StudentInfo //{arrays of information of different students}
	result := db.Find(&students)      //retrieves all the information from models.Student table

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
	var student models.StudentInfo
	json.NewDecoder(r.Body).Decode(&student)

	rowsUpdated := updateStudent(student, int64(id))
	msg := fmt.Sprintf("number of rows updated: %d", rowsUpdated)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

func updateStudent(student models.StudentInfo, id int64) int64 {
	db := database.Database_connection()
	result := db.Model(&models.StudentInfo{}).Where("id = ?", id).Updates(student)
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
	result := db.Delete(&models.BookInfo{}, id)
	if result.Error != nil {
		log.Fatalf("couldnt delete the row %v", result.Error)
	}
	rowsDeleted := result.RowsAffected
	return rowsDeleted
}
