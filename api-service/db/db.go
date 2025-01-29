package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	for i := 0; i < 30; i++ { // try in 30 sec
		err = connectDB()
		if err == nil {
			log.Println("Successfully connected to database")
			return
		}
		log.Printf("Failed to connect to database, attempt %d/30: %v", i+1, err)
		time.Sleep(time.Second) // wait 1 second until the next attempt
	}
	log.Fatal("Could not connect to database after 30 attempts:", err)
}

func connectDB() error {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// check connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	return nil
}
