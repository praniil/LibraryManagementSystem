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

	insertId := insertStudent(student)
	res := Response{
		ID:      insertId,
		Message: "A new Student crested",
	}
	json.NewEncoder(w).Encode(res)
}

func insertStudent(student models.Student) int64 {
	db := database.Database_connection()
	db.AutoMigrate(&models.Student{}, &models.Book{})
	result := db.Create(&student)
	if result.Error != nil {
		log.Fatalf("unable to create a table for students. %v", result.Error)
	}
	fmt.Println("Student Created with id: ", student.ID)
	return int64(student.ID)
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
