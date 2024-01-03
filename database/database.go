package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"			//the import statement ensures that when program starts, the init() function of the package gets executed and it registers itself as the driver for psql with sql
	"os"
)

func Database() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %w", err)
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Create the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error connectiong to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging databse: %w", err)
	}

	return db, nil

}
