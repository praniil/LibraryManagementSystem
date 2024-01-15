package controllers

import (
	"encoding/json"
	"fmt"
	"libraryManagementSystem/database"
	"libraryManagementSystem/models"
	"log"
	"net/http"
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

func UpdateStudent(w http.)
