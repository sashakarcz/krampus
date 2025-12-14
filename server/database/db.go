package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize opens the database connection and runs migrations
func Initialize(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// Run migrations
	if err = RunMigrations(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
