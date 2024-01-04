package controllers

import (
	"encoding/json"
	"fmt"
	"libraryManagementSystem/models"
	"log"
	"net/http"
	"libraryManagementSystem/database"
	"gorm.io/gorm"
)

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

}

func insertBook (book models.Book) int64 {
	db:= database.Database_connection();
	db.AutoMigrate()
}
