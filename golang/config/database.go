package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// InitDB initializes and returns a connection to the database
func InitDB() (*gorm.DB, error) {
	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	sslmode := "disable"
	port := os.Getenv("DATABASE_PORT")
	if port == "" {
		port = "5432" // Default to 5432 if not specified
	}

	// Connect to PostgreSQL without specifying the database name to check for its existence
	defaultConnStr := fmt.Sprintf("host=%s port=%s user=%s sslmode=%s password=%s", host, port, user, sslmode, password)
	db, err := sql.Open("postgres", defaultConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Check if the database exists
	var exists int
	db.QueryRow("SELECT 1 FROM pg_database WHERE datname=$1", dbname).Scan(&exists)

	// If the database does not exist, create it
	if exists == 0 {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			return nil, err
		}
	}

	// Now, connect to the newly created or verified existing database using GORM
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", host, port, user, dbname, sslmode, password)
	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
