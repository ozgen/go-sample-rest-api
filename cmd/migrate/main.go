package main

import (
	"go-sample-rest-api/config"
	"go-sample-rest-api/db"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {

	// Establish a new PostgreSQL connection
	db, err := db.NewPostgresStorageConn(config.Envs)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close() // Ensure the database connection is closed when main finishes

	// Check the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	// Setup the migration driver
	driver, err := postgres.WithInstance(db, &postgres.Config{}) // Ensure db.DB is the actual sql.DB instance
	if err != nil {
		log.Fatalf("Failed to create a migration driver: %v", err)
	}

	// Create a new migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations", // Ensure the path is correct
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create a migrate instance: %v", err)
	}

	// Execute migration based on command argument
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "up":
			err = m.Up()
			handleMigrationError(err)
		case "down":
			err = m.Down()
			handleMigrationError(err)
		default:
			log.Println("Invalid command")
			log.Println("Usage: migrate [up|down]")
		}

		// Output current migration version and status
		version, dirty, err := m.Version()
		if err != nil {
			log.Printf("Could not retrieve migration version: %v", err)
		} else {
			log.Printf("Current migration version: %d, dirty state: %v", version, dirty)
		}
	} else {
		log.Println("No command specified")
		log.Println("Usage: migrate [up|down]")
	}
}

func handleMigrationError(err error) {
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No change in migration needed.")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		log.Println("Migration successful.")
	}
}
