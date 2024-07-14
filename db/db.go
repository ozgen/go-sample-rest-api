package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"go-sample-rest-api/config"
)

// NewPostgresStorageConn establishes a new connection to the PostgreSQL database
// and returns the sql.DB object to be used by the application.
func NewPostgresStorageConn(conf config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		conf.DBUser, conf.DBPassword, conf.DBName, conf.DBHost, conf.DBPort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, err
	}

	// Check if the connection is successful
	if err = db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		db.Close() // Only close the db on error where it cannot be used.
		return nil, err
	}
	log.Print("Successfully connected to db!")

	return db, nil
}
