package db

import (
	"database/sql"
)

// DB is an interface for interacting with the database
type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// SQLDB implements the DB interface using an *sql.DB object
type SQLDB struct {
	db *sql.DB
}

// NewSQLDB creates a new SQLDB
func NewSQLDB(db *sql.DB) *SQLDB {
	return &SQLDB{db: db}
}

// Query executes a SQL query that returns rows, typically a SELECT
func (d *SQLDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

// QueryRow executes a SQL query that is expected to return at most one row
func (d *SQLDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
}

// Exec executes a SQL query that doesn't return rows, typically an INSERT, UPDATE, or DELETE
func (d *SQLDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}
