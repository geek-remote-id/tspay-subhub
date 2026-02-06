package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

// Connect initializes the database connection
func Connect(connStr string) (*sql.DB, error) {
	if connStr == "" {
		return nil, fmt.Errorf("connection string is empty")
	}

	// Add timezone parameter if not already present
	// This ensures all timestamps are in Asia/Jakarta timezone (UTC+7)
	if !contains(connStr, "timezone=") {
		if contains(connStr, "?") {
			connStr += "&timezone=Asia/Jakarta"
		} else {
			connStr += "?timezone=Asia/Jakarta"
		}
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return db, nil
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
