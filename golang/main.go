package main

import (
	"expense-tracker/config"
	"expense-tracker/models"
	"expense-tracker/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db, err := config.InitDB() // Initialize the database connection
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})

	// Set up routes
	routes.SetupRoutes(r, db)

	r.Run(":8080")
}
