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

	var book models.BookInfo
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

func insertBook(book models.BookInfo) int64 {
	db := database.Database_connection()
	// if exists := db.Migrator().HasTable(&models.BookInfo{}); exists {
	// 	fmt.Println("Table books exists")
	// } else {
	// 	db.AutoMigrate(&models.BookInfo{}, &models.StudentInfo{})
	// }
	db.AutoMigrate(&models.BookInfo{}, &models.StudentInfo{})
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

func getBook(id int64) (models.BookInfo, error) {
	db := database.Database_connection()

	var book models.BookInfo
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

func getAllBooks() ([]models.BookInfo, error) {
	db := database.Database_connection()
	var books []models.BookInfo //{arrays of information of different books}
	result := db.Find(&books)   //retrieves all the information from models.Book table

	if result.Error != nil {
		log.Fatalf("unable to find books. %v", result.Error)
	}
	return books, nil
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"]) //{id} extract
	if err != nil {
		log.Fatalf("couldnot extract id from the url, %v", err)
	}
	var book models.BookInfo
	json.NewDecoder(r.Body).Decode(&book)

	rowsUpdated := updatebook(book, int64(id))
	msg := fmt.Sprintf("number of rows updated: %d", rowsUpdated)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

func updatebook(book models.BookInfo, id int64) int64 {
	db := database.Database_connection()
	result := db.Model(&models.BookInfo{}).Where("id = ?", id).Updates(book)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Printf("Record not found with id: %d", id)
	}
	if result.Error != nil {
		log.Fatalf("unable to update the record . %v", result.Error)
	}
	rowsUpdated := result.RowsAffected
	return rowsUpdated

}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println("unable to convert string to int type")
	}

	deletedRow := deleteRow(int64(id))
	msg := fmt.Sprintf("Total no of deleted rows: %v", deletedRow)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func deleteRow(id int64) int64 {
	db := database.Database_connection()

	result := db.Delete(&models.BookInfo{}, id) //soft deletion for hard deletion := db.Unscoped().delete()
	if result.Error != nil {
		log.Fatalf("failed to delete the book info")
	}
	rowsDeleted := result.RowsAffected
	return rowsDeleted
}
