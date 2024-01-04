package main

import (
	"fmt"
	"libraryManagementSystem/database"
)

func main() {
	db := database.Database_connection()
	fmt.Println("connected to the databse")
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		defer sqlDB.Close()

		fmt.Println("Database connection closed")
	}()
}
