package main

import (
	"fmt"
	"libraryManagementSystem/database"
	"log"
)

func main() {
	db, err := database.Database()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Successfully connected to the postgresql database")

}
