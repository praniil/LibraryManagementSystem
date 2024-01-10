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

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	//get the book id from request params
	params := mux.Vars(r) //gets paramaeter from the url
	// /api/getbook/{id}
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert the string to int %v", err)
	}

	book, err := getBook(int64(id))
	if err != nil {
		fmt.Println("unable to get the book")
	}

	json.NewEncoder(w).Encode(book)

}

func getBook(id int64) (models.Book, error) {
	db := database.Database_connection()

	var book models.Book
	result := db.First(&book, id) //db.first returns first record that matches the condition and returns to &book
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("no rows were returned")
		return book, nil
	} else if result.Error != nil {
		log.Fatalf("unable to query the book. %v", result.Error)
		return book, result.Error
	}
	fmt.Println("book in original format", book)
	return book, nil

	//original format from database {1 University Physics 2000 Conceptual Physics 1982-02-16}
	//later on encoded to json format while giving response to the user
}

func GetAllBook(w http.ResponseWriter, r *http.Request) {
	//r *http.Request http.Request struct contains information about an incomming HTTP request from a client, using a pointer to this struct allows the function to access and modify the req data
	//its like func insert(node *Node)
	w.Header().Set("Content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	books, err := getAllBooks()
	if err != nil {
		log.Fatalf("unable to find books. %v", err)
	}

	json.NewEncoder(w).Encode(books)
}

func getAllBooks() ([]models.Book, error) {
	db := database.Database_connection()
	var books []models.Book   //{arrays of information of different books}
	result := db.Find(&books) //retrieves all the information from models.Book table

	if result.Error != nil {
		log.Fatalf("unable to find books. %v", result.Error)
	}
	return books, nil
}


