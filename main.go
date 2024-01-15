package main

import (
	"fmt"
	"libraryManagementSystem/database"
	"libraryManagementSystem/models"
	"libraryManagementSystem/router"
	"log"
	"net/http"
)

func main() {
	db := database.Database_connection()

	var students []models.Student
	var id int64
	id = 1
	db.Where("book_id = ?", id).Find(&students)

	for _, student := range students {
		fmt.Printf("Students ID: %d, Name: %s \n", student.ID, student.FullName)
	}
	fmt.Println("connected to the databse")
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		defer sqlDB.Close()

		fmt.Println("Database connection closed")
	}()
	r := router.Router()
	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
