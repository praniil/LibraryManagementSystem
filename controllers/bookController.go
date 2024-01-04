package controllers

import (
	"encoding/json"
	"fmt"
	"libraryManagementSystem/database"
	"libraryManagementSystem/models"
	"log"
	"net/http"
)

type Response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	//.Decode(&book) ==> decoded data into book

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertID := insertBook(book)

	res := Response{
		ID:      insertID,
		Message: "User created successfully",
	}
	json.NewEncoder(w).Encode(res)

}

func insertBook(book models.Book) int64 {
	db := database.Database_connection()
	db.AutoMigrate(&models.Book{})
	result := db.Create(&book)
	if result.Error != nil {
		panic(fmt.Sprintf("Failed to execute the query: %v", result.Error))
	}
	fmt.Printf("Inserted a single record %v \n", book.ID)
	return int64(book.ID)

}
