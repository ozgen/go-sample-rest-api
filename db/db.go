package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"go-sample-rest-api/config"
	"go-sample-rest-api/logging"
)

// NewPostgresStorageConn establishes a new connection to the PostgreSQL database
// and returns the sql.DB object to be used by the application.
func NewPostgresStorageConn(conf config.Config) (*sql.DB, error) {
	log := logging.GetLogger()
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		conf.DBUser, conf.DBPassword, conf.DBName, conf.DBHost, conf.DBPort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.WithFields(logrus.Fields{
			"host":  conf.DBHost,
			"port":  conf.DBPort,
			"error": err,
		}).Error("Failed to open database connection")
		return nil, err
	}

	// Check if the connection is successful
	if err = db.Ping(); err != nil {
		log.WithFields(logrus.Fields{
			"host":  conf.DBHost,
			"port":  conf.DBPort,
			"error": err,
		}).Error("Failed to ping database")
		db.Close() // Only close the db on error where it cannot be used.
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"host": conf.DBHost,
		"port": conf.DBPort,
	}).Info("Successfully connected to the db!")
	return db, nil
}
