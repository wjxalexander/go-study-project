package store

import (
	"database/sql"
	"fmt"

	// implicit import for driver registration
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		// https://pkg.go.dev/fmt#Errorf
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Add enhanced configuration to the connection pool settings with:
	// db.SetMaxOpenConns(), db.SetMaxIdleConns(), and db.SetConnMaxIdleTime()
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}
	fmt.Println("Database opened successfully")
	return db, nil
}
